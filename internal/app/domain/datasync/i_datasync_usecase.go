package datasync

import (
	"github.com/labstack/echo"
)

type IDataSyncUS interface {
	UGetAllAccount(c echo.Context) error
}
