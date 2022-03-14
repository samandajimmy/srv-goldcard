package update_limit

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// UseCase represent the registrations usecases
type UseCase interface {
	DecreasedSTL(c echo.Context, pl model.PayloadCoreDecreasedSTL) model.ResponseErrors
	InquiryUpdateLimit(c echo.Context, pl model.PayloadInquiryUpdateLimit) (model.RespUpdateLimitInquiry, model.ResponseErrors)
	PostUpdateLimit(c echo.Context, pcds model.PayloadUpdateLimit) model.ResponseErrors
	CoreGtePayment(c echo.Context, pl model.PayloadCoreGtePayment) model.ResponseErrors
	GetSavingAccount(c echo.Context, plAcc model.PayloadAccNumber) (interface{}, error)
}
