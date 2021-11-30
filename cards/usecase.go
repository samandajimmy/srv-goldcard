package cards

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

type UseCase interface {
	BlockCard(c echo.Context, pl models.PayloadCardBlock) error
	GetCardStatus(c echo.Context, pl models.PayloadAccNumber) (models.RespCardStatus, error)
	CloseCard(c echo.Context, pl models.PayloadCIF) error
}
