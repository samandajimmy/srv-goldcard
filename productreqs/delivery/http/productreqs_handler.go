package http

import (
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/productreqs"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

var response models.Response

// ProductreqsHandler represent the httphandler for Product Requirements
type ProductreqsHandler struct {
	productReqsUseCase productreqs.UseCase
}

// NewProductreqsHandler represent to registration gold card
func NewProductreqsHandler(echoGroup models.EchoGroup, pu productreqs.UseCase) {
	handler := &ProductreqsHandler{
		productReqsUseCase: pu,
	}

	// End Point For External
	echoGroup.API.GET("/product/requirements", handler.productRequirements)
}

func (preq *ProductreqsHandler) productRequirements(c echo.Context) error {
	response = models.Response{}
	respErrors := &models.ResponseErrors{}

	responseData, err := preq.productReqsUseCase.ProductRequirements(c)

	if err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)

		return c.JSON(http.StatusBadRequest, response)
	}

	response.SetResponse(responseData, respErrors)

	return c.JSON(getStatusCode(err), response)
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
