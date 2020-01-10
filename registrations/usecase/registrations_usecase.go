package usecase

import (
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"reflect"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type registrationsUseCase struct {
	regRepo registrations.Repository
}

// RegistrationsUseCase represent Registrations Use Case
func RegistrationsUseCase(
	regRepository registrations.Repository,
) registrations.UseCase {
	return &registrationsUseCase{
		regRepo: regRepository,
	}
}

func (reg *registrationsUseCase) PostAddress(c echo.Context, pl models.PayloadAddress) error {
	// get account by appNumber
	acc, err := reg.regRepo.GetAccountByAppNumber(c, pl.ApplicationNumber)

	if err != nil {
		return models.ErrAppNumberNotFound
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

	return nil
}

func (reg *registrationsUseCase) PostRegistration(c echo.Context, payload models.PayloadRegistration) (string, error) {
	appNumber, _ := uuid.NewRandom()

	// get BRI bank_id
	bankID, err := reg.regRepo.GetBankIDByCode(c, models.BriBankCode)

	if err != nil {
		return "", models.ErrBankNotFound
	}

	// get pegadaian emergency_contact_id
	ecID, err := reg.regRepo.GetEmergencyContactIDByType(c, models.EmergencyContactDef)

	if err != nil {
		return "", models.ErrEmergecyContactNotFound
	}

	app := models.Applications{ApplicationNumber: appNumber.String()}
	acc := models.Account{CIF: payload.CIF, BankID: bankID, EmergencyContactID: ecID}
	pi := models.PersonalInformation{HandPhoneNumber: payload.HandPhoneNumber}
	err = reg.regRepo.CreateApplication(c, app, acc, pi)

	if err != nil {
		return "", models.ErrCreateApplication
	}

	return appNumber.String(), nil
}

func (reg *registrationsUseCase) PostPersonalInfo(c echo.Context, pl models.PayloadPersonalInformation) error {
	// check duplication/blacklist by BRI
	resp := models.BriResponse{}
	requestData := map[string]interface{}{
		"nik":       pl.Nik,
		"birthDate": pl.BirthDate,
	}
	reqBody := models.BriRequest{RequestData: requestData}
	err := models.BriPost("/v1/cobranding/deduplication", reqBody, &resp)

	if err != nil {
		return models.ErrExternalAPI
	}

	// to set variation of BRI response
	resp.SetRC()

	if resp.ResponseCode == "00" {
		return models.ErrBlacklisted
	}

	// get account by appNumber
	acc, err := reg.regRepo.GetAccountByAppNumber(c, pl.ApplicationNumber)

	if err != nil {
		return models.ErrAppNumberNotFound
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

	return nil
}

func (reg *registrationsUseCase) PostCardLimit(c echo.Context, pl models.PayloadCardLimit) error {
	// get account by appNumber
	acc, err := reg.regRepo.GetAccountByAppNumber(c, pl.ApplicationNumber)

	if err != nil {
		return models.ErrAppNumberNotFound
	}

	acc.Card.CardLimit = pl.CardLimit
	err = reg.regRepo.UpdateCardLimit(c, acc)

	if err != nil {
		return models.ErrUpdateRegData
	}

	return nil
}

// PostAddress representation update address to database
func (reg *registrationsUseCase) PostSavingAccount(c echo.Context, pl models.PayloadSavingAccount) error {
	// get account by appNumber
	acc, err := reg.regRepo.GetAccountByAppNumber(c, pl.ApplicationNumber)

	if err != nil {
		return models.ErrAppNumberNotFound
	}

	acc.Application.SavingAccount = pl.AccountNumber

	err = reg.regRepo.PostSavingAccount(c, acc)

	if err != nil {
		return models.ErrPostSavingAccountFailed
	}

	return nil
}

func (reg *registrationsUseCase) FinalRegistration(c echo.Context, pl models.PayloadAppNumber) error {
	// get account by appNumber
	briPl, err := reg.regRepo.GetAllRegData(c, pl.ApplicationNumber)
	acc, err := reg.regRepo.GetAccountByAppNumber(c, pl.ApplicationNumber)

	if err != nil {
		return models.ErrAppNumberNotFound
	}

	resp := models.BriResponse{}
	reqBody := models.BriRequest{RequestData: briPl}
	err = models.BriPost("/v1/cobranding/register", reqBody, &resp)

	if err != nil {
		return models.ErrExternalAPI
	}

	// to set variation of BRI response
	resp.SetRC()

	if resp.ResponseCode != "00" {
		return models.DynamicErr(models.ErrBriAPIRequest, []interface{}{resp.ResponseCode, resp.ResponseMessage})
	}

	// update brixkey id
	brixkey, ok := resp.ResponseData["briXkey"].(string)

	if !ok {
		logger.Make(c, nil).Debug(models.ErrSetVar)

		return models.ErrSetVar
	}

	acc.BrixKey = brixkey
	err = reg.regRepo.UpdateBrixkeyID(c, acc)

	if err != nil {
		return models.ErrUpdateBrixkey
	}

	// upload document to BRI API
	err = reg.uploadAppDoc(c, acc)

	if err != nil {
		return err
	}

	return nil
}

func (reg *registrationsUseCase) GetAppStatus(c echo.Context, pl models.PayloadAppNumber) (models.AppStatus, error) {
	var appStatus models.AppStatus
	// Get account by app number
	acc, err := reg.regRepo.GetAccountByAppNumber(c, pl.ApplicationNumber)

	if err != nil {
		return appStatus, models.ErrAppIDNotFound
	}

	// Request API BRI 
	resp := models.BriResponse{}
	requestData := map[string]interface{}{
		"briXkey": acc.BrixKey,
	}
	reqBody := models.BriRequest{RequestData: requestData}
	err = models.BriPost("/v1/cobranding/card/appstatus", reqBody, &resp)

	if err != nil {
		return appStatus, models.ErrExternalAPI
	}

	// to set variation of BRI response
	resp.SetRC()

	if resp.ResponseCode != "00" {
		return appStatus, models.DynamicErr(models.ErrBriAPIRequest, []interface{}{resp.ResponseCode, resp.ResponseMessage})
	}


	// Update Status Application
	if _, ok := resp.ResponseData["appStatus"].(string); !ok {
		return appStatus, models.ErrSetVar
	}

	acc.Application.Status = resp.ResponseData["appStatus"].(string)
	acc.Application.ID = acc.ApplicationID
	appStatus, err = reg.regRepo.UpdateAppStatus(c, acc.Application)

	if err != nil {
		return appStatus, models.ErrUpdateAppStatus
	}

	return appStatus, nil
} 

func (reg *registrationsUseCase) uploadAppDoc(c echo.Context, acc models.Account) error {
	var docs []models.AppDocument
	docNames := []string{"KtpImageBase64", "NpwpImageBase64", "SelfieImageBase64"}
	docIDs := []string{"KtpDocID", "NpwpDocID", "SelfieDocID"}
	app, err := reg.regRepo.GetAppByID(c, acc.ApplicationID)

	if err != nil {
		return models.ErrAppIDNotFound
	}

	r := reflect.ValueOf(&app)

	for idx, docName := range docNames {
		f := reflect.Indirect(r).FieldByName(docName)
		fDocID := reflect.Indirect(r).FieldByName(docIDs[idx])

		if f.IsZero() {
			continue
		}

		doc := models.AppDocument{
			BriXkey:    acc.BrixKey,
			DocType:    models.DefAppDocType,
			FileName:   docName,
			FileExt:    models.DefAppDocFileExt,
			Base64file: "data:image/jpeg;base64," + f.String(),
		}

		docs = append(docs, doc)
		resp := models.BriResponse{}
		reqBody := models.BriRequest{RequestData: doc}
		err = models.BriPost("/v1/cobranding/document", reqBody, &resp)

		if err != nil {
			return models.ErrExternalAPI
		}

		// to set variation of BRI response
		resp.SetRC()

		if resp.ResponseCode != "00" {
			return models.DynamicErr(models.ErrBriAPIRequest, []interface{}{resp.ResponseCode,
				resp.ResponseMessage})
		}

		if _, ok := resp.ResponseData["documentId"].(string); !ok {
			return models.ErrDocIDNotFound
		}

		fDocID.SetString(resp.ResponseData["documentId"].(string))
	}

	// update document id to application data
	err = reg.regRepo.UpdateAppDocID(c, app)

	if err != nil {
		return models.ErrUpdateAppDocID
	}

	return nil
}

func (reg *registrationsUseCase) sendApplicationNotif(payload map[string]string) error {
	response := map[string]interface{}{}
	pds, err := models.NewPdsAPI(echo.MIMEApplicationForm)

	if err != nil {
		return err
	}

	req, err := pds.Request("/goldcard/status_pengajuan_notif", echo.POST, payload)

	if err != nil {
		return err
	}

	_, err = pds.Do(req, &response)

	if err != nil {
		return err
	}

	return nil
}
