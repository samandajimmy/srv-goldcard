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

func (reg *registrationsUseCase) PostAddress(c echo.Context, registrations *models.Registrations) error {
	err := reg.regRepo.PostAddress(c, registrations)

	if err != nil {
		return models.ErrPostAddressFailed
	}

	return nil
}

func (reg *registrationsUseCase) GetAddress(c echo.Context, phoneNo string) (map[string]interface{}, error) {
	res, err := reg.regRepo.GetAddress(c, phoneNo)

	if err != nil {
		return nil, models.ErrPostAddressFailed
	}

	response := map[string]interface{}{"address": res}

	return response, nil
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

func (reg *registrationsUseCase) PostPersonalInfo(c echo.Context, payload models.PayloadPersonalInformation) error {
	// get account by appNumber
	acc, err := reg.regRepo.GetAccountByAppNumber(c, payload.ApplicationNumber)

	if err != nil {
		return models.ErrAppNumberNotFound
	}

	err = acc.MappingRegistrationData(c, payload)

	if err != nil {
		return models.ErrMappingData
	}

	// update account data
	err = reg.regRepo.UpdateAllRegistrationData(c, acc)

	if err != nil {
		return models.ErrUpdateRegData
	}

	return nil
}

// PostAddress representation update address to database
func (reg *registrationsUseCase) PostSavingAccount(c echo.Context, applications *models.Applications) error {
	err := reg.regRepo.PostSavingAccount(c, applications)

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
