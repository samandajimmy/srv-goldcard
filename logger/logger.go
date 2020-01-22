package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/fatih/structs"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

const timestampFormat = "2006-01-02 15:04:05.000"

var (
	errReqLogger      = errors.New("Error during creating a request logger")
	errEchoContextNil = errors.New("Echo context tidak boleh nil")
)

type requestLogger struct {
	RequestID string      `json:"requestID,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
}

// Init function to make an initial logger
func Init() {
	logrus.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: timestampFormat,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			tmp := strings.Split(f.File, "/")
			filename := tmp[len(tmp)-1]
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	logrus.SetFormatter(formatter)
	logrus.SetLevel(logrus.DebugLevel)
}

// Make is to get a log parameter
func Make(c echo.Context, payload interface{}) *logrus.Entry {
	var rl requestLogger
	logrus.SetReportCaller(true)

	if c != nil {
		rl.RequestID = c.Response().Header().Get(echo.HeaderXRequestID)
	}

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

// GetEchoRID to get echo request ID
func GetEchoRID(c echo.Context) string {
	if c == nil {
		Make(nil, nil).Fatal(errEchoContextNil)
	}

	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// MakeStructToJSON to get a json string of struct
// JUST FOR DEBUGGING TOOL
func MakeStructToJSON(strct interface{}) {
	b, err := json.Marshal(strct)

	if err != nil {
		Make(nil, nil).Fatal(err)

		return
	}

	fmt.Println()
	MakeWithoutReportCaller(nil, nil).Debug(string(b))
	fmt.Println()
}
