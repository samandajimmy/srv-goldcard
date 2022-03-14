package repository

import (
	"encoding/json"
	"srv-goldcard/internal/app/domain/transaction"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/api"
	"srv-goldcard/internal/pkg/logger"
	"strconv"

	"github.com/labstack/echo"
)

type restTransactions struct{}

// NewRestActivations will create an object that represent the activation.RestRepository interface
func NewRestTransactions() transaction.RestRepository {
	return &restTransactions{}
}

func (ra *restTransactions) GetBRICardInformation(c echo.Context, acc model.Account) (model.BRICardBalance, error) {
	var briCardBal model.BRICardBalance
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey": acc.BrixKey,
	}

	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, "/card/information", reqBRIBody.RequestData, &respBRI)

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

func (rt *restTransactions) CorePaymentInquiry(c echo.Context, pl model.PlPaymentInquiry, acc model.Account) (map[string]interface{}, error) {
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
		logger.Make(c, nil).Debug(model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc}))
		return response, model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc})
	}

	return respSwitching.ResponseData, nil
}

func (rt *restTransactions) PostPaymentTransactionToCore(c echo.Context, bill model.Billing) error {
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
		logger.Make(c, nil).Debug(model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc}))
		return model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc})
	}

	return nil
}

func (rt *restTransactions) PostPaymentCoreNotif(c echo.Context, acc model.Account, pl model.PlPaymentTrxCore) error {
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
		logger.Make(c, nil).Debug(model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc}))
		return model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc})
	}

	return nil
}

func (rt *restTransactions) PostPaymentBRI(c echo.Context, acc model.Account, amount int64) error {
	resp := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey":        acc.BrixKey,
		"amount":         amount,
		"productRequest": acc.ProductRequest,
	}
	errBRI := api.RetryableBriPost(c, "/limit/payment", requestDataBRI, &resp)

	if errBRI != nil {
		logger.Make(c, nil).Debug(errBRI)
		return errBRI
	}

	return nil
}

// GetBRIPendingTrx to get pending trx for single account from BRI
func (rt *restTransactions) GetBRIPendingTrx(c echo.Context, acc model.Account, startDate string, endDate string) (model.RespBRIPendingTrxData, error) {
	respBRIPendTrxData := model.RespBRIPendingTrxData{}
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey":   acc.BrixKey,
		"startDate": startDate,
		"endDate":   endDate,
	}
	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, "/trx/pending", reqBRIBody.RequestData, &respBRI)

	if respBRI.ResponseCode == "NF" {
		return respBRIPendTrxData, nil
	}

	if errBRI != nil {
		logger.Make(c, nil).Debug(errBRI)

		return respBRIPendTrxData, errBRI
	}

	if respBRI.ResponseCode != "00" {
		logger.Make(c, nil).Debug(model.DynamicErr(model.ErrBriAPIRequest, []interface{}{respBRI.ResponseCode, respBRI.ResponseData}))

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
func (rt *restTransactions) GetBRIPostedTrx(c echo.Context, briXkey string) (model.RespBRIPostedTransaction, error) {
	respBRIPosted := model.RespBRIPostedTransaction{}
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey": briXkey,
	}
	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, "/trx/trx_posting", reqBRIBody.RequestData, &respBRI)

	if respBRI.ResponseCode == "5X" {
		return respBRIPosted, nil
	}

	if errBRI != nil {
		logger.Make(c, nil).Debug(errBRI)

		return respBRIPosted, errBRI
	}

	if respBRI.ResponseCode != "00" {
		logger.Make(c, nil).Debug(model.DynamicErr(model.ErrBriAPIRequest, []interface{}{respBRI.ResponseCode, respBRI.ResponseData}))

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
