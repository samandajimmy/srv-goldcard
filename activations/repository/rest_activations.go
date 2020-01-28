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

// GetDetailGoldUser to get detail gold user from core
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

	if r.ResponseCode != "00" {
		logger.Make(c, nil).Debug(models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode, r.ResponseDesc}))

		return nilMap, err
	}

	return r.ResponseData, nil
}
