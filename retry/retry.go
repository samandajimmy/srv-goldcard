package retry

import (
	"fmt"
	"gade/srv-goldcard/logger"
	"os"
	"strconv"

	"github.com/avast/retry-go"
	"github.com/labstack/echo"
)

// Do is function to initiate retry helper
func Do(c echo.Context, fnName string, fn retry.RetryableFunc) error {
	msg := fmt.Sprintf("try to attempt %s", fnName)
	logger.MakeWithoutReportCaller(c, nil).Debug(msg)
	err := retry.Do(fn, attempt(), onRetry(c, fnName))

	if err != nil {
		return err
	}

	logger.MakeWithoutReportCaller(c, nil).Debug(fnName + " attempted successfully")

	return nil
}

func attempt() retry.Option {
	attempt, err := strconv.Atoi(os.Getenv(`RETRY_ON_ERR`))

	if err != nil {
		logger.Make(nil, nil).Fatal(err)
	}

	return retry.Attempts(uint(attempt))
}

func onRetry(c echo.Context, fnName string) retry.Option {
	return retry.OnRetry(func(n uint, err error) {
		msg := fmt.Sprintf("has been attempted %s #%d time with error: %s", fnName, n+1, err)
		logger.MakeWithoutReportCaller(c, nil).Debug(msg)
	})
}
