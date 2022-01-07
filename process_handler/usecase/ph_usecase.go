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

func (ph *processHandUseCase) IsProcessedAppExisted(c echo.Context, acc models.Account) (bool, error) {
	ps := models.ProcessStatus{
		ProcessID: acc.Application.ID,
		TblName:   models.ApplicationTableName,
	}

	psExisting, err := ph.phRepo.GetProcessHandler(ps)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return false, err
	}

	return psExisting.ID != 0, nil
}

func (ph *processHandUseCase) UpsertAppProcess(c echo.Context, acc *models.Account, errStr string) error {
	ps := models.ProcessStatus{
		Reason:      errStr,
		ProcessID:   acc.Application.ID,
		ProcessType: models.FinalAppProcessType,
		TblName:     models.ApplicationTableName,
		Error:       "true",
	}

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
