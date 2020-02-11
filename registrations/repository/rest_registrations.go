package repository

import (
	"gade/srv-goldcard/activations"
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"strconv"

	"github.com/labstack/echo"
)

type restRegistrations struct {
	aRepo activations.Repository
}

// NewRestRegistrations will create an object that represent the registrations.Repository interface
func NewRestRegistrations(aRepo activations.Repository) registrations.RestRepository {
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
		logger.Make(c, nil).Debug(models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode, r.ResponseDesc}))

		return hargaEmas, nil
	}

	currSTL, err := strconv.ParseInt(r.ResponseData["hargaEmas"], 10, 64)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return 0, err
	}

	return currSTL, nil
}

func (rr *restRegistrations) OpenGoldcard(c echo.Context, acc models.Account, isRecalculate bool) error {
	const (
		isBlokirTrue      = "1"
		isRecalculateTrue = "1"
	)

	r := api.SwitchingResponse{}
	body := map[string]interface{}{
		"isBlokir":      isBlokirTrue,
		"noRek":         acc.Application.SavingAccount,
		"gramTransaksi": acc.Card.CardLimit,
		"stlPengajuan":  acc.Card.CurrentSTL,
	}

	if isRecalculate {
		body["isRecalculate"] = isRecalculateTrue
	}

	req := api.MappingRequestSwitching(body)
	err := api.RetryableSwitchingPost(c, req, "/goldcard/open", &r)

	if err != nil {
		return err
	}

	return nil
}
