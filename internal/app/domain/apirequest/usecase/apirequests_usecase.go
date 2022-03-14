package usecase

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"

	"srv-goldcard/internal/app/domain/apirequest"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"

	"github.com/labstack/echo"
)

// ARUseCase is variable to store api_requests usecase
var ARUseCase apirequest.UseCase

type apirequestsUseCase struct {
	arRepo apirequest.Repository
}

// APIRequestsUseCase represent APIRequests Use Case
func APIRequestsUseCase(arRepo apirequest.Repository) apirequest.UseCase {
	return &apirequestsUseCase{
		arRepo: arRepo,
	}
}

func (arus *apirequestsUseCase) PostAPIRequest(c echo.Context, statusCode int, api, req, resp interface{}) error {
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

	apiRequest := model.APIRequest{
		RequestID:    logger.GetEchoRID(c),
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
