package middleware

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"os"
	"reflect"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"

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

var echGroup model.EchoGroup

// InitMiddleware to generate all middleware that app need
func InitMiddleware(ech *echo.Echo, echoGroup model.EchoGroup) {
	cm := &customMiddleware{ech}
	echGroup = echoGroup

	ech.Use(middleware.RequestIDWithConfig(middleware.DefaultRequestIDConfig))
	cm.customLogging()
	cm.customBodyDump()
	ech.Use(middleware.Recover())
	cm.cors()
	cm.basicAuth()
	cm.jwtAuth()
	cm.customValidation()
}

func (cm *customMiddleware) customBodyDump() {
	cm.e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Handler: func(c echo.Context, req, resp []byte) {
			bodyParser(c, &req)
			reqBody := c.Request()

			logger.MakeWithoutReportCaller(c, req).Info("Request payload for endpoint " + reqBody.Method + " " + reqBody.URL.Path)
			logger.MakeWithoutReportCaller(c, resp).Info("Response payload for endpoint " + reqBody.Method + " " + reqBody.URL.Path)
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
	_ = validator.RegisterValidation("isRequiredWith", customValidator.isRequiredWith)
	_ = validator.RegisterValidation("base64", customValidator.base64)
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

// begin custom validator
func (cv *customValidator) isRequiredWith(fl validator.FieldLevel) bool {
	field := fl.Field()
	otherField, _, _, _ := fl.GetStructFieldOK2()

	if otherField.IsValid() && otherField.Interface() != reflect.Zero(otherField.Type()).Interface() {
		if field.IsValid() && field.Interface() == reflect.Zero(field.Type()).Interface() {
			return false
		}
	}

	return true
}

func (cv *customValidator) base64(fl validator.FieldLevel) bool {
	field := fl.Field()

	// if field is nil
	if field.Interface() == reflect.Zero(field.Type()).Interface() {
		return true
	}

	// check field base64 or not, if error then false
	_, err := base64.StdEncoding.DecodeString(field.Interface().(string))

	return err == nil
}

func bodyParser(c echo.Context, pl *[]byte) {
	if string(*pl) == "" {
		rawQuery := c.Request().URL.RawQuery
		m, err := url.ParseQuery(rawQuery)

		if err != nil {
			logger.Make(nil, nil).Fatal(err)
		}

		*pl, err = json.Marshal(m)

		if err != nil {
			logger.Make(nil, nil).Fatal(err)
		}
	}
}
