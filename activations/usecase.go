package activations

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the activations usecases
type UseCase interface {
	PostActivations(echo.Context, models.PayloadActivations) error
}
