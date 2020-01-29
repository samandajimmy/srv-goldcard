package usecase

import (
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"gade/srv-goldcard/retry"
	"reflect"

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
	// get account by appNumber
	acc, err := reg.checkApplication(c, pl)

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
	acc, err := reg.checkApplication(c, payload)

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
	acc, err := reg.checkApplication(c, pl)

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
	acc, err := reg.checkApplication(c, pl)
	r := api.SwitchingResponse{}

	if err != nil {
		return err
	}

	// Gold Card inquiry Registrations to core
	body := map[string]interface{}{
		"noRek": acc.Application.SavingAccount,
	}
	req := api.MappingRequestSwitching(body)
	err = api.RetryableSwitchingPost(c, req, "/goldcard/inquiry", &r)

	if err != nil {
		return err
	}

	// Validation response
	if r.ResponseCode != api.SwitchingRCInquiryAllow {
		return models.ErrInquiryReg
	}

	acc.Card.CardLimit = pl.CardLimit
	acc.Card.GoldLimit = pl.GoldLimit
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

// PostOccupation representation update occupation to database
func (reg *registrationsUseCase) PostOccupation(c echo.Context, pl models.PayloadOccupation) error {
	var city string
	var zipcode string
	acc, err := reg.checkApplication(c, pl)

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

// PostAddress representation update address to database
func (reg *registrationsUseCase) PostSavingAccount(c echo.Context, pl models.PayloadSavingAccount) error {
	// get account by appNumber
	acc, err := reg.checkApplication(c, pl)

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
	acc, err := reg.checkApplication(c, pl)

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
		return err
	}

	// open and lock gold limit to core
	if err := reg.rrr.OpenGoldcard(c, acc, false); err != nil {
		return err
	}

	// concurrently apply the goldcard application to BRI
	go reg.briApply(c, &acc, briPl)

	// concurrently update application status and current_step
	go func() {
		acc.Application.SetStatus(models.AppStatusProcessed)
		acc.Application.CurrentStep = models.AppStepCompleted
		_ = reg.regRepo.UpdateAppStatus(c, acc.Application)
		_ = reg.regRepo.UpdateApplication(c, acc.Application, []string{"status", "current_step"})
	}()

	return nil
}

func (reg *registrationsUseCase) GetAppStatus(c echo.Context, pl models.PayloadAppNumber) (models.AppStatus, error) {
	var appStatus models.AppStatus
	// Get account by app number
	acc, err := reg.checkApplication(c, pl)

	if err != nil {
		return appStatus, err
	}

	// concurrently get app status from BRI API then update to our DB
	go func() {
		resp := api.BriResponse{}
		reqBody := map[string]interface{}{
			"briXkey": acc.BrixKey,
		}

		err := api.BriPost(c, "/v1/cobranding/card/appstatus", reqBody, &resp)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			return
		}

		// update application status
		data := resp.Data[0]
		if _, ok := data["appStatus"].(string); !ok {
			logger.Make(c, nil).Debug(err)
			return
		}

		acc.Application.ID = acc.ApplicationID
		acc.Application.SetStatus(data["appStatus"].(string))
		err = reg.regRepo.UpdateAppStatus(c, acc.Application)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			return
		}
	}()

	appStatus, err = reg.regRepo.GetAppStatus(c, acc.Application)

	if err != nil {
		return appStatus, models.ErrGetAppStatus
	}

	return appStatus, nil
}

func (reg *registrationsUseCase) briApply(c echo.Context, acc *models.Account, pl models.PayloadBriRegister) {
	err := reg.briRegister(c, acc, pl)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return
	}

	// upload document to BRI API
	err = reg.uploadAppDocs(c, acc)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return
	}
}

func (reg *registrationsUseCase) briRegister(c echo.Context, acc *models.Account, pl models.PayloadBriRegister) error {
	if acc.BrixKey != "" {
		return nil
	}

	resp := api.BriResponse{}
	reqBody := api.BriRequest{RequestData: pl}
	err := api.RetryableBriPost(c, "/v1/cobranding/register", reqBody, &resp)

	if err != nil {
		return err
	}

	// update brixkey id
	if _, ok := resp.DataOne["briXkey"].(string); !ok {
		logger.Make(c, nil).Debug(models.ErrSetVar)

		return models.ErrSetVar
	}

	acc.BrixKey = resp.DataOne["briXkey"].(string)
	// concurrently update brixkey from BRI API
	go func() {
		_ = reg.regRepo.UpdateBrixkeyID(c, *acc)
	}()

	return nil
}

func (reg *registrationsUseCase) uploadAppDocs(c echo.Context, acc *models.Account) error {
	for _, doc := range acc.Application.Documents {
		// concurrently upload application documents to BRI
		go func(doc models.Document) {
			_ = reg.uploadAppDoc(c, acc.BrixKey, doc)
		}(doc)
	}

	return nil
}

func (reg *registrationsUseCase) uploadAppDoc(c echo.Context, brixkey string, doc models.Document) error {
	if doc.DocID != "" {
		return nil
	}

	briReq := models.AppDocument{
		BriXkey:    brixkey,
		DocType:    models.DefAppDocType,
		FileName:   doc.FileName,
		FileExt:    doc.FileExtension,
		Base64file: "data:image/jpeg;base64," + doc.FileBase64,
	}

	resp := api.BriResponse{}
	reqBody := api.BriRequest{RequestData: briReq}
	err := api.RetryableBriPost(c, "/v1/cobranding/document", reqBody, &resp)

	if err != nil {
		return err
	}

	if _, ok := resp.DataOne["documentId"].(string); !ok {
		return models.ErrDocIDNotFound
	}

	doc.DocID = resp.DataOne["documentId"].(string)
	// concurrently insert or update application document
	go func() {
		_ = reg.regRepo.UpsertAppDocument(c, doc)
	}()

	return nil
}

func (reg *registrationsUseCase) checkApplication(c echo.Context, pl interface{}) (models.Account, error) {
	r := reflect.ValueOf(pl)
	appNumber := r.FieldByName("ApplicationNumber")

	if appNumber.IsZero() {
		return models.Account{}, nil
	}

	acc := models.Account{Application: models.Applications{ApplicationNumber: appNumber.String()}}
	err := reg.regRepo.GetAccountByAppNumber(c, &acc)

	if err != nil {
		return models.Account{}, models.ErrAppNumberNotFound
	}

	if acc.BrixKey != "" {
		return models.Account{}, models.ErrAppNumberCompleted
	}

	return acc, nil
}

func (reg *registrationsUseCase) upsertDocument(c echo.Context, app models.Applications) error {
	if len(app.Documents) == 0 {
		return nil
	}

	for _, doc := range app.Documents {
		err := reg.regRepo.UpsertAppDocument(c, doc)

		if err != nil {
			return err
		}
	}

	return nil
}

func (reg *registrationsUseCase) updateSTLPrice(c echo.Context, acc models.Account) {
	hargeEmas, err := reg.rrr.GetCurrentGoldSTL(c)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return
	}

	acc.Card.CurrentSTL = hargeEmas
	err = reg.regRepo.UpdateCardLimit(c, acc)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return
	}
}
