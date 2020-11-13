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
	res, err := ph.phRepo.GetProcessHandler(ps)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}

	if res.ID != 0 {
		// Update Process Handler reason
		res.Reason = res.Reason + "||" + ps.Reason
		res.Error = ps.Error
		res.ErrorCount += 1
		_ = ph.updateReasonProcStatus(c, res)

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

// UpdateCounterError Method Update counter error on table process_statuses
func (ph *processHandUseCase) UpdateCounterError(c echo.Context, acc models.Account) {
	var ps models.ProcessStatus
	err := ps.MapUpdateProcessStatus(models.ApplicationTableName, acc.Application.ID)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return
	}

	psOld, err := ph.phRepo.GetProcessHandler(ps)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return
	}

	col := []string{"updated_at", "error_count"}

	ps.ErrorCount = psOld.ErrorCount + 1
	err = ph.phRepo.UpdateProcessHandler(ps, col)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return
	}
}

// UpdateErrorStatus Method Update error status on table process_statuses
func (ph *processHandUseCase) UpdateErrorStatus(c echo.Context, acc models.Account) error {
	var ps models.ProcessStatus
	err := ps.MapUpdateProcessStatus(models.ApplicationTableName, acc.Application.ID)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}

	col := []string{"updated_at", "error"}

	ps.Error = "false"
	err = ph.phRepo.UpdateProcessHandler(ps, col)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}

	return nil
}

func (ph *processHandUseCase) updateReasonProcStatus(c echo.Context, ps models.ProcessStatus) error {
	col := []string{"reason", "updated_at", "error", "error_count"}

	err := ph.phRepo.UpdateProcessHandler(ps, col)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}

	return nil
}
