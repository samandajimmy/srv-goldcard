package repository

import (
	"encoding/json"
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

func (ra *restTransactions) GetBRICardInformation(c echo.Context, acc models.Account) (models.BRICardBalance, error) {
	var briCardBal models.BRICardBalance
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey": acc.BrixKey,
	}

	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, "/v1/cobranding/card/information", reqBRIBody.RequestData, &respBRI)

	if errBRI != nil {
		logger.Make(c, nil).Debug(errBRI)

		return briCardBal, errBRI
	}

	mrshlCardInfo, err := json.Marshal(respBRI.DataOne)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return briCardBal, err
	}

	err = json.Unmarshal(mrshlCardInfo, &briCardBal)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return briCardBal, err
	}

	return briCardBal, nil
}
