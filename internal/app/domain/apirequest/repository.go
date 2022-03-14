package apirequest

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// Repository represent the apirequest's repository contract
type Repository interface {
	InserAPIRequest(echo.Context, model.APIRequest) error
}
