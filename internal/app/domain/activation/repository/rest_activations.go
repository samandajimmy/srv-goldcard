package repository

import (
	"srv-goldcard/internal/app/domain/activation"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/api"
	"srv-goldcard/internal/pkg/logger"
	"time"

	"github.com/labstack/echo"
)

type restActivations struct{}

// NewRestActivations will create an object that represent the activation.RestRepository interface
func NewRestActivations() activation.RestRepository {
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
		logger.Make(c, nil).Debug(model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode, r.ResponseDesc}))

		return nilMap, err
	}

	return r.ResponseData, nil
}

func (ra *restActivations) ActivationsToCore(c echo.Context, acc *model.Account) error {
	respSwitching := api.SwitchingResponse{}
	cardNumber := ""

	if (acc.Card.CardStatus != model.CardStatuses{}) {
		cardNumber = acc.Card.CardStatus.LastEncryptedCardNumber
	}

	requestDataSwitching := map[string]interface{}{
		"cif":        acc.CIF,
		"noRek":      acc.Application.SavingAccount,
		"branchCode": acc.BranchCode,
		"cardNumber": cardNumber,
	}

	req := api.MappingRequestSwitching(requestDataSwitching)
	errSwitching := api.RetryableSwitchingPost(c, req, "/goldcard/aktivasi", &respSwitching)

	if errSwitching != nil {
		return errSwitching
	}

	if respSwitching.ResponseCode != api.APIRCSuccess {
		logger.Make(c, nil).Debug(model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc}))
		return model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc})
	}

	if _, ok := respSwitching.ResponseData["noRekPembayaran"].(string); !ok {
		logger.Make(c, nil).Debug(model.ErrSetVar)

		return model.ErrSetVar
	}

	acc.AccountNumber = respSwitching.ResponseData["noRekPembayaran"].(string)

	return nil
}

func (ra *restActivations) ActivationsToBRI(c echo.Context, acc model.Account, pa model.PayloadActivations) error {
	respBRI := api.BriResponse{}
	date, err := time.Parse(model.DDMMYYYY, pa.BirthDate)

	if err != nil {
		return err
	}

	birthDate := date.Format(model.DateFormatDef)

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
		logger.Make(c, nil).Debug(err)

		return errBRI
	}

	return nil
}
