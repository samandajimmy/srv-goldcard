package main

import (
	"database/sql"
	"fmt"
	"gade/srv-goldcard/middleware"
	"gade/srv-goldcard/models"
	"net/http"
	"os"
	"strconv"
	"time"

	_productreqsHttpsDelivery "gade/srv-goldcard/productreqs/delivery/http"
	_productreqsUseCase "gade/srv-goldcard/productreqs/usecase"
	_registrationsHttpDelivery "gade/srv-goldcard/registrations/delivery/http"
	_registrationsRepository "gade/srv-goldcard/registrations/repository"
	_registrationsUseCase "gade/srv-goldcard/registrations/usecase"
	_tokenHttpDelivery "gade/srv-goldcard/tokens/delivery/http"
	_tokenRepository "gade/srv-goldcard/tokens/repository"
	_tokenUseCase "gade/srv-goldcard/tokens/usecase"

	"github.com/go-pg/pg/v9"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"

	"gade/srv-goldcard/logger"
)

var ech *echo.Echo

func init() {
	ech = echo.New()
	ech.Debug = true
	loadEnv()
	logger.Init()
}

func main() {
	dbConn, dbpg := getDBConn()
	migrate := dataMigrations(dbConn)

	defer dbConn.Close()
	defer migrate.Close()

	echoGroup := models.EchoGroup{
		Admin: ech.Group("/admin"),
		API:   ech.Group("/api"),
		Token: ech.Group("/token"),
	}

	contextTimeout, err := strconv.Atoi(os.Getenv(`CONTEXT_TIMEOUT`))

	if err != nil {
		logger.Make(nil, nil).Debug(err)
	}

	timeoutContext := time.Duration(contextTimeout) * time.Second

	// load all middlewares
	middleware.InitMiddleware(ech, echoGroup)

	// TOKEN
	tokenRepository := _tokenRepository.NewPsqlTokenRepository(dbConn)
	tokenUseCase := _tokenUseCase.NewTokenUseCase(tokenRepository, timeoutContext)
	_tokenHttpDelivery.NewTokensHandler(echoGroup, tokenUseCase)

	// REGISTRATIONS
	registrationsRepository := _registrationsRepository.NewPsqlRegistrationsRepository(dbConn, dbpg)
	registrationsUserCase := _registrationsUseCase.RegistrationsUseCase(registrationsRepository)
	_registrationsHttpDelivery.NewRegistrationsHandler(echoGroup, registrationsUserCase)

	// PRODUCT REQUIREMENTS
	productreqsUseCase := _productreqsUseCase.ProductReqsUseCase()
	_productreqsHttpsDelivery.NewProductreqsHandler(echoGroup, productreqsUseCase)

	// PING
	ech.GET("/ping", ping)

	ech.Start(":" + os.Getenv(`PORT`))

}

func ping(echTx echo.Context) error {
	response := models.Response{}
	response.Status = models.StatusSuccess
	response.Message = "PONG!!"

	return echTx.JSON(http.StatusOK, response)
}

func getDBConn() (*sql.DB, *pg.DB) {
	dbHost := os.Getenv(`DB_HOST`)
	dbPort := os.Getenv(`DB_PORT`)
	dbUser := os.Getenv(`DB_USER`)
	dbPass := os.Getenv(`DB_PASS`)
	dbName := os.Getenv(`DB_NAME`)

	connection := fmt.Sprintf("postgres://%s%s@%s%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	dbConn, err := sql.Open(`postgres`, connection)

	if err != nil {
		logger.Make(nil, nil).Debug(err)
	}

	err = dbConn.Ping()

	if err != nil {
		logger.Make(nil, nil).Debug(err)
		os.Exit(1)
	}

	// go-pg connection initiation
	dbOpt, err := pg.ParseURL(connection)

	if err != nil {
		logger.Make(nil, nil).Debug(err)
	}

	dbpg := pg.Connect(dbOpt)

	if os.Getenv(`DB_LOGGER`) == "true" {
		dbpg.AddQueryHook(logger.DbLogger{})
	}

	return dbConn, dbpg
}

func dataMigrations(dbConn *sql.DB) *migrate.Migrate {
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})

	migrations, err := migrate.NewWithDatabaseInstance(
		"file://migrations/",
		os.Getenv(`DB_USER`), driver)

	if err != nil {
		logger.Make(nil, nil).Debug(err)
	}

	if err := migrations.Up(); err != nil {
		logger.Make(nil, nil).Debug(err)
	}

	return migrations
}

func loadEnv() {
	// check .env file existence
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		return
	}

	err := godotenv.Load()

	if err != nil {
		logger.Make(nil, nil).Fatal("Error loading .env file")
	}

	return
}
