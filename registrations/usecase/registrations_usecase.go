package usecase

import (
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"gade/srv-goldcard/retry"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type registrationsUseCase struct {
	regRepo registrations.Repository
	rrr     registrations.RestRepository
}

// RegistrationsUseCase represent Registrations Use Case
func RegistrationsUseCase(regRepo registrations.Repository, rrr registrations.RestRepository) registrations.UseCase {
	return &registrationsUseCase{regRepo, rrr}
}

func (reg *registrationsUseCase) PostAddress(c echo.Context, pl models.PayloadAddress) error {
	// get core service health status
	if err := reg.regRepo.GetCoreServiceStatus(c); err != nil {
		return err
	}

	// get account by appNumber
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	// validate new data address or not
	if pl.IsNew == models.UseExistingAddress {
		return nil
	}

	err = acc.MappingAddressData(c, pl)

	if err != nil {
		return models.ErrMappingData
	}

	// get zipcode
	addrData := models.AddressData{City: pl.AddressCity, Province: pl.Province,
		Subdistrict: pl.Subdistrict, Village: pl.Village}
	zipcode, err := reg.regRepo.GetZipcode(c, addrData)

	if err != nil {
		return models.ErrZipcodeNotFound
	}

	acc.Correspondence.Zipcode = zipcode
	err = reg.regRepo.PostAddress(c, acc)

	if err != nil {
		return models.ErrPostAddressFailed
	}

	// update app current step
	acc.Application.CurrentStep = models.AppStepAddress
	_ = reg.regRepo.UpdateApplication(c, acc.Application, []string{"current_step"})

	return nil
}

func (reg *registrationsUseCase) PostRegistration(c echo.Context, payload models.PayloadRegistration) (models.RespRegistration, error) {
	var respRegNil models.RespRegistration
	acc, err := reg.CheckApplication(c, payload)

	if err != nil && err != models.ErrAppNumberNotFound {
		return respRegNil, err
	}

	// if application exist, return app status
	if acc.ID != 0 {
		return models.RespRegistration{
			ApplicationNumber: acc.Application.ApplicationNumber,
			ApplicationStatus: acc.Application.Status,
			CurrentStep:       acc.Application.CurrentStep,
		}, nil
	}

	appNumber, _ := uuid.NewRandom()

	// get BRI bank_id
	bankID, err := reg.regRepo.GetBankIDByCode(c, models.BriBankCode)

	if err != nil {
		return respRegNil, models.ErrBankNotFound
	}

	// get pegadaian emergency_contact_id
	ecID, err := reg.regRepo.GetEmergencyContactIDByType(c, models.EmergencyContactDef)

	if err != nil {
		return respRegNil, models.ErrEmergecyContactNotFound
	}

	app := models.Applications{ApplicationNumber: appNumber.String(), Status: models.AppStatusOngoing}
	acc = models.Account{CIF: payload.CIF, BranchCode: payload.BranchCode, ProductRequest: models.DefBriProductRequest,
		BillingCycle: models.DefBriBillingCycle, CardDeliver: models.DefBriCardDeliver, BankID: bankID, EmergencyContactID: ecID}
	pi := models.PersonalInformation{HandPhoneNumber: payload.HandPhoneNumber}
	err = reg.regRepo.CreateApplication(c, app, acc, pi)

	if err != nil {
		return respRegNil, models.ErrCreateApplication
	}

	return models.RespRegistration{ApplicationNumber: app.ApplicationNumber}, nil
}

func (reg *registrationsUseCase) PostPersonalInfo(c echo.Context, pl models.PayloadPersonalInformation) error {
	// check duplication/blacklist by BRI
	resp := api.BriResponse{}
	requestData := map[string]interface{}{
		"nik":       pl.Nik,
		"birthDate": pl.BirthDate,
	}
	reqBody := api.BriRequest{RequestData: requestData}
	err := api.RetryableBriPost(c, "/v1/cobranding/deduplication", reqBody, &resp)

	if err != nil {
		return err
	}

	// get account by appNumber
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	err = acc.MappingRegistrationData(c, pl)

	if err != nil {
		return models.ErrMappingData
	}

	// get zipcode
	addrData := models.AddressData{City: pl.AddressCity, Province: pl.Province,
		Subdistrict: pl.Subdistrict, Village: pl.Village}
	zipcode, err := reg.regRepo.GetZipcode(c, addrData)

	if err != nil {
		return models.ErrZipcodeNotFound
	}

	// update account data
	acc.PersonalInformation.Zipcode = zipcode
	err = reg.regRepo.UpdateAllRegistrationData(c, acc)

	if err != nil {
		return models.ErrUpdateRegData
	}

	// concurrently insert or update all possible documents
	go retry.DoConcurrent(c, "upsertDocument", func() error {
		return reg.upsertDocument(c, acc.Application)
	})

	// update app current step
	acc.Application.CurrentStep = models.AppStepPersonalInfo
	// concurrently update current_step
	go func() {
		_ = reg.regRepo.UpdateApplication(c, acc.Application, []string{"current_step"})
	}()

	return nil
}

func (reg *registrationsUseCase) PostCardLimit(c echo.Context, pl models.PayloadCardLimit) error {
	// get account by appNumber
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	// Gold Card inquiry Registrations to core
	r := api.SwitchingResponse{}
	body := map[string]interface{}{
		"noRek": acc.Application.SavingAccount,
	}
	req := api.MappingRequestSwitching(body)
	err = api.RetryableSwitchingPost(c, req, "/goldcard/inquiry", &r)

	if err != nil {
		return err
	}

	if r.ResponseCode != api.SwitchingRCInquiryAllow {
		logger.Make(c, nil).Debug(models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode, r.ResponseDesc}))

		return models.ErrInquiryReg
	}

	// get current STL
	currStl, err := reg.rrr.GetCurrentGoldSTL(c)

	if err != nil {
		return models.ErrGetCurrSTL
	}

	acc.Card.CardLimit = pl.CardLimit
	acc.Card.GoldLimit = acc.Card.SetGoldLimit(pl.CardLimit, currStl) // convert limit to gold with current STL added with reserved locking balance
	err = reg.regRepo.UpdateCardLimit(c, acc)

	if err != nil {
		return models.ErrUpdateCardLimit
	}

	// Get STL Price
	go reg.updateSTLPrice(c, acc)

	// update app current step
	acc.Application.CurrentStep = models.AppStepCardLimit
	// concurrently update current_step
	go func() {
		_ = reg.regRepo.UpdateApplication(c, acc.Application, []string{"current_step"})
	}()

	return nil
}

