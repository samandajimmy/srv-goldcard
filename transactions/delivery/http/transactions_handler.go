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

	// Endpoint For third party
	echoGroup.API.POST("/transactions/payment/:source", handler.paymentTransaction)

	// Endpoint For PDS
	echoGroup.API.GET("/transactions/history", handler.HistoryTransactions)
	echoGroup.API.GET("/transactions/balance", handler.GetCardBalance)
	echoGroup.API.POST("/transactions/payment/inquiry", handler.PaymentInquiry)
	echoGroup.API.POST("/transactions/payment/core", handler.paymentTransactionCore)
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
	var plListTrx models.PayloadListTrx
	th.response, th.respErrors = models.NewResponse()

	if err := c.Bind(&plListTrx); err != nil {
		th.respErrors.SetTitle(models.MessageUnprocessableEntity)
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	if err := c.Validate(plListTrx); err != nil {
		th.respErrors.SetTitle(err.Error())
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	result, err := th.transactionsUseCase.GetTransactionsHistory(c, plListTrx)

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

func (th *TransactionsHandler) paymentTransaction(c echo.Context) error {
	var pbpt models.PayloadPaymentTransactions
	th.response, th.respErrors = models.NewResponse()

	if err := c.Bind(&pbpt); err != nil {
		th.respErrors.SetTitle(models.MessageUnprocessableEntity)
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	// init source variable
	pbpt.Source = c.Param("source")

	if err := c.Validate(pbpt); err != nil {
		th.respErrors.SetTitle(err.Error())
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	err := th.transactionsUseCase.PostPaymentTransaction(c, pbpt)

	if err.Title != "" {
		th.response.SetResponse("", &err)

		return th.response.Body(c, nil)
	}

	th.response.SetResponse("", &err)
	return th.response.Body(c, nil)
}

func (th *TransactionsHandler) paymentTransactionCore(c echo.Context) error {
	var pl models.PlPaymentTrxCore
	th.response, th.respErrors = models.NewResponse()

	if err := c.Bind(&pl); err != nil {
		th.respErrors.SetTitle(models.MessageUnprocessableEntity)
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	if err := c.Validate(pl); err != nil {
		th.respErrors.SetTitle(err.Error())
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	err := th.transactionsUseCase.PostPaymentTrxCore(c, pl)

	if err.Title != "" {
		th.response.SetResponse("", &err)

		return th.response.Body(c, nil)
	}

	th.response.SetResponse("", &err)
	return th.response.Body(c, nil)
}

func (th *TransactionsHandler) PaymentInquiry(c echo.Context) error {
	var ppi models.PlPaymentInquiry
	th.response, th.respErrors = models.NewResponse()

	if err := c.Bind(&ppi); err != nil {
		th.respErrors.SetTitle(models.MessageUnprocessableEntity)
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	if err := c.Validate(ppi); err != nil {
		th.respErrors.SetTitle(err.Error())
		th.response.SetResponse("", &th.respErrors)

		return th.response.Body(c, err)
	}

	refTrx, err := th.transactionsUseCase.PaymentInquiry(c, ppi)

	if err.Title != "" {
		th.response.SetResponse("", &err)

		return th.response.Body(c, nil)
	}

	response := map[string]string{"refTrx": refTrx}
	th.response.SetResponse(response, &err)
	return th.response.Body(c, nil)
}
