package http

import (
	"gade/srv-goldcard/applications"
	"gade/srv-goldcard/models"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

var response models.Response

// ApplicationsHandler represent the httphandler for registrations
type ApplicationsHandler struct {
	applicationsUseCase applications.UseCase
}

// NewApplicationsHandler represent to registration gold card
func NewApplicationsHandler(echoGroup models.EchoGroup, appliUseCase applications.UseCase) {
	handler := &ApplicationsHandler{
		applicationsUseCase: appliUseCase,
	}

	// End Point For CMS
	echoGroup.API.POST("/saving_account", handler.PostSavingAccount)
}

// PostAddress a handler to update Address in table personal_informations
func (appli *ApplicationsHandler) PostSavingAccount(c echo.Context) error {
	respErrors := &models.ResponseErrors{}
	logger := models.RequestLogger{}
	response = models.Response{}
	var applications models.Applications

	c.Bind(&applications)
	logger.DataLog(c, applications).Info("Start of Post Saving Account")
	err := appli.applicationsUseCase.PostSavingAccount(c, &applications)

	if err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)
		logger.DataLog(c, response).Info("End of Post Saving Account")

		return c.JSON(getStatusCode(err), response)
	}

	response.Code = "00"
	response.Status = models.StatusSuccess
	response.Message = models.MessageUpdateSuccess
	logger.DataLog(c, response).Info("End of Post Saving Account")
	return c.JSON(getStatusCode(err), response)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if strings.Contains(err.Error(), "400") {
		return http.StatusBadRequest
	}

	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusOK
	}
}
