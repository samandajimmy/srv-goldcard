package usecase

import (
	"os"
	"reflect"
	"strconv"
	"time"

	"srv-goldcard/internal/app/domain/activation"
	"srv-goldcard/internal/app/domain/process_handler"
	"srv-goldcard/internal/app/domain/registration"
	"srv-goldcard/internal/app/domain/transaction"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/api"
	"srv-goldcard/internal/pkg/logger"
	"srv-goldcard/internal/pkg/retry"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type registrationsUseCase struct {
	regRepo  registration.Repository
	rrr      registration.RestRepository
	phUC     process_handler.UseCase
	tUseCase transaction.UseCase
	arRepo   activation.RestRepository
}

// RegistrationsUseCase represent Registrations Use Case
func RegistrationsUseCase(regRepo registration.Repository, rrr registration.RestRepository, phUC process_handler.UseCase, tUseCase transaction.UseCase, arRepo activation.RestRepository) registration.UseCase {
	return &registrationsUseCase{regRepo, rrr, phUC, tUseCase, arRepo}
}

func (reg *registrationsUseCase) PostAddress(c echo.Context, pl model.PayloadAddress) error {
	// get core service health status
	if err := reg.regRepo.GetCoreServiceStatus(c); err != nil {
		return err
	}

	// get account by appNumber
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	// set card deliver
	acc.CardDeliver = pl.CardDeliver
	err = reg.regRepo.PostAddress(c, acc)

	if err != nil {
		return model.ErrPostAddressFailed
	}

	// update app current step
	acc.Application.CurrentStep = model.AppStepAddress
	_ = reg.regRepo.UpdateApplication(c, acc.Application, []string{"current_step"})

	return nil
}

func (reg *registrationsUseCase) GetAddress(c echo.Context, pl model.PayloadAppNumber) (model.RespGetAddress, error) {
	respGetAddr := model.RespGetAddress{}
	// get account by appNumber
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return respGetAddr, err
	}

	pi := acc.PersonalInformation
	occ := acc.Occupation

	// set card deliver
	respGetAddr.CardDeliver = acc.CardDeliver

	// set office address data in resp get address struct
	respGetAddr.Office.AddressLine1 = occ.OfficeAddress1
	respGetAddr.Office.AddressLine2 = occ.OfficeAddress2
	respGetAddr.Office.AddressLine3 = occ.OfficeAddress3
	respGetAddr.Office.Province = occ.OfficeProvince
	respGetAddr.Office.City = occ.OfficeCity
	respGetAddr.Office.Subdistrict = occ.OfficeSubdistrict
	respGetAddr.Office.Village = occ.OfficeVillage
	respGetAddr.Office.Zipcode = occ.OfficeZipcode

	// set domicile address data in resp get address struct
	respGetAddr.Domicile.AddressLine1 = pi.AddressLine1
	respGetAddr.Domicile.AddressLine2 = pi.AddressLine2
	respGetAddr.Domicile.AddressLine3 = pi.AddressLine3
	respGetAddr.Domicile.Province = pi.AddressProvince
	respGetAddr.Domicile.City = pi.AddressCity
	respGetAddr.Domicile.Subdistrict = pi.AddressSubdistrict
	respGetAddr.Domicile.Village = pi.AddressVillage
	respGetAddr.Domicile.Zipcode = pi.Zipcode

	return respGetAddr, nil
}

