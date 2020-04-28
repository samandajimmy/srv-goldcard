package repository

import (
	"encoding/json"
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/cards"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

type restCards struct{}

// NewRestCards will create an object that represent the cards.RestRepository interface
func NewRestCards() cards.RestRepository {
	return &restCards{}
}

func (rc *restCards) GetBRICardBlockStatus(c echo.Context, acc models.Account, pl models.PayloadCardBlock) (models.BRICardBlockStatus, error) {
	var briCardBlockStatus models.BRICardBlockStatus
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey": acc.BrixKey,
	}

	// Set request path based on reason code
	requestPath := models.RequestPathCardBlock
	if pl.ReasonCode == models.ReasonCodeStolen {
		requestPath = models.RequestPathCardStolen
	}

	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, requestPath, reqBRIBody.RequestData, &respBRI)

	if errBRI != nil {
		logger.Make(c, nil).Debug(errBRI)

		return briCardBlockStatus, errBRI
	}

	mrshlCardInfo, err := json.Marshal(respBRI.DataOne)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return briCardBlockStatus, err
	}

	err = json.Unmarshal(mrshlCardInfo, &briCardBlockStatus)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return briCardBlockStatus, err
	}

	return briCardBlockStatus, nil
}
