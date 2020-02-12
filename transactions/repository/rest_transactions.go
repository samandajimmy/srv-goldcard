package repository

import (
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/transactions"

	"github.com/labstack/echo"
)

type restTransactions struct{}

// NewRestActivations will create an object that represent the activations.RestRepository interface
func NewRestTransactions() transactions.RestRepository {
	return &restTransactions{}
}

func (ra *restTransactions) GetBRICardInformation(c echo.Context, acc models.Account) (map[string]interface{}, error) {
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey": acc.BrixKey,
	}
	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, "/v1/cobranding/card/information", reqBRIBody.RequestData, &respBRI)

	if errBRI != nil {
		return nil, errBRI
	}

	if respBRI.ResponseCode != "00" {
		logger.Make(c, nil).Debug(models.DynamicErr(models.ErrBriAPIRequest, []interface{}{respBRI.ResponseCode, respBRI.ResponseData}))

		return nil, errBRI
	}

	return respBRI.DataOne, nil
}
