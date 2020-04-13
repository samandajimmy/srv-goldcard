package update_limits

import (
	"gade/srv-goldcard/models"
)

// Repository represent the transactions Repository
type Repository interface {
	GetParameterByKey(key string) (models.Parameter, error)
}
