package http

import (
	"net/http"
	"srv-goldcard/internal/app/domain/billing"
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// BillingsHandler represent the httphandler for Product Requirements
type BillingsHandler struct {
	response        model.Response
	respErrors      model.ResponseErrors
	billingsUseCase billing.UseCase
}

// NewBillingsHandler represent to registration gold card
func NewBillingsHandler(echoGroup model.EchoGroup, bl billing.UseCase) {
	handler := &BillingsHandler{
		billingsUseCase: bl,
	}

	// End Point For BRI
	echoGroup.API.POST("/billings/summary/bri", handler.BRIPegadaianBillings)

	// End Point For External
	echoGroup.API.GET("/billings/statements", handler.billingStatement)
}

func (bhn *BillingsHandler) billingStatement(c echo.Context) error {
	var pan model.PayloadAccNumber
	bhn.response, bhn.respErrors = model.NewResponse()

	if err := c.Bind(&pan); err != nil {
		bhn.respErrors.SetTitle(model.MessageUnprocessableEntity)
		bhn.response.SetResponse("", &bhn.respErrors)

		return bhn.response.Body(c, err)
	}

	if err := c.Validate(pan); err != nil {
		bhn.respErrors.SetTitle(err.Error())
		bhn.response.SetResponse("", &bhn.respErrors)

		return bhn.response.Body(c, err)
	}

	responseData, err := bhn.billingsUseCase.GetBillingStatement(c, pan)

	if err != nil {
		bhn.respErrors.SetTitle(err.Error())
		bhn.response.SetResponse("", &bhn.respErrors)

		return c.JSON(http.StatusBadRequest, bhn.response)
	}

	bhn.response.SetResponse(responseData, &bhn.respErrors)

	if (responseData == model.BillingStatement{}) {
		bhn.response.SetResponse("", &bhn.respErrors)
	}

	return bhn.response.Body(c, nil)
}

// BRIPegadaianBillings a handler to post pegadaian billings from bri
func (bhn *BillingsHandler) BRIPegadaianBillings(c echo.Context) error {
	var pbpb model.PayloadBRIPegadaianBillings
	bhn.response, bhn.respErrors = model.NewResponse()

	if err := c.Bind(&pbpb); err != nil {
		bhn.respErrors.SetTitle(model.MessageUnprocessableEntity)
		bhn.response.SetResponse("", &bhn.respErrors)

		return bhn.response.Body(c, err)
	}

	if err := c.Validate(pbpb); err != nil {
		bhn.respErrors.SetTitle(err.Error())
		bhn.response.SetResponse("", &bhn.respErrors)

		return bhn.response.Body(c, err)
	}

	err := bhn.billingsUseCase.PostBRIPegadaianBillings(c, pbpb)

	if err.Title != "" {
		bhn.response.SetResponse("", &err)

		return bhn.response.Body(c, nil)
	}

	bhn.response.SetResponse("", &err)
	return bhn.response.Body(c, nil)
}
