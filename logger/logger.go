package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"gade/srv-goldcard/models"
	"runtime"
	"strings"

	"github.com/fatih/structs"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type requestLogger struct {
	RequestID string      `json:"requestID"`
	Payload   interface{} `json:"payload,omitempty"`
}

// Init function to make an initial logger
func Init() {
	logrus.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: models.DateTimeFormatMillisecond + "000",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			tmp := strings.Split(f.File, "/")
			filename := tmp[len(tmp)-1]
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	logrus.SetFormatter(formatter)
	logrus.SetLevel(logrus.DebugLevel)
}

var errReqLogger = errors.New("Error during creating a request logger")

// Make is to get a log parameter
func Make(c echo.Context, payload interface{}) *logrus.Entry {
	var rl requestLogger
	logrus.SetReportCaller(true)

	if c == nil {
		return logrus.WithFields(logrus.Fields{})
	}

	rl.RequestID = c.Response().Header().Get(echo.HeaderXRequestID)

	plStr, ok := payload.(string)
	pl, err := json.Marshal(payload)

	if err != nil {
		logrus.WithFields(logrus.Fields{"requestID": rl.RequestID}).Debug(err)

		return logrus.WithFields(logrus.Fields{})
	}

	if payload != nil {
		rl.Payload = string(pl)
	}

	if ok {
		rl.Payload = plStr
	}

	return logrus.WithFields(logrus.Fields{"params": structs.Map(rl)})
}

// MakeWithoutReportCaller to get a log without report caller
func MakeWithoutReportCaller(c echo.Context, payload interface{}) *logrus.Entry {
	log := Make(c, payload)
	logrus.SetReportCaller(false)

	return log
}
