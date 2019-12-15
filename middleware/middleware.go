package middleware

import (
	"gade/srv-goldcard/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"
	"os"
	"reflect"
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

	ech.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "requestID=${id}, method=${method}, status=${status}, path=${path}, latency=${latency_human} " +
			"host=${host}, remote_ip=${remote_ip}, user_agent=${user_agent}, error=${error} \n",
	}))

	ech.Use(middleware.Recover())
	cm.cors()
	cm.basicAuth()
	// cm.jwtAuth() // klo gk di tutup gk bisa request.
	cm.customValidation()
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
