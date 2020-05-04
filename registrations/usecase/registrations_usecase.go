package usecase

import (
	"gade/srv-goldcard/activations"
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/process_handler"
	"gade/srv-goldcard/registrations"
	"gade/srv-goldcard/retry"
	"gade/srv-goldcard/transactions"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type registrationsUseCase struct {
	regRepo  registrations.Repository
	rrr      registrations.RestRepository
	phUC     process_handler.UseCase
	tUseCase transactions.UseCase
	arRepo   activations.RestRepository
}

// RegistrationsUseCase represent Registrations Use Case
func RegistrationsUseCase(regRepo registrations.Repository, rrr registrations.RestRepository, phUC process_handler.UseCase, tUseCase transactions.UseCase, arRepo activations.RestRepository) registrations.UseCase {
	return &registrationsUseCase{regRepo, rrr, phUC, tUseCase, arRepo}
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
	_, err = reg.regRepo.UpdateCardLimit(c, acc, false)

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

func (reg *registrationsUseCase) FinalRegistration(c echo.Context, pl models.PayloadAppNumber, fn models.FuncAfterGC) error {
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

	// Generate Application Form BRI Document
	err = reg.GenerateApplicationFormDocument(c, acc)

	if err != nil {
		return err
	}

	// Generate Slip TE Document
	err = reg.GenerateSlipTEDocument(c, acc)

	if err != nil {
		return err
	}

	// open and lock gold limit to core
	errAppBri := make(chan error)
	errAppCore := make(chan error)
	accChan := make(chan models.Account)
	go func() {
		// this validation for check is core already open before
		if acc.Application.CoreOpen {
			errAppCore <- nil
			return
		}

		err = reg.rrr.OpenGoldcard(c, acc, false)

		if err != nil {
			// insert error to process handler
			// change error status become true on table proess_statuses
			go reg.upsertProcessHandler(c, &acc, err)
			logger.Make(c, nil).Debug(err)
			errAppCore <- err
			return
		}
		// update Core open Status
		reg.coreOpenStatus(c, acc)
		errAppCore <- nil
	}()

	// concurrently update application status and current_step
	go func() {
		acc.Application.SetStatus(models.AppStatusProcessed)
		acc.Application.CurrentStep = models.AppStepCompleted
		_ = reg.regRepo.UpdateAppStatus(c, acc.Application)
		_ = reg.regRepo.UpdateApplication(c, acc.Application, []string{"status", "current_step"})
		accChannel, err := reg.CheckApplication(c, pl)

		if err != nil {
			accChan <- models.Account{}
		}

		accChan <- accChannel
	}()

	// channeling after core open goldcard finish
	err = fn(c, &acc, briPl, accChan, errAppBri, errAppCore)

	if err != nil {
		return err
	}

	return nil
}

func (reg *registrationsUseCase) FinalRegistrationPdsApi(c echo.Context, pl models.PayloadAppNumber) error {
	err := reg.FinalRegistration(c, pl, reg.concurrentlyAfterOpenGoldcard)

	if err != nil {
		return err
	}

	return nil
}

func (reg *registrationsUseCase) FinalRegistrationScheduler(c echo.Context, pl models.PayloadAppNumber) error {
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	err = reg.FinalRegistration(c, pl, reg.afterOpenGoldcard)

	if err != nil {
		// counter error on table process_statuses
		go reg.phUC.UpdateCounterError(c, acc)
		return err
	}

	// update error status to false on table process_statuses.
	err = reg.phUC.UpdateErrorStatus(c, acc)

	if err != nil {
		return err
	}

	return nil
}

func (reg *registrationsUseCase) concurrentlyAfterOpenGoldcard(c echo.Context, acc *models.Account,
	briPl models.PayloadBriRegister, accChan chan models.Account, errAppBri, errAppCore chan error) error {
	go func() {
		_ = reg.afterOpenGoldcard(c, acc, briPl, accChan, errAppBri, errAppCore)
	}()

	return nil
}

func (reg *registrationsUseCase) upsertProcessHandler(c echo.Context, acc *models.Account, errCore error) {
	var ps models.ProcessStatus
	err := ps.MapInsertProcessStatus(models.FinalAppProcessType, models.ApplicationTableName, acc.Application.ID, errCore)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return
	}

	err = reg.phUC.PostProcessHandler(c, ps)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return
	}
}

// function to update status core open if success
func (reg *registrationsUseCase) coreOpenStatus(c echo.Context, acc models.Account) {
	err := reg.regRepo.UpdateCoreOpen(c, &acc)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return
	}
}

// Function to Generate Application Form
func (reg *registrationsUseCase) GenerateApplicationFormDocument(c echo.Context, acc models.Account) error {
	// Get Document (ktp, npwp, selfie, slip_te, and app_form)
	docs, err := reg.regRepo.GetDocumentByApplicationId(acc.ApplicationID)

	if err != nil {
		return models.ErrGetDocument
	}

	// Mapping Application Form Data and Generate PDF
	appFormData := models.ApplicationForm{}

	paramsAppForm := map[string]interface{}{
		"docs": docs,
		"acc":  acc,
	}

	err = appFormData.MappingApplicationForm(paramsAppForm)

	if err != nil {
		return models.ErrMappingData
	}

	// concurrently insert or update all possible documents
	go retry.DoConcurrent(c, "upsertDocument", func() error {
		return reg.upsertDocument(c, appFormData.Account.Application)
	})

	return nil
}

// Function to Generate Slip TE
func (reg *registrationsUseCase) GenerateSlipTEDocument(c echo.Context, acc models.Account) error {
	// get user effective balance
	userDetail, err := reg.arRepo.GetDetailGoldUser(c, acc.Application.SavingAccount)

	if err != nil {
		return err
	}

	if _, ok := userDetail["saldoEfektif"].(string); !ok {
		return err
	}

	goldEffBalance, err := strconv.ParseFloat(userDetail["saldoEfektif"].(string), 64)

	if err != nil {
		return err
	}

	// Get Signatory Name for Slip TE Document
	signatoryName, err := reg.regRepo.GetSignatoryNameParam(c)

	if err != nil {
		return err
	}

	// Get Signatory Nip for Slip TE Document
	signatoryNip, err := reg.regRepo.GetSignatoryNipParam(c)

	if err != nil {
		return err
	}

	// Mapping Application Form Data and Generate PDF
	slipTeData := models.SlipTE{}

	paramsSlipTe := map[string]interface{}{
		"acc":            acc,
		"signatoryName":  signatoryName,
		"signatoryNip":   signatoryNip,
		"goldEffBalance": goldEffBalance,
	}

	err = slipTeData.MappingSlipTe(paramsSlipTe)

	if err != nil {
		logger.Make(nil, nil).Debug(err)
		return models.ErrMappingData
	}

	//  insert or update all possible documents
	err = reg.upsertDocument(c, slipTeData.Account.Application)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}

	return nil
}

func (reg *registrationsUseCase) ResetRegistration(c echo.Context, pl models.PayloadAppNumber) error {
	// Get account
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	// Deactive account
	err = reg.regRepo.DeactiveAccount(c, acc)

	if err != nil {
		return err
	}

	return nil
}
