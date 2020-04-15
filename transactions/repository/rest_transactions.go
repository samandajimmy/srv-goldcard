package repository

import (
	"encoding/json"
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/transactions"
	"strconv"

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

func (rt *restTransactions) CorePaymentInquiry(c echo.Context, pl models.PlPaymentInquiry) (map[string]interface{}, error) {
	response := map[string]interface{}{}
	respSwitching := api.SwitchingResponse{}
	requestDataSwitching := map[string]interface{}{
		"amount":         pl.PaymentAmount,
		"jenisTransaksi": "CC",
		"norek":          pl.AccountNumber,
	}

	req := api.MappingRequestSwitching(requestDataSwitching)
	errSwitching := api.RetryableSwitchingPost(c, req, "/gadai/inquiry", &respSwitching)

	if errSwitching != nil {
		return response, errSwitching
	}

	if respSwitching.ResponseCode != api.APIRCSuccess {
		logger.Make(c, nil).Debug(models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc}))
		return response, models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc})
	}

	return respSwitching.ResponseData, nil
}

func (rt *restTransactions) PostPaymentTransactionToCore(c echo.Context, bill models.Billing) error {
	respSwitching := api.SwitchingResponse{}
	requestDataSwitching := map[string]interface{}{
		"cif":          bill.Account.CIF,
		"noRek":        bill.Account.Application.SavingAccount,
		"nominal":      strconv.FormatInt(bill.DebtAmount, 10),
		"norekTagihan": bill.Account.AccountNumber,
		"reffBiller":   bill.RefBilling,
		"isUpdate":     models.BillIsUpdate,
	}

	req := api.MappingRequestSwitching(requestDataSwitching)
	errSwitching := api.RetryableSwitchingPost(c, req, "/goldcard/transaksi/sendTagihan", &respSwitching)

	if errSwitching != nil {
		return errSwitching
	}

	if respSwitching.ResponseCode != api.APIRCSuccess {
		logger.Make(c, nil).Debug(models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc}))
		return models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc})
	}

	return nil
}
