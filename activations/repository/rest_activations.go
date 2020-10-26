package repository

import (
	"gade/srv-goldcard/activations"
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"time"

	"github.com/labstack/echo"
)

type restActivations struct{}

// NewRestActivations will create an object that represent the activations.RestRepository interface
func NewRestActivations() activations.RestRepository {
	return &restActivations{}
}

func (ra *restActivations) GetDetailGoldUser(c echo.Context, accNumber string) (map[string]interface{}, error) {
	nilMap := map[string]interface{}{}
	r := api.SwitchingResponse{}
	reqBody := map[string]interface{}{"norek": accNumber}
	req := api.MappingRequestSwitching(reqBody)
	err := api.RetryableSwitchingPost(c, req, "/portofolio/dettabemas", &r)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return nilMap, err
	}

	if r.ResponseCode != api.APIRCSuccess {
		logger.Make(c, nil).Debug(models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode, r.ResponseDesc}))

		return nilMap, err
	}

	return r.ResponseData, nil
}

func (ra *restActivations) ActivationsToCore(c echo.Context, acc *models.Account) error {
	respSwitching := api.SwitchingResponse{}
	requestDataSwitching := map[string]interface{}{
		"cif":        acc.CIF,
		"noRek":      acc.Application.SavingAccount,
		"branchCode": acc.BranchCode,
	}

	req := api.MappingRequestSwitching(requestDataSwitching)
	errSwitching := api.RetryableSwitchingPost(c, req, "/goldcard/aktivasi", &respSwitching)

	if errSwitching != nil {
		return errSwitching
	}

	if respSwitching.ResponseCode != api.APIRCSuccess {
		logger.Make(c, nil).Debug(models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc}))
		return models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc})
	}

	if _, ok := respSwitching.ResponseData["noRekPembayaran"].(string); !ok {
		logger.Make(c, nil).Debug(models.ErrSetVar)

		return models.ErrSetVar
	}

	acc.AccountNumber = respSwitching.ResponseData["noRekPembayaran"].(string)

	return nil
}

func (ra *restActivations) ActivationsToBRI(c echo.Context, acc models.Account, pa models.PayloadActivations) error {
	respBRI := api.BriResponse{}
	date, err := time.Parse(models.DDMMYYYY, pa.BirthDate)

	if err != nil {
		return err
	}

	birthDate := date.Format(models.DateFormatDef)

	requestDataBRI := map[string]interface{}{
		"briXkey":        acc.BrixKey,
		"expDate":        pa.ExpDate,
		"lastFourDigits": pa.LastFourDigits,
		"firstSixDigits": pa.FirstSixDigits,
		"birthDate":      birthDate,
	}

	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, "/card/activation", reqBRIBody.RequestData, &respBRI)

	if errBRI != nil {
		return errBRI
	}

	return nil
}