func (reg *registrationsUseCase) PostRegistration(c echo.Context, payload model.PayloadRegistration) (model.RespRegistration, error) {
	var respRegNil model.RespRegistration

	r := reflect.ValueOf(payload)
	payloadAppNumber := r.FieldByName("ApplicationNumber")
	app := model.Applications{}

	if payloadAppNumber.IsZero() {
		app = reg.CheckApplicationByCIF(c, payload)
	}

	if app.ID != 0 {
		return model.RespRegistration{
			ApplicationNumber: app.ApplicationNumber,
			ApplicationStatus: app.Status,
			CurrentStep:       app.CurrentStep,
		}, nil
	}

	acc, err := reg.CheckApplication(c, payload)

	if err != nil && err != model.ErrAppNumberNotFound {
		return respRegNil, err
	}

	// if application exist, return app status
	if acc.ID != 0 {
		return model.RespRegistration{
			ApplicationNumber: acc.Application.ApplicationNumber,
			ApplicationStatus: acc.Application.Status,
			CurrentStep:       acc.Application.CurrentStep,
		}, nil
	}

	appNumber, _ := uuid.NewRandom()

	// get BRI bank_id
	bankID, err := reg.regRepo.GetBankIDByCode(c, model.BriBankCode)

	if err != nil {
		return respRegNil, model.ErrBankNotFound
	}

	// get pegadaian emergency_contact_id
	ecID, err := reg.regRepo.GetEmergencyContactIDByType(c, model.EmergencyContactDef)

	if err != nil {
		return respRegNil, model.ErrEmergecyContactNotFound
	}

	// set expiryAt time
	now := time.Now()
	expiryDur, _ := strconv.ParseInt(os.Getenv(`APP_TIMEOUT_DURATION`), 10, 64)
	expiryAt := now.Add(time.Duration(expiryDur) * time.Second)

	app = model.Applications{
		ApplicationNumber: appNumber.String(),
		Status:            model.AppStatusOngoing,
		ExpiredAt:         expiryAt,
		CreatedAt:         now,
	}
	pi := model.PersonalInformation{HandPhoneNumber: payload.HandPhoneNumber}
	acc = model.Account{CIF: payload.CIF, BranchCode: payload.BranchCode, ProductRequest: model.DefBriProductRequest,
		BillingCycle: model.DefBriBillingCycle, CardDeliver: model.BriCardDeliverHome, BankID: bankID,
		EmergencyContactID: ecID, Application: app, PersonalInformation: pi}
	err = reg.regRepo.CreateApplication(c, app, acc, pi)

	if err != nil {
		return respRegNil, model.ErrCreateApplication
	}

	// run job for expiration process on background
	diff := app.ExpiredAt.Sub(app.CreatedAt)
	delay := time.Duration(diff.Seconds())
	reg.appTimeoutJob(c, acc, diff, delay)

	return model.RespRegistration{
		ApplicationNumber: app.ApplicationNumber,
		ApplicationStatus: app.Status,
	}, nil
}

func (reg *registrationsUseCase) PostPersonalInfo(c echo.Context, pl model.PayloadPersonalInformation) error {
	// check duplication/blacklist by BRI
	resp := api.BriResponse{}
	requestData := map[string]interface{}{
		"nik":       pl.Nik,
		"birthDate": pl.BirthDate,
	}
	reqBody := api.BriRequest{RequestData: requestData}
	err := api.RetryableBriPost(c, "/deduplication", reqBody, &resp)

	if err != nil {
		return err
	}

	// get account by appNumber
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	// get zipcode
	addrData := model.AddressData{City: pl.AddressCity, Province: pl.Province,
		Subdistrict: pl.Subdistrict, Village: pl.Village, AddressLine1: pl.AddressLine1}
	zipcode, err := reg.regRepo.GetZipcode(c, addrData)

	if err != nil {
		return model.ErrZipcodeNotFound
	}

	pl.Zipcode = zipcode
	addrData.Zipcode = zipcode
	err = acc.MappingRegistrationData(pl, addrData)

	if err != nil {
		return err
	}

	// update account data
	err = reg.regRepo.UpdateAllRegistrationData(c, acc)

	if err != nil {
		return model.ErrUpdateRegData
	}

	// concurrently insert or update all possible documents
	go retry.DoConcurrent(c, "upsertDocument", func() error {
		return reg.upsertDocument(c, acc.Application)
	})

	// update app current step
	acc.Application.CurrentStep = model.AppStepPersonalInfo
	// concurrently update current_step
	go func() {
		_ = reg.regRepo.UpdateApplication(c, acc.Application, []string{"current_step"})
	}()

	return nil
}

