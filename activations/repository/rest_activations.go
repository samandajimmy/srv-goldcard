package repository

import (
	"gade/srv-goldcard/activations"
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

type restActivations struct{}

// NewRestActivations will create an object that represent the activations.RestRepository interface
func NewRestActivations() activations.RestRepository {
	return &restActivations{}
}

func (ra *restActivations) GetDetailGoldUser(c echo.Context, accNumber string) (map[string]string, error) {
	nilMap := map[string]string{}
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

func (ra *restActivations) ActivationsToCore(c echo.Context, acc models.Account) error {
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

	return nil
}

func (ra *restActivations) ActivationsToBRI(c echo.Context, acc models.Account, pa models.PayloadActivations) error {
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey":       acc.BrixKey,
		"expDate":       pa.ExpDate,
		"lastSixDigits": pa.LastSixDigits,
	}
	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, "/v1/cobranding/card/activation", reqBRIBody.RequestData, &respBRI)

	if errBRI != nil {
		return errBRI
	}

	return nil
}
