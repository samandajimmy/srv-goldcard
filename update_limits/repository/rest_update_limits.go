package repository

import (
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/update_limits"
	"strconv"

	"github.com/labstack/echo"
)

type restUpdateLimits struct{}

// NewRestActivations will create an object that represent the activations.RestRepository interface
func NewRestUpdateLimits() update_limits.RestRepository {
	return &restUpdateLimits{}
}

func (rul *restUpdateLimits) CorePostUpdateLimit(c echo.Context, savingAccNum string, card models.Card, cif string) error {
	const (
		isBlokirTrue      = "1"
		isRecalculateTrue = "1"
		reqStatusRequest  = "RQ"
	)

	r := api.SwitchingResponse{}
	body := map[string]interface{}{
		"isBlokir":         isBlokirTrue,
		"noRek":            savingAccNum,
		"cif":              cif,
		"nominalTransaksi": strconv.FormatInt(card.CardLimit, 10),
		"isRecalculate":    isRecalculateTrue,
		"reqStatus":        reqStatusRequest,
	}

	req := api.MappingRequestSwitching(body)
	err := api.RetryableSwitchingPost(c, req, "/goldcard/updateLimit/register", &r)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	if r.ResponseCode != "00" {
		logger.Make(c, nil).Debug(models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode, r.ResponseDesc}))

		return models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode,
			r.ResponseDesc})
	}

	return nil
}

func (rul *restUpdateLimits) BRIPostUpdateLimit(c echo.Context, acc models.Account, doc models.Document) error {
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey":         acc.BrixKey,
		"limit":           acc.Card.CardLimit,
		"productRequest":  models.DefBriProductRequestUpLimit,
		"handPhoneNumber": acc.PersonalInformation.HandPhoneNumber,
		"email":           acc.PersonalInformation.Email,
		"file":            map[string]string{"OTHR": models.MapBRIExtBase64File[doc.FileExtension] + doc.FileBase64},
	}

	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, "/v1/cobranding/limit/update", reqBRIBody.RequestData, &respBRI)

	// response code SD when try to attempt update limit to BRI more than one times in a day
	if respBRI.ResponseCode == "SD" {
		return models.ErrSameDayUpdateLimitAttempt
	}

	if errBRI != nil {
		logger.Make(c, nil).Debug(errBRI)

		return models.ErrPostUpdateLimitToBRI
	}

	if respBRI.ResponseCode != "00" {
		logger.Make(c, nil).Debug(models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{respBRI.ResponseCode, respBRI.ResponseMessage}))

		return models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{respBRI.ResponseCode,
			respBRI.ResponseMessage})
	}

	return nil
}