func (reg *registrationsUseCase) PostCardLimit(c echo.Context, pl model.PayloadCardLimit) error {
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
		logger.Make(c, nil).Debug(model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode, r.ResponseDesc}))

		return model.ErrInquiryReg
	}

	// get current STL
	currStl, err := reg.rrr.GetCurrentGoldSTL(c)

	if err != nil {
		return model.ErrGetCurrSTL
	}

	acc.Card.CardLimit = pl.CardLimit
	acc.Card.GoldLimit = acc.Card.SetGoldLimit(pl.CardLimit, currStl)
	acc.Card.StlLimit = currStl
	acc.Card.StlBalance = currStl
	acc.Application.CardLimit = pl.CardLimit
	err = reg.regRepo.UpdateCardLimit(c, acc, nil)

	if err != nil {
		return model.ErrUpdateCardLimit
	}

	// update app current step
	acc.Application.CurrentStep = model.AppStepCardLimit
	// concurrently update current_step
	go func() {
		_ = reg.regRepo.UpdateApplication(c, acc.Application, []string{"current_step"})
	}()

	return nil
}

func (reg *registrationsUseCase) PostOccupation(c echo.Context, pl model.PayloadOccupation) error {
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	// get zipcode
	addrData := model.AddressData{City: pl.OfficeCity, Province: pl.OfficeProvince,
		Subdistrict: pl.OfficeSubdistrict, Village: pl.OfficeVillage,
		AddressLine1: pl.OfficeAddress1}
	// job category need company name
	inclCompanyName := []int{1, 4, 5}
	// if job category is must included with company name
	if result := model.ArrayContains(inclCompanyName, int(pl.JobCategory)); result {
		addrData.AddressLine1 = pl.Company + " " + pl.OfficeAddress1
	}

	addrData.Zipcode, err = reg.regRepo.GetZipcode(c, addrData)

	if err != nil {
		return model.ErrZipcodeNotFound
	}

	pl.JobTitle = model.DefJobTitle
	err = acc.Occupation.MappingOccupation(pl, addrData)

	if err != nil {
		return err
	}

	err = reg.regRepo.PostOccupation(c, acc)

	if err != nil {
		return model.ErrUpdateOccData
	}

	// update app current step
	acc.Application.CurrentStep = model.AppStepOccupation
	_ = reg.regRepo.UpdateApplication(c, acc.Application, []string{"current_step"})

	return nil
}

func (reg *registrationsUseCase) PostSavingAccount(c echo.Context, pl model.PayloadSavingAccount) error {
	// get account by appNumber
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	acc.Application.SavingAccount = pl.AccountNumber
	err = reg.regRepo.PostSavingAccount(c, acc)

	if err != nil {
		return model.ErrPostSavingAccountFailed
	}

	// update app current step
	acc.Application.CurrentStep = model.AppStepSavingAcc
	// concurrently update current_step
	go func() {
		_ = reg.regRepo.UpdateApplication(c, acc.Application, []string{"current_step"})
	}()

	return nil
}

func (reg *registrationsUseCase) FinalRegistration(c echo.Context, pl model.PayloadAppNumber, fn model.FuncAfterGC) error {
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	// get account by appNumber
	briPl, err := reg.regRepo.GetAllRegData(c, pl.ApplicationNumber)

	if err != nil {
		return model.ErrAppData
	}

	// validasi bri register payload
	if err := c.Validate(briPl); err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}

	// generate application document and slip te
	err = reg.generateOtherDocs(c, &acc)

	if err != nil {
		return err
	}

	// open and lock gold limit to core
	accChan := make(chan model.Account)

	// do open goldcard on core
	err = reg.openCore(c, &acc)

	if err != nil {
		// send notif app failed
		_ = reg.appNotification(c, acc, "failed", true)

		return err
	}

	// concurrently update application status and current_step
	go func() {
		acc.Application.SetStatus(model.AppStatusProcessed)
		acc.Application.CurrentStep = model.AppStepCompleted
		_ = reg.regRepo.UpdateAppStatus(c, acc.Application)
		_ = reg.regRepo.UpdateApplication(c, acc.Application, []string{"status", "current_step"})
		accChannel, err := reg.CheckApplication(c, pl)

		if err != nil {
			accChan <- model.Account{}
		}

		accChan <- accChannel
	}()

	// channeling after core open goldcard finish
	err = fn(c, &acc, briPl, accChan)

	if err != nil {
		return err
	}

	return nil
}

