package http

import (
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/transactions"

	"github.com/labstack/echo"
)

// RegistrationsHandler represent the httphandler for registrations
type TransactionsHandler struct {
	response            models.Response
	respErrors          models.ResponseErrors
	transactionsUseCase transactions.UseCase
}

func NewTransactionsHandler(
	echoGroup models.EchoGroup,
	trUseCase transactions.UseCase) {
	handler := &TransactionsHandler{
		transactionsUseCase: trUseCase,
	}

	// End Point For BRI
	echoGroup.API.POST("/transactions/bri", handler.BRIPendingTransactions)
}

// Registrations a handler to handle goldcard registrations
func (th *TransactionsHandler) BRIPendingTransactions(c echo.Context) error {
	var pbpt models.PayloadBRIPendingTransactions
	th.response, th.respErrors = models.NewResponse()

	if err := c.Bind(&pbpt); err != nil {
		th.respErrors.SetTitle(models.MessageUnprocessableEntity)
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	if err := c.Validate(pbpt); err != nil {
		th.respErrors.SetTitle(err.Error())
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	err := th.transactionsUseCase.PostBRIPendingTransactions(c, pbpt)

	if err != nil {
		th.respErrors.SetTitle(err.Error())
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	th.response.SetResponse("", &th.respErrors)

	return th.response.Body(c, err)
}