func (reg *registrationsUseCase) PostOccupation(c echo.Context, pl models.PayloadOccupation) error {
	var city string
	var zipcode string
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	err = acc.Occupation.MappingOccupation(pl)

	if err != nil {
		return models.ErrMappingData
	}

	city, zipcode, err = reg.regRepo.GetCityFromZipcode(c, acc)

	if err != nil {
		return models.ErrMappingData
	}

	acc.Occupation.OfficeCity = city
	acc.Occupation.OfficeZipcode = zipcode
	acc.Occupation.JobTitle = models.DefJobTitle

	err = reg.regRepo.PostOccupation(c, acc)

	if err != nil {
		return models.ErrUpdateOccData
	}

	// update app current step
	acc.Application.CurrentStep = models.AppStepOccupation
	_ = reg.regRepo.UpdateApplication(c, acc.Application, []string{"current_step"})

	return nil
}

func (reg *registrationsUseCase) PostSavingAccount(c echo.Context, pl models.PayloadSavingAccount) error {
	// get account by appNumber
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	acc.Application.SavingAccount = pl.AccountNumber

	err = reg.regRepo.PostSavingAccount(c, acc)

	if err != nil {
		return models.ErrPostSavingAccountFailed
	}

	// update app current step
	acc.Application.CurrentStep = models.AppStepSavingAcc
	// concurrently update current_step
	go func() {
		_ = reg.regRepo.UpdateApplication(c, acc.Application, []string{"current_step"})
	}()

	return nil
}

func (reg *registrationsUseCase) FinalRegistration(c echo.Context, pl models.PayloadAppNumber) error {
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	// get account by appNumber
	briPl, err := reg.regRepo.GetAllRegData(c, pl.ApplicationNumber)

	if err != nil {
		return models.ErrAppData
	}

	// validasi bri register payload
	if err := c.Validate(briPl); err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	// open and lock gold limit to core
	errAppBri := make(chan error)
	errAppCore := make(chan error)
	accChan := make(chan models.Account)

	go func() {
		err := reg.rrr.OpenGoldcard(c, acc, false)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			errAppCore <- err
			return
		}

		errAppCore <- nil
	}()

	// channeling after core open goldcard finish
	go reg.afterOpenGoldcard(c, &acc, briPl, accChan, errAppBri, errAppCore)

	// concurrently update application status and current_step
	go func() {
		acc.Application.SetStatus(models.AppStatusProcessed)
		acc.Application.CurrentStep = models.AppStepCompleted
		_ = reg.regRepo.UpdateAppStatus(c, acc.Application)
		_ = reg.regRepo.UpdateApplication(c, acc.Application, []string{"status", "current_step"})
		accChannel, _ := reg.CheckApplication(c, pl)
		accChan <- accChannel
	}()

	return nil
}
