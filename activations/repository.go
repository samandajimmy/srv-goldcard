package activations

import "github.com/labstack/echo"

// Repository represent the activations repository contract
type Repository interface {
}

// RestRepository represent the rest activations repository contract
type RestRepository interface {
	GetDetailGoldUser(c echo.Context, accNumber string) (map[string]string, error)
}
