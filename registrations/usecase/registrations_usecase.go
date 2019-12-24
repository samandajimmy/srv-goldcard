package usecase

import (
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type registrationUseCase struct {
}

// NewRegistrationUseCase will create new an registrationUseCase object representation of registrations.UseCase interface
func NewRegistrationUseCase() registrations.UseCase {
	return &registrationUseCase{}
}

func (reg *registrationUseCase) sendApplicationNotif(payload map[string]string) error {
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
