package http

import (
	"gade/srv-goldcard/models"

	"gade/srv-goldcard/cards"

	"github.com/labstack/echo"
)

// CardsHandler represent the httphandler for Cards
type CardsHandler struct {
	response     models.Response
	respErrors   models.ResponseErrors
	cardsUseCase cards.UseCase
}

// NewCardsHandler represent to hancle cards
func NewCardsHandler(echoGroup models.EchoGroup, ca cards.UseCase) {
	handler := &CardsHandler{
		cardsUseCase: ca,
	}

	// End Point For External
	echoGroup.API.POST("/cards/block", handler.CardBlock)
}

func (chn *CardsHandler) CardBlock(c echo.Context) error {
	var pcb models.PayloadCardBlock
	chn.response, chn.respErrors = models.NewResponse()

	if err := c.Bind(&pcb); err != nil {
		chn.respErrors.SetTitle(models.MessageUnprocessableEntity)
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
