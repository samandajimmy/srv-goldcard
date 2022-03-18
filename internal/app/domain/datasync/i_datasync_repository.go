package datasync

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

type IDataSyncRp interface {
	RGetAllAccount(c echo.Context) ([]model.SyncActivation, error)
}
