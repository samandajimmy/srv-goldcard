package http

import (
	billings "gade/srv-goldcard/billings"
	"gade/srv-goldcard/models"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// BillingsHandler represent the httphandler for Product Requirements
type BillingsHandler struct {
	response        models.Response
	respErrors      models.ResponseErrors
	billingsUseCase billings.UseCase
}

// NewBillingsHandler represent to registration gold card
func NewBillingsHandler(echoGroup models.EchoGroup, bl billings.UseCase) {
	handler := &BillingsHandler{
		billingsUseCase: bl,
	}

	// End Point For BRI
	echoGroup.API.POST("/billings/summary/bri", handler.BRIPegadaianBillings)

	// End Point For External
	echoGroup.API.GET("/billings/statements", handler.billingStatement)
}

func (bhn *BillingsHandler) billingStatement(c echo.Context) error {
	var pan models.PayloadAccNumber
	bhn.response, bhn.respErrors = models.NewResponse()

	if err := c.Bind(&pan); err != nil {
		bhn.respErrors.SetTitle(models.MessageUnprocessableEntity)
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

	return c.JSON(getStatusCode(err), bhn.response)
}

// BRIPegadaianBillings a handler to post pegadaian billings from bri
func (bhn *BillingsHandler) BRIPegadaianBillings(c echo.Context) error {
	var pbpb models.PayloadBRIPegadaianBillings
	bhn.response, bhn.respErrors = models.NewResponse()

	if err := c.Bind(&pbpb); err != nil {
		bhn.respErrors.SetTitle(models.MessageUnprocessableEntity)
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

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if strings.Contains(err.Error(), "400") {
		return http.StatusBadRequest
	}

	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusOK
	}
}
