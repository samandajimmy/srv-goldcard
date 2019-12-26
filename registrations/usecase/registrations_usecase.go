package usecase

import (
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type registrationsUseCase struct {
	registrationsRepository registrations.Repository
}

// RegistrationsUseCase represent Registrations Use Case
func RegistrationsUseCase(
	regRepository registrations.Repository,
) registrations.UseCase {
	return &registrationsUseCase{
		registrationsRepository: regRepository,
	}
}

// PostAddress representation update address to database
func (reg *registrationsUseCase) PostAddress(c echo.Context, registrations *models.Registrations) (string, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	err := reg.registrationsRepository.PostAddress(c, registrations)

	if err != nil {
		requestLogger.Debug(models.ErrPostAddressFailed)

		return "", models.ErrPostAddressFailed
	}

	return "", nil
}

// PostAddress representation get address from database
func (reg *registrationsUseCase) GetAddress(c echo.Context, phoneNo string) (map[string]interface{}, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	res, err := reg.registrationsRepository.GetAddress(c, phoneNo)

	response := map[string]interface{}{"address": res}

	if err != nil {
		requestLogger.Debug(models.ErrPostAddressFailed)

		return nil, models.ErrPostAddressFailed
	}

	return response, nil
}

func (reg *registrationsUseCase) sendApplicationNotif(payload map[string]string) error {
	response := map[string]interface{}{}
	pds, err := models.NewPdsAPI(echo.MIMEApplicationForm)

	if err != nil {
		logrus.Debug(err)

		return err
	}

	req, err := pds.Request("/goldcard/status_pengajuan_notif", echo.POST, payload)

	if err != nil {
		logrus.Debug(err)

		return err
	}

	_, err = pds.Do(req, &response)

	if err != nil {
		logrus.Debug(err)

		return err
	}

	return nil
}
