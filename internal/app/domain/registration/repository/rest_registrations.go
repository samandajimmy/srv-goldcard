package repository

import (
	"fmt"
	"srv-goldcard/internal/app/domain/activation"
	"srv-goldcard/internal/app/domain/registration"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/api"
	"srv-goldcard/internal/pkg/logger"
	"strconv"

	"github.com/labstack/echo"
)

type restRegistrations struct {
	aRepo activation.Repository
}

// NewRestRegistrations will create an object that represent the registration.Repository interface
func NewRestRegistrations(aRepo activation.Repository) registration.RestRepository {
	return &restRegistrations{aRepo}
}

// GetCurrentGoldSTL to get current STL from core
func (rr *restRegistrations) GetCurrentGoldSTL(c echo.Context) (int64, error) {
	// get stored gold price
	hargaEmas, err := rr.aRepo.GetStoredGoldPrice(c)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		hargaEmas = 0
	}

	r := api.SwitchingResponse{}
	STLBody := map[string]interface{}{}
	req := api.MappingRequestSwitching(STLBody)
	err = api.RetryableSwitchingPost(c, req, "/param/stl", &r)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return hargaEmas, nil
	}

	if r.ResponseCode != "00" {
		logger.Make(c, nil).Debug(model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode, r.ResponseDesc}))

		return hargaEmas, nil
	}

	if _, ok := r.ResponseData["hargaEmas"].(string); !ok {
		logger.Make(c, nil).Debug(model.ErrSetVar)

		return 0, model.ErrSetVar
	}

	currSTL, err := strconv.ParseInt(r.ResponseData["hargaEmas"].(string), 10, 64)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return 0, err
	}

	return currSTL, nil
}

func (rr *restRegistrations) OpenGoldcard(c echo.Context, acc model.Account, isRecalculate bool) error {
	const (
		isBlokirTrue      = "1"
		isRecalculateTrue = "1"
	)

	r := api.SwitchingResponse{}
	body := map[string]interface{}{
		"isBlokir":         isBlokirTrue,
		"noRek":            acc.Application.SavingAccount,
		"gramTransaksi":    fmt.Sprintf("%f", acc.Card.GoldLimit),
		"nominalTransaksi": strconv.FormatInt(acc.Card.CardLimit, 10),
	}

	if isRecalculate {
		body["isRecalculate"] = isRecalculateTrue
	}

	req := api.MappingRequestSwitching(body)
	err := api.RetryableSwitchingPost(c, req, "/goldcard/open", &r)

	if err != nil {
		return err
	}

	if r.ResponseCode != "00" {
		return model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode,
			r.ResponseDesc})
	}

	return nil
}

func (rr *restRegistrations) CloseGoldcard(c echo.Context, acc model.Account) error {
	const isBlokirFalse = "0"
	r := api.SwitchingResponse{}
	body := map[string]interface{}{
		"isBlokir":      isBlokirFalse,
		"noRek":         acc.Application.SavingAccount,
		"gramTransaksi": fmt.Sprintf("%f", acc.Card.GoldLimit),
	}

	req := api.MappingRequestSwitching(body)
	err := api.RetryableSwitchingPost(c, req, "/goldcard/close", &r)

	if err != nil {
		return err
	}

	if r.ResponseCode != api.APIRCSuccess {
		return model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode,
			r.ResponseDesc})
	}

	return nil
}

func (rr *restRegistrations) SendNotification(c echo.Context, notif model.PdsNotification, notifType string) error {
	resp := api.PdsResponse{}
	reqBody := notif
	endpoint := "/notification/send"

	switch notifType {
	case "email":
		endpoint += "/email"
	case "mobile":
		endpoint += "/mobile"
	default:
		endpoint += ""
	}

	err := api.RetryablePdsPost(c, endpoint, reqBody, &resp, echo.MIMEApplicationJSON)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (rr *restRegistrations) AuthLogin() error {
	resp := api.PdsResponse{}
	reqBody := map[string]string{
		"email":    "082141217929",
		"password": "gadai123",
		"agen":     "android",
		"version":  "3",
	}
	endpoint := "/auth/login/new"
	err := api.RetryablePdsPost(nil, endpoint, reqBody, &resp, echo.MIMEApplicationForm)

	if err != nil {
		logger.Make(nil, nil).Debug(err)

		return err
	}

	return nil
}
