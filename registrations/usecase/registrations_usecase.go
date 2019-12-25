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
func RegistrationsUseCase(
	regRepository registrations.Repository,
) registrations.UseCase {
	return &registrationsUseCase{
		registrationsRepository: regRepository,
	}
}

// PostAddress representation update address to database
func (reg *registrationsUseCase) PostAddress(c echo.Context, registrations *models.Registrations) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	err := reg.registrationsRepository.PostAddress(c, registrations)

	if err != nil {
		requestLogger.Debug(models.ErrPostAddressFailed)

		return models.ErrPostAddressFailed
	}

	return nil
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