func (reg *registrationsUseCase) openCore(c echo.Context, acc *model.Account) error {
	// this validation for check is core already open before
	if acc.Application.CoreOpen {
		return nil
	}

	err := reg.rrr.OpenGoldcard(c, *acc, false)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		// set application status rejected
		acc.Application.SetStatus(model.AppStatusRejected)
		_ = reg.regRepo.UpdateAppStatus(c, acc.Application)

		return model.ErrCoreOpen
	}

	// update Core open Status
	_ = reg.regRepo.UpdateCoreOpen(c, acc)

	return nil
}

func (reg *registrationsUseCase) FinalRegistrationPdsApi(c echo.Context, pl model.PayloadAppNumber) error {
	err := reg.FinalRegistration(c, pl, reg.concurrentlyAfterOpenGoldcard)

	if err != nil {
		return err
	}

	return nil
}

func (reg *registrationsUseCase) FinalRegistrationScheduler(c echo.Context, pl model.PayloadAppNumber) error {
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	err = reg.FinalRegistration(c, pl, reg.afterOpenGoldcard)

	if err != nil {
		return err
	}

	// update error status to false on table process_statuses.
	err = reg.phUC.UpdateErrorStatus(c, acc)

	if err != nil {
		return err
	}

	return nil
}

func (reg *registrationsUseCase) concurrentlyAfterOpenGoldcard(c echo.Context, acc *model.Account,
	briPl model.PayloadBriRegister, accChan chan model.Account) error {
	go func() {
		_ = reg.afterOpenGoldcard(c, acc, briPl, accChan)
	}()

	return nil
}

// Function to Generate Application Form
func (reg *registrationsUseCase) GenerateApplicationFormDocument(c echo.Context, acc *model.Account) error {
	if len(acc.Application.Documents) == 0 {
		return model.ErrMappingData
	}

	// Mapping Application Form Data and Generate PDF
	appFormData := model.ApplicationForm{}
	appFormData.Account = *acc
	err := appFormData.MappingApplicationForm()

	if err != nil {
		return model.ErrMappingData
	}

	err = reg.upsertDocument(c, appFormData.Account.Application)

	if err != nil {
		return model.ErrMappingData
	}

	*acc = appFormData.Account

	return nil
}

// Function to Generate Slip TE
func (reg *registrationsUseCase) GenerateSlipTEDocument(c echo.Context, acc *model.Account) error {
	if len(acc.Application.Documents) == 0 {
		return model.ErrMappingData
	}

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
	slipTeData := model.SlipTE{}
	slipTeData.Account = *acc

	paramsSlipTe := map[string]interface{}{
		"signatoryName":  signatoryName,
		"signatoryNip":   signatoryNip,
		"goldEffBalance": goldEffBalance,
	}

	err = slipTeData.MappingSlipTe(paramsSlipTe)

	if err != nil {
		logger.Make(nil, nil).Debug(err)
		return model.ErrMappingData
	}

	//  insert or update all possible documents
	err = reg.upsertDocument(c, slipTeData.Account.Application)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}

	*acc = slipTeData.Account

	return nil
}

func (reg *registrationsUseCase) ResetRegistration(c echo.Context, pl model.PayloadAppNumber) error {
	// Get account
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	// send close goldcard account to core
	err = reg.rrr.CloseGoldcard(c, acc)

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

func (reg *registrationsUseCase) ForceDeliver(c echo.Context, pl model.PayloadAppNumber) error {
	// Get account
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return err
	}

	err = reg.regRepo.ForceDeliverAccount(c, acc)

	if err != nil {
		return err
	}

	return nil
}
