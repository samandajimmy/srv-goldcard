package repository

import (
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/billings"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"

	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

type restBillings struct{}

// NewRestBillings will create an object that represent the activations.RestRepository interface
func NewRestBillings() billings.RestRepository {
	return &restBillings{}
}

func (ra *restBillings) GetBillingsStatement(c echo.Context, acc models.Account) (models.BillingStatement, error) {
	dateNow := time.Now()
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey": acc.BrixKey,
		"year":    strconv.Itoa(dateNow.Year()),
		"month":   fmt.Sprintf("%02d", dateNow.Month()),
	}

	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, "/trx/inquiry", reqBRIBody.RequestData, &respBRI)

	if respBRI.ResponseCode == "5X" {
		return models.BillingStatement{}, nil
	}

	if errBRI != nil {
		logger.Make(c, nil).Debug(respBRI)

		return models.BillingStatement{}, errBRI
	}

	response := respBRI.DataOne["listOfStatements"].(map[string]interface{})["statementHeader"].(map[string]interface{})

	return models.BillingStatement{
		BillingAmount:     int64(response["currentBalance"].(float64)),
		BillingPrintDate:  response["statementDate"].(string),
		BillingDueDate:    response["paymentDueDate"].(string),
		BillingMinPayment: int64(response["monthlyPayment"].(float64)),
	}, nil
}
