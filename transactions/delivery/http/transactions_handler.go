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

	// Endpoint For BRI
	echoGroup.API.POST("/transactions/bri", handler.BRIPendingTransactions)

	// Endpoint For PDS
	echoGroup.API.GET("/transactions/history", handler.HistoryTransactions)
	echoGroup.API.POST("/transactions/balance", handler.GetCardBalance)
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

	if err.Title != "" {
		th.response.SetResponse("", &err)

		return th.response.Body(c, nil)
	}

	th.response.SetResponse("", &err)
	return th.response.Body(c, nil)
}

func (th *TransactionsHandler) HistoryTransactions(c echo.Context) error {
	var pht models.PayloadHistoryTransactions
	th.response, th.respErrors = models.NewResponse()

	if err := c.Bind(&pht); err != nil {
		th.respErrors.SetTitle(models.MessageUnprocessableEntity)
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	if err := c.Validate(pht); err != nil {
		th.respErrors.SetTitle(err.Error())
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	result, err := th.transactionsUseCase.GetTransactionsHistory(c, pht)

	if err.Title != "" {
		th.response.SetResponse("", &err)

		return th.response.Body(c, nil)
	}

	th.response.SetResponse(result, &err)
	return th.response.Body(c, nil)
}

// Registrations a handler to handle goldcard get card information
func (th *TransactionsHandler) GetCardBalance(c echo.Context) error {
	var pan models.PayloadAccNumber
	th.response, th.respErrors = models.NewResponse()

	if err := c.Bind(&pan); err != nil {
		th.respErrors.SetTitle(models.MessageUnprocessableEntity)
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	if err := c.Validate(pan); err != nil {
		th.respErrors.SetTitle(err.Error())
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	resp, err := th.transactionsUseCase.GetCardBalance(c, pan)

	if err != nil {
		th.respErrors.SetTitle(err.Error())
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	th.response.SetResponse(resp, &th.respErrors)
	return th.response.Body(c, err)
}
