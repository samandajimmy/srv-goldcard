package middleware

import (
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"os"
	"reflect"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type customMiddleware struct {
	e *echo.Echo
}

var echGroup models.EchoGroup

// InitMiddleware to generate all middleware that app need
func InitMiddleware(ech *echo.Echo, echoGroup models.EchoGroup) {
	cm := &customMiddleware{ech}
	echGroup = echoGroup

	ech.Use(middleware.RequestIDWithConfig(middleware.DefaultRequestIDConfig))
	cm.customLogging()
	cm.customBodyDump()
	ech.Use(middleware.Recover())
	cm.cors()
	cm.basicAuth()
	// cm.jwtAuth() // klo gk di tutup gk bisa request.
	cm.customValidation()
}

func (cm *customMiddleware) customBodyDump() {
	cm.e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Handler: func(c echo.Context, req, resp []byte) {
			reqBody := c.Request()
			reqStr := string(req)
			respStr := string(resp)

			logger.MakeWithoutReportCaller(c, reqStr).Info("Request payload for endpoint " + reqBody.Method + " " + reqBody.URL.String())
			logger.MakeWithoutReportCaller(c, respStr).Info("Response payload for endpoint " + reqBody.Method + " " + reqBody.URL.String())
		},
	}))
}

func (cm *customMiddleware) customLogging() {
	cm.e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logrus.SetReportCaller(false)
			req := c.Request()
			res := c.Response()
			reqID := req.Header.Get(echo.HeaderXRequestID)

			if reqID == "" {
				reqID = res.Header().Get(echo.HeaderXRequestID)
			}

			logrus.WithFields(logrus.Fields{
				"requestID":  reqID,
				"method":     req.Method,
				"status":     res.Status,
				"host":       req.Host,
				"user_agent": req.UserAgent(),
				"uri":        req.URL.String(),
				"ip":         c.RealIP(),
			}).Info("Incoming request")
			return next(c)
		}
	})
}

func (cm *customMiddleware) customValidation() {
	validator := validator.New()
	customValidator := customValidator{}
	validator.RegisterValidation("isRequiredWith", customValidator.isRequiredWith)
	customValidator.validator = validator
	cm.e.Validator = &customValidator
}

func (cm customMiddleware) cors() {
	cm.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"Access-Control-Allow-Origin"},
		AllowMethods: []string{"*"},
	}))
}

func (cm customMiddleware) basicAuth() {
	echGroup.Token.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == os.Getenv(`BASIC_USERNAME`) && password == os.Getenv(`BASIC_PASSWORD`) {
			return true, nil
		}

		return false, nil
	}))
}

func (cm customMiddleware) jwtAuth() {
	echGroup.Admin.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(os.Getenv(`JWT_SECRET`)),
	}))

	echGroup.API.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(os.Getenv(`JWT_SECRET`)),
	}))

}

func (cv *customValidator) isRequiredWith(fl validator.FieldLevel) bool {
	field := fl.Field()
	otherField, _, _ := fl.GetStructFieldOK()

	if otherField.IsValid() && otherField.Interface() != reflect.Zero(otherField.Type()).Interface() {
		if field.IsValid() && field.Interface() == reflect.Zero(field.Type()).Interface() {
			return false
		}
	}

	return true
}
