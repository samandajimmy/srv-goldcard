package usecase

import (
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/process_handler"
	"gade/srv-goldcard/registrations"

	"github.com/google/uuid"
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

func (ph *processHandUseCase) ProcHandFinalApp(c echo.Context, applicationNumber, processID, processType, status string, errStatus bool) {
	if processID == "" {
		uuid, _ := uuid.NewRandom()
		processID = uuid.String()
	}

	err := ph.PostProcessHandler(c, processID, processType, status)

	if err != nil {
		logger.Make(c, nil).Debug(err)
	}
	go func() {
		err = ph.regRepo.UpdateAppError(c, applicationNumber, processID, errStatus)

		if err != nil {
			logger.Make(c, nil).Debug(err)
		}
	}()
}

func (ph *processHandUseCase) PostProcessHandler(c echo.Context, processID, processType, status string) error {
	var ps models.ProcessStatus
	ps, err := ph.phRepo.GetProcessHandler(processID)

	if err != nil {
		return err
	}

	if ps.ProcessID == "" {
		// Map insert process status
		err := ps.MapInsertProcessStatus(processID, processType, status)

		if err != nil {
			return err
		}

		// Insert Process Handler
		err = ph.phRepo.PostProcessHandler(ps)

		if err != nil {
			return err
		}

		return nil
	}

	// Map update process status
	err = ps.MapUpdateProcessStatus(status)

	if err != nil {
		return nil
	}

	// Update Process
	err = ph.phRepo.PutProcessHandler(ps)

	if err != nil {
		return err
	}

	return nil
}
