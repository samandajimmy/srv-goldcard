package usecase

import (
	"srv-goldcard/internal/app/domain/process_handler"
	"srv-goldcard/internal/app/domain/registration"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"

	"github.com/labstack/echo"
)

type processHandUseCase struct {
	phRepo  process_handler.Repository
	regRepo registration.Repository
}

// ProcessHandUseCase represent Process Handler Use Case
func ProcessHandUseCase(phRepo process_handler.Repository, regRepo registration.Repository) process_handler.UseCase {
	return &processHandUseCase{
		phRepo:  phRepo,
		regRepo: regRepo,
	}
}

func (ph *processHandUseCase) IsProcessedAppExisted(c echo.Context, acc model.Account) (bool, error) {
	ps := model.ProcessStatus{
		ProcessID: acc.Application.ID,
		TblName:   model.ApplicationTableName,
	}

	psExisting, err := ph.phRepo.GetProcessHandler(ps)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return false, err
	}

	return psExisting.ID != 0, nil
}

func (ph *processHandUseCase) UpsertAppProcess(c echo.Context, acc *model.Account, errStr string) error {
	ps := model.ProcessStatus{
		Reason:      errStr,
		ProcessID:   acc.Application.ID,
		ProcessType: model.FinalAppProcessType,
		TblName:     model.ApplicationTableName,
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
func (ph *processHandUseCase) UpdateErrorStatus(c echo.Context, acc model.Account) error {
	var ps model.ProcessStatus
	err := ps.MapUpdateProcessStatus(model.ApplicationTableName, acc.Application.ID)

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

func (ph *processHandUseCase) updateReasonProcStatus(c echo.Context, ps model.ProcessStatus) error {
	col := []string{"reason", "updated_at", "error", "error_count"}

	err := ph.phRepo.UpdateProcessHandler(ps, col)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}

	return nil
}
