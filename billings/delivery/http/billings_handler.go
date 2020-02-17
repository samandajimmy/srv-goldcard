package http

import (
	billings "gade/srv-goldcard/billings"
	"gade/srv-goldcard/models"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

var response models.Response

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

	// End Point For External
	echoGroup.API.GET("/billings/billing-statement", handler.billingStatement)
}

func (bh *BillingsHandler) billingStatement(c echo.Context) error {
	var pan models.PayloadAccNumber
	// response = models.Response{}
	// respErrors := &models.ResponseErrors{}

	bh.response, bh.respErrors = models.NewResponse()

	if err := c.Bind(&pan); err != nil {
		bh.respErrors.SetTitle(models.MessageUnprocessableEntity)
		bh.response.SetResponse("", &bh.respErrors)

		return bh.response.Body(c, err)
	}

	if err := c.Validate(pan); err != nil {
		bh.respErrors.SetTitle(err.Error())
		bh.response.SetResponse("", &bh.respErrors)

		return bh.response.Body(c, err)
	}

	responseData, err := bh.billingsUseCase.GetBillingStatement(c, pan)

	if err != nil {
		bh.respErrors.SetTitle(err.Error())
		bh.response.SetResponse("", &bh.respErrors)

		return c.JSON(http.StatusBadRequest, bh.response)
	}

	bh.response.SetResponse(responseData, &bh.respErrors)

	return c.JSON(getStatusCode(err), bh.response)
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
