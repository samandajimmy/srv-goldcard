package repository

import (
	"srv-goldcard/internal/app/domain/update_limit"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/api"
	"srv-goldcard/internal/pkg/logger"
	"strconv"

	"github.com/labstack/echo"
)

type restUpdateLimits struct{}

// NewRestActivations will create an object that represent the activation.RestRepository interface
func NewRestUpdateLimits() update_limit.RestRepository {
	return &restUpdateLimits{}
}

func (rul *restUpdateLimits) CorePostUpdateLimit(c echo.Context, savingAccNum string, card model.Card, cif string) error {
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
		logger.Make(c, nil).Debug(model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode, r.ResponseDesc}))

		return model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode,
			r.ResponseDesc})
	}

	return nil
}

func (rul *restUpdateLimits) BRIPostUpdateLimit(c echo.Context, acc model.Account, doc model.Document) error {
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey":         acc.BrixKey,
		"limit":           acc.Card.CardLimit,
		"productRequest":  model.DefBriProductRequest,
		"handPhoneNumber": acc.PersonalInformation.HandPhoneNumber,
		"email":           acc.PersonalInformation.Email,
		"file":            map[string]string{"OTHR": model.MapBRIExtBase64File[doc.FileExtension] + doc.FileBase64},
	}

	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, "/limit/update", reqBRIBody.RequestData, &respBRI)

	// response code SD when try to attempt update limit to BRI more than one times in a day
	if respBRI.ResponseCode == "SD" {
		return model.ErrSameDayUpdateLimitAttempt
	}

	if errBRI != nil {
		logger.Make(c, nil).Debug(errBRI)

		return model.ErrPostUpdateLimitToBRI
	}

	if respBRI.ResponseCode != "00" {
		logger.Make(c, nil).Debug(model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{respBRI.ResponseCode, respBRI.ResponseMessage}))

		return model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{respBRI.ResponseCode,
			respBRI.ResponseMessage})
	}

	return nil
}

func (rul *restUpdateLimits) CorePostInquiryUpdateLimit(c echo.Context, cif string, savingAccNum string, nominalLimit int64) string {
	// reqStatus code
	// EQ Inquiry
	// RQ Request Nasabah
	// AP Approve Dari Bank
	// XX Pembatalan Dari Bank
	const (
		reqStatus         = "EQ"
		isBlokirTrue      = "1"
		isRecalculateTrue = "1"
	)

	r := api.SwitchingResponse{}
	body := map[string]interface{}{
		"isBlokir":         isBlokirTrue,
		"noRek":            savingAccNum,
		"cif":              cif,
		"nominalTransaksi": nominalLimit,
		"isRecalculate":    isRecalculateTrue,
		"reqStatus":        reqStatus,
	}

	req := api.MappingRequestSwitching(body)
	err := api.RetryableSwitchingPost(c, req, "/goldcard/updateLimit/inquiry", &r)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return "99"
	}

	if r.ResponseCode != "00" {
		logger.Make(c, nil).Debug(model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode, r.ResponseDesc}))
	}

	return r.ResponseCode
}
