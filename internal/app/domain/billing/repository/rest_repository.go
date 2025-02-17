package repository

import (
	"srv-goldcard/internal/app/domain/billing"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/api"
	"srv-goldcard/internal/pkg/logger"

	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

type restBillings struct{}

// NewRestBillings will create an object that represent the activation.RestRepository interface
func NewRestBillings() billing.RestRepository {
	return &restBillings{}
}

func (ra *restBillings) GetBillingsStatement(c echo.Context, acc model.Account) (model.BillingStatement, error) {
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
		return model.BillingStatement{}, nil
	}

	if errBRI != nil {
		logger.Make(c, nil).Debug(respBRI)

		return model.BillingStatement{}, errBRI
	}

	response := respBRI.DataOne["listOfStatements"].(map[string]interface{})["statementHeader"].(map[string]interface{})

	return model.BillingStatement{
		BillingAmount:     int64(response["currentBalance"].(float64)),
		BillingPrintDate:  response["statementDate"].(string),
		BillingDueDate:    response["paymentDueDate"].(string),
		BillingMinPayment: int64(response["monthlyPayment"].(float64)),
	}, nil
}
