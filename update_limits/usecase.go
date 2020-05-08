package update_limits

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the registrations usecases
type UseCase interface {
	DecreasedSTL(c echo.Context, pl models.PayloadCoreDecreasedSTL) models.ResponseErrors
	InquiryUpdateLimit(c echo.Context, pl models.PayloadInquiryUpdateLimit) models.ResponseErrors
	PostUpdateLimit(c echo.Context, pcds models.PayloadUpdateLimit) models.ResponseErrors
	CoreGtePayment(c echo.Context, pl models.PayloadCoreGtePayment) models.ResponseErrors
	GetSavingAccount(c echo.Context, plAppNumber models.PayloadAppNumber) (interface{}, error)
}
