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

func (rc *restCards) PostCardReplaceBRI(c echo.Context, pl models.PayloadBriXkey) error {
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey": pl.BriXkey,
	}

	requestPath := "/card/replace"
	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, requestPath, reqBRIBody.RequestData, &respBRI)

	if errBRI != nil {
		logger.Make(c, nil).Debug(errBRI)

		return errBRI
	}

	if respBRI.ResponseCode != "00" {
		return models.ErrReplaceCard
	}

	return nil
}

func (rc *restCards) CoreBlockaCard(c echo.Context, acc models.Account, cardBlock models.CardBlock) error {
	if cardBlock.Description == "" {
		cardBlock.Description = "Card Blocked"
	}

	r := api.SwitchingResponse{}
	body := map[string]interface{}{
		"isBlokir":   "0",
		"noRek":      acc.Application.SavingAccount,
		"cardNumber": acc.Card.CardNumber,
		"reportDesc": cardBlock.Description,
	}

	req := api.MappingRequestSwitching(body)
	err := api.RetryableSwitchingPost(c, req, "/goldcard/protect", &r)

	if err != nil {
		return err
	}

	if r.ResponseCode != "00" {
		return models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode,
			r.ResponseDesc})
	}

	return nil
}

func (rc *restCards) PdsSetNullAppAccNumber(c echo.Context, cif models.PayloadCIF) error {
	resp := api.PdsResponse{}
	reqBody := cif
	endpoint := "/goldcard/set_null_gcnumber"

	err := api.RetryablePdsPost(c, endpoint, reqBody, &resp, echo.MIMEApplicationForm)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}
