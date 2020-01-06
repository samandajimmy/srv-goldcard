package usecase

import (
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"

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
