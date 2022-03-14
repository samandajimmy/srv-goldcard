package process_handler

import "srv-goldcard/internal/app/model"

// Repository represent the process handler repository contract
type Repository interface {
	PostProcessHandler(ps model.ProcessStatus) error
	GetProcessHandler(ps model.ProcessStatus) (model.ProcessStatus, error)
	UpdateProcessHandler(ps model.ProcessStatus, col []string) error
}
