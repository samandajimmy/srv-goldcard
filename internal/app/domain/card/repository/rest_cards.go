package repository

import (
	"encoding/json"
	"srv-goldcard/internal/app/domain/card"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/api"
	"srv-goldcard/internal/pkg/logger"
	"time"

	"github.com/labstack/echo"
)

type restCards struct{}

// NewRestCards will create an object that represent the card.RestRepository interface
func NewRestCards() card.RestRepository {
	return &restCards{}
}

func (rc *restCards) GetBRICardBlockStatus(c echo.Context, acc model.Account, pl model.PayloadCardBlock) (model.BRICardBlockStatus, error) {
	var briCardBlockStatus model.BRICardBlockStatus
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey": acc.BrixKey,
	}

	// Set request path based on reason code
	requestPath := model.RequestPathCardBlock
	if pl.ReasonCode == model.ReasonCodeStolen {
		requestPath = model.RequestPathCardStolen
	}

	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, requestPath, reqBRIBody.RequestData, &respBRI)

	// RC 73 means that cards already blocked in BRI, so we may skip this process
	if respBRI.ResponseCode == "73" {
		briCardBlockStatus.ReportDesc = "Card already blocked in BRI, trying to block core"
		briCardBlockStatus.ReportingDate = time.Now().Unix() * 1000
		return briCardBlockStatus, nil
	}

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

func (rc *restCards) PostCardReplaceBRI(c echo.Context, pl model.PayloadBriXkey) error {
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
		return model.ErrReplaceCard
	}

	return nil
}

func (rc *restCards) CoreBlockaCard(c echo.Context, acc model.Account, cardBlock model.CardBlock) error {
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
		return model.DynamicErr(model.ErrSwitchingAPIRequest, []interface{}{r.ResponseCode,
			r.ResponseDesc})
	}

	return nil
}

func (rc *restCards) PdsSetNullAppAccNumber(c echo.Context, cif model.PayloadCIF) error {
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
