package http

import (
	"srv-goldcard/internal/app/domain/card"
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// CardsHandler represent the httphandler for Cards
type CardsHandler struct {
	response     model.Response
	respErrors   model.ResponseErrors
	cardsUseCase card.UseCase
}

// NewCardsHandler represent to hancle cards
func NewCardsHandler(echoGroup model.EchoGroup, ca card.UseCase) {
	handler := &CardsHandler{
		cardsUseCase: ca,
	}

	// End Point For External
	echoGroup.API.POST("/cards/block", handler.CardBlock)
	echoGroup.API.GET("/cards/status", handler.CardStatus)
	echoGroup.API.POST("/cards/close", handler.CardClose)
}

func (chn *CardsHandler) CardBlock(c echo.Context) error {
	var pcb model.PayloadCardBlock
	chn.response, chn.respErrors = model.NewResponse()

	if err := c.Bind(&pcb); err != nil {
		chn.respErrors.SetTitle(model.MessageUnprocessableEntity)
		chn.response.SetResponse("", &chn.respErrors)

		return chn.response.Body(c, err)
	}

	if err := c.Validate(pcb); err != nil {
		chn.respErrors.SetTitle(err.Error())
		chn.response.SetResponse("", &chn.respErrors)

		return chn.response.Body(c, err)
	}

	err := chn.cardsUseCase.BlockCard(c, pcb)

	if err != nil {
		chn.respErrors.SetTitle(err.Error())
		chn.response.SetResponse("", &chn.respErrors)

		return chn.response.Body(c, err)
	}

	chn.response.SetResponse("", &chn.respErrors)
	return chn.response.Body(c, err)
}

func (chn *CardsHandler) CardStatus(c echo.Context) error {
	var pl model.PayloadAccNumber
	chn.response, chn.respErrors = model.NewResponse()

	if err := c.Bind(&pl); err != nil {
		chn.respErrors.SetTitle(model.MessageUnprocessableEntity)
		chn.response.SetResponse("", &chn.respErrors)

		return chn.response.Body(c, err)
	}

	if err := c.Validate(pl); err != nil {
		chn.respErrors.SetTitle(err.Error())
		chn.response.SetResponse("", &chn.respErrors)

		return chn.response.Body(c, err)
	}

	resp, err := chn.cardsUseCase.GetCardStatus(c, pl)

	if err != nil {
		chn.respErrors.SetTitle(err.Error())
		chn.response.SetResponse("", &chn.respErrors)

		return chn.response.Body(c, err)
	}

	chn.response.SetResponse(resp, &chn.respErrors)
	return chn.response.Body(c, err)
}

func (chn *CardsHandler) CardClose(c echo.Context) error {
	var pl model.PayloadCIF
	chn.response, chn.respErrors = model.NewResponse()

	if err := c.Bind(&pl); err != nil {
		chn.respErrors.SetTitle(model.MessageUnprocessableEntity)
		chn.response.SetResponse("", &chn.respErrors)

		return chn.response.Body(c, err)
	}

	if err := c.Validate(pl); err != nil {
		chn.respErrors.SetTitle(err.Error())
		chn.response.SetResponse("", &chn.respErrors)

		return chn.response.Body(c, err)
	}

	err := chn.cardsUseCase.CloseCard(c, pl)

	if err != nil {
		chn.respErrors.SetTitle(err.Error())
		chn.response.SetResponse("", &chn.respErrors)

		return chn.response.Body(c, err)
	}

	chn.response.SetResponse("", &chn.respErrors)
	return chn.response.Body(c, err)
}
