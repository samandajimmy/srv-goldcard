package usecase

import (
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/process_handler"
	"gade/srv-goldcard/registrations"

	"github.com/labstack/echo"
)

type processHandUseCase struct {
	phRepo  process_handler.Repository
	regRepo registrations.Repository
}

// ProcessHandUseCase represent Process Handler Use Case
func ProcessHandUseCase(phRepo process_handler.Repository, regRepo registrations.Repository) process_handler.UseCase {
	return &processHandUseCase{
		phRepo:  phRepo,
		regRepo: regRepo,
	}
}

// PostProcessHandler represent Usecase push process handler
func (ph *processHandUseCase) PostProcessHandler(c echo.Context, ps models.ProcessStatus) error {
	n, err := ph.checkProcessStatuses(c, ps)

	if err != nil {
		return err
	}

	if !n {
		return nil
	}

	// Insert Process Handler
	err = ph.phRepo.PostProcessHandler(ps)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}

	return nil
}

func (ph *processHandUseCase) checkProcessStatuses(c echo.Context, ps models.ProcessStatus) (bool, error) {
	insert, err := ph.phRepo.GetProcessHandler(ps)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return insert, err
	}

	return insert, nil
}
