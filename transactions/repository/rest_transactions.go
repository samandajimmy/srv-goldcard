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

func (rt *restTransactions) CorePaymentInquiry(c echo.Context, pl models.PlPaymentInquiry, acc models.Account) (map[string]interface{}, error) {
	response := map[string]interface{}{}
	respSwitching := api.SwitchingResponse{}
	requestDataSwitching := map[string]interface{}{
		"cif":            acc.CIF,
		"noRek":          acc.Application.SavingAccount,
		"norekTagihan":   acc.AccountNumber,
		"nominal":        pl.PaymentAmount,
		"jenisTransaksi": "CC",
	}

	req := api.MappingRequestSwitching(requestDataSwitching)
	errSwitching := api.RetryableSwitchingPost(c, req, "/goldcard/transaksi/inquiryTagihan", &respSwitching)

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
		"cif":            bill.Account.CIF,
		"noRek":          bill.Account.Application.SavingAccount,
		"nominal":        strconv.FormatInt(bill.DebtAmount, 10),
		"norekTagihan":   bill.Account.AccountNumber,
		"branchCode":     bill.Account.BranchCode,
		"jenisTransaksi": "PAYMENT",
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

func (rt *restTransactions) PostPaymentCoreNotif(c echo.Context, acc models.Account, pl models.PlPaymentTrxCore) error {
	respSwitching := api.SwitchingResponse{}
	requestDataSwitching := map[string]interface{}{
		"noRek":         acc.Application.SavingAccount,
		"norekTagihan":  acc.AccountNumber,
		"nominal":       pl.PaymentAmount,
		"branchCode":    acc.BranchCode,
		"cif":           acc.CIF,
		"reffSwitching": pl.RefTrx,
	}

	req := api.MappingRequestSwitching(requestDataSwitching)
	errSwitching := api.RetryableSwitchingPost(c, req, "/goldcard/tagihan/payment", &respSwitching)

	if errSwitching != nil {
		return errSwitching
	}

	if respSwitching.ResponseCode != api.APIRCSuccess {
		logger.Make(c, nil).Debug(models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc}))
		return models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc})
	}

	return nil
}

func (rt *restTransactions) PostPaymentBRI(c echo.Context, acc models.Account, amount int64) error {
	resp := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey":        acc.BrixKey,
		"amount":         amount,
		"productRequest": acc.ProductRequest,
	}
	errBRI := api.RetryableBriPost(c, "/v1/cobranding/limit/payment", requestDataBRI, &resp)

	if errBRI != nil {
		logger.Make(c, nil).Debug(errBRI)
		return errBRI
	}

	return nil
}

// GetBRIPendingTrx to get pending trx for single account from BRI
func (rt *restTransactions) GetBRIPendingTrx(c echo.Context, acc models.Account, startDate string, endDate string) (models.RespBRIPendingTrxData, error) {
	respBRIPendTrxData := models.RespBRIPendingTrxData{}
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey":   acc.BrixKey,
		"startDate": startDate,
		"endDate":   endDate,
	}
	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, "/v1/cobranding/trx/pending", reqBRIBody.RequestData, &respBRI)

	if respBRI.ResponseCode == "NF" {
		return respBRIPendTrxData, nil
	}

	if errBRI != nil {
		logger.Make(c, nil).Debug(errBRI)

		return respBRIPendTrxData, errBRI
	}

	requestData := respBRI.DataOne["requestData"].([]interface{})
	mrshlBRIPendTrxInq, err := json.Marshal(requestData[0])

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return respBRIPendTrxData, err
	}

	err = json.Unmarshal(mrshlBRIPendTrxInq, &respBRIPendTrxData)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return respBRIPendTrxData, err
	}

	return respBRIPendTrxData, nil
}

// GetBRIPosted to get posted trx for single account from BRI
func (rt *restTransactions) GetBRIPostedTrx(c echo.Context, briXkey string) (models.RespBRIPostedTransaction, error) {
	respBRIPosted := models.RespBRIPostedTransaction{}
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey": briXkey,
	}
	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, "/v1/cobranding/trx/trx_posting", reqBRIBody.RequestData, &respBRI)

	if respBRI.ResponseCode == "5X" {
		return respBRIPosted, nil
	}

	if errBRI != nil {
		logger.Make(c, nil).Debug(errBRI)

		return respBRIPosted, errBRI
	}

	mrshlBRIBilInq, err := json.Marshal(respBRI.ResponseData)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return respBRIPosted, err
	}

	err = json.Unmarshal(mrshlBRIBilInq, &respBRIPosted)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return respBRIPosted, err
	}

	return respBRIPosted, nil
}
