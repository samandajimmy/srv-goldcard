package card

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

type UseCase interface {
	BlockCard(c echo.Context, pl model.PayloadCardBlock) error
	GetCardStatus(c echo.Context, pl model.PayloadAccNumber) (model.RespCardStatus, error)
	CloseCard(c echo.Context, pl model.PayloadCIF) error
}
