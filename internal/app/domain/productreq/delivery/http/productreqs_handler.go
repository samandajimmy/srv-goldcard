package http

import (
	"net/http"
	"srv-goldcard/internal/app/domain/productreq"
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// ProductreqsHandler represent the httphandler for Product Requirements
type ProductreqsHandler struct {
	response           model.Response
	respErrors         model.ResponseErrors
	productReqsUseCase productreq.UseCase
}

// NewProductreqsHandler represent to registration gold card
func NewProductreqsHandler(echoGroup model.EchoGroup, pu productreq.UseCase) {
	handler := &ProductreqsHandler{
		productReqsUseCase: pu,
	}

	// End Point For External
	echoGroup.API.GET("/product/requirements", handler.productRequirements)
	echoGroup.API.POST("/product/public-holiday", handler.InsertPublicHolidayDate)
	echoGroup.API.GET("/product/public-holiday", handler.GetPublicHolidayDate)
}

func (preq *ProductreqsHandler) productRequirements(c echo.Context) error {
	preq.response, preq.respErrors = model.NewResponse()
	responseData, err := preq.productReqsUseCase.ProductRequirements(c)

	if err != nil {
		preq.respErrors.SetTitle(err.Error())
		preq.response.SetResponse("", &preq.respErrors)

		return c.JSON(http.StatusBadRequest, preq.response)
	}

	preq.response.SetResponse(responseData, &preq.respErrors)

	return preq.response.Body(c, nil)
}

// InsertPublicHolidayDate a handler to handle public holiday insert
func (preq *ProductreqsHandler) InsertPublicHolidayDate(c echo.Context) error {
	var phd model.PayloadInsertPublicHoliday
	preq.response, preq.respErrors = model.NewResponse()

	if err := c.Bind(&phd); err != nil {
		preq.respErrors.SetTitle(model.MessageUnprocessableEntity)
		preq.response.SetResponse("", &preq.respErrors)

		return preq.response.Body(c, err)
	}

	if err := c.Validate(phd); err != nil {
		preq.respErrors.SetTitle(err.Error())
		preq.response.SetResponse("", &preq.respErrors)

		return preq.response.Body(c, err)
	}

	resp, err := preq.productReqsUseCase.InsertPublicHolidayDate(c, phd)

	if err != nil {
		preq.respErrors.SetTitle(err.Error())
		preq.response.SetResponse("", &preq.respErrors)

		return preq.response.Body(c, err)
	}

	preq.response.SetResponse(resp, &preq.respErrors)

	return preq.response.Body(c, err)
}

// GetPublicHolidayDate a handler to handle public holiday get
func (preq *ProductreqsHandler) GetPublicHolidayDate(c echo.Context) error {
	preq.response, preq.respErrors = model.NewResponse()

	resp, err := preq.productReqsUseCase.GetPublicHolidayDate(c)

	if err != nil {
		preq.respErrors.SetTitle(err.Error())
		preq.response.SetResponse("", &preq.respErrors)

		return preq.response.Body(c, err)
	}

	preq.response.SetResponse(resp, &preq.respErrors)

	return preq.response.Body(c, err)
}
