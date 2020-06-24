package repository

import (
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/billings"

	"github.com/labstack/echo"
	"time"
	"strconv"
	"fmt"
)

type restBillings struct{}

// NewRestBillings will create an object that represent the activations.RestRepository interface
func NewRestBillings() billings.RestRepository {
	return &restBillings{}
}

func (ra *restBillings) GetBillingsStatement(c echo.Context, acc models.Account) (map[string]interface{}, string) {
	dateNow := time.Now()
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey": acc.BrixKey,
		"year" : strconv.Itoa(dateNow.Year()),
		"month" : fmt.Sprintf("%02d", dateNow.Month()),
	}

	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, "/v1/cobranding/trx/inquiry", reqBRIBody.RequestData, &respBRI)

	if errBRI != nil {
		logger.Make(c, nil).Debug(errBRI)

		return respBRI.DataOne, respBRI.ResponseCode
	}

	return respBRI.DataOne, ""
}