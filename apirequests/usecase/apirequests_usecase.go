package usecase

import (
	"encoding/json"
	"gade/srv-goldcard/apirequests"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"net/http"
	"net/url"
	"reflect"

	"github.com/labstack/echo"
)

// ARUseCase is variable to store api_requests usecase
var ARUseCase apirequests.UseCase

type apirequestsUseCase struct {
	arRepo apirequests.Repository
}

// APIRequestsUseCase represent APIRequests Use Case
func APIRequestsUseCase(arRepo apirequests.Repository) apirequests.UseCase {
	return &apirequestsUseCase{
		arRepo: arRepo,
	}
}

func (arus *apirequestsUseCase) PostAPIRequest(c echo.Context, reqID string, statusCode int, api, req, resp interface{}) error {
	rapi := reflect.ValueOf(api)
	host := reflect.Indirect(rapi).FieldByName("Host").Interface().(*url.URL)
	status := "success"

	if statusCode != http.StatusOK {
		status = "error"
	}

	reqJSON, err := json.Marshal(req)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	respJSON, err := json.Marshal(resp)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	apiRequest := models.APIRequest{
		RequestID:    reqID,
		HostName:     host.Host,
		Endpoint:     reflect.Indirect(rapi).FieldByName("Endpoint").String(),
		Status:       status,
		RequestData:  reqJSON,
		ResponseData: respJSON,
	}

	err = arus.arRepo.InserAPIRequest(c, apiRequest)

	if err != nil {

		return err
	}

	return nil
}
