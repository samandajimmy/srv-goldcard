package process_handler

import "github.com/labstack/echo"

// UseCase represent the process handler usecases
type UseCase interface {
	ProcHandFinalApp(c echo.Context, applicationNumber, processID, processType, status string, errStatus bool)
	PostProcessHandler(c echo.Context, process_id, process_type, status string) error
	StatProcessCheck(c echo.Context, processID, status string) (bool, error)
}
