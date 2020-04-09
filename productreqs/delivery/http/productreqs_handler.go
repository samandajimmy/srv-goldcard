package http

import (
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/productreqs"
	"net/http"

	"github.com/labstack/echo"
)

// ProductreqsHandler represent the httphandler for Product Requirements
type ProductreqsHandler struct {
	response           models.Response
	respErrors         models.ResponseErrors
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
	preq.response, preq.respErrors = models.NewResponse()
	responseData, err := preq.productReqsUseCase.ProductRequirements(c)

	if err != nil {
		preq.respErrors.SetTitle(err.Error())
		preq.response.SetResponse("", &preq.respErrors)

		return c.JSON(http.StatusBadRequest, preq.response)
	}

	preq.response.SetResponse(responseData, &preq.respErrors)

	return preq.response.Body(c, nil)
}
