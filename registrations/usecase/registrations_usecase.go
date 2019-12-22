package usecase

import (
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"

	"github.com/labstack/echo"
)

type registrationsUseCase struct {
	registrationsUseCase    registrations.UseCase
	registrationsRepository registrations.Repository
}

// RegistrationsUseCase represent Registrations Use Case
func RegistrationsUseCase() registrations.UseCase {
	return &registrationsUseCase{}
}

func (reg *registrationsUseCase) PostAddress(c echo.Context, registrations *models.Registrations) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	if registrations.PhoneNumber == "" {
		requestLogger.Debug(models.ErrBadParamInput)

		return models.ErrBadParamInput
	}

	if registrations.ResidenceAddress == "" {
		requestLogger.Debug(models.ErrAddressEmpty)

		return models.ErrAddressEmpty
	}

	err := reg.registrationsRepository.PostAddress(c, registrations)

	if err != nil {
		requestLogger.Debug(models.ErrPostAddressFailed)

		return models.ErrPostAddressFailed
	}

	return nil
}
