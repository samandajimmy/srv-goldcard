package apirequests

import (
	"github.com/labstack/echo"
)

// UseCase represent the apirequest's usecases
type UseCase interface {
	PostAPIRequest(c echo.Context, statusCode int, api, req, resp interface{}) error
}
