package http

import (
	"srv-goldcard/internal/app/domain/datasync"
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

type datasyncHandler struct {
	response   model.Response
	respErrors model.ResponseErrors
	dsUs       datasync.IDataSyncUS
}

func NewDatasyncHandler(echoGroup model.EchoGroup, dsUs datasync.IDataSyncUS) {
	handler := &datasyncHandler{dsUs: dsUs}

	echoGroup.API.GET("/data-sync/activation", handler.hActivation)
}

func (dsh *datasyncHandler) hActivation(c echo.Context) error {
	dsh.response, dsh.respErrors = model.NewResponse()
	err := dsh.dsUs.UGetAllAccount(c)

	if err != nil {
		dsh.respErrors.SetTitle(err.Error())
		dsh.response.SetResponse("", &dsh.respErrors)

		return dsh.response.Body(c, err)
	}

	dsh.response.SetResponse([]model.Account{}, &dsh.respErrors)
	return dsh.response.Body(c, err)
}
