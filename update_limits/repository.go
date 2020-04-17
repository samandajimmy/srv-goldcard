package update_limits

import "github.com/labstack/echo"

// Repository represent the transactions Repository
type Repository interface {
	GetEmailByKey(c echo.Context) (string, error)
}
