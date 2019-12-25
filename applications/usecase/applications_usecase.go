package usecase

import (
	"gade/srv-goldcard/applications"
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

type applicationsUseCase struct {
	applicationsUseCase    applications.UseCase
	applicationsRepository applications.Repository
}

// RegistrationsUseCase represent Registrations Use Case
func ApplicationsUseCase(appliRepository applications.Repository) applications.UseCase {
	return &applicationsUseCase{
		applicationsRepository: appliRepository,
	}
}

// PostAddress representation update address to database
func (appli *applicationsUseCase) PostSavingAccount(c echo.Context, applications *models.Applications) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	if applications.ApplicationNumber == "" || applications.SavingAccount == "" {
		requestLogger.Debug(models.ErrBadParamInput)
		return models.ErrBadParamInput
	}

	err := appli.applicationsRepository.PostSavingAccount(c, applications)

	if err != nil {
		requestLogger.Debug(models.ErrSavingAccountFailed)

		return models.ErrSavingAccountFailed
	}

	return nil
}
