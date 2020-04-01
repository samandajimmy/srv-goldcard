package process_handler

import "gade/srv-goldcard/models"

// Repository represent the process handler repository contract
type Repository interface {
	PostProcessHandler(ps models.ProcessStatus) error
	GetProcessHandler(ps models.ProcessStatus) (bool, error)
}
