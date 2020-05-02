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

	_activationHttpDelivery "gade/srv-goldcard/activations/delivery/http"
	_activationRepository "gade/srv-goldcard/activations/repository"
	_activationUseCase "gade/srv-goldcard/activations/usecase"
	_apiRequestsRepository "gade/srv-goldcard/apirequests/repository"
	_apiRequestsUseCase "gade/srv-goldcard/apirequests/usecase"
	_billingsHttpDelivery "gade/srv-goldcard/billings/delivery/http"
	_billingsRepository "gade/srv-goldcard/billings/repository"
	_billingsUseCase "gade/srv-goldcard/billings/usecase"
	_processHandlerRepository "gade/srv-goldcard/process_handler/repository"
	_processHandlerUseCase "gade/srv-goldcard/process_handler/usecase"
	_productreqsHttpsDelivery "gade/srv-goldcard/productreqs/delivery/http"
	_productreqsUseCase "gade/srv-goldcard/productreqs/usecase"
	_registrationsHttpDelivery "gade/srv-goldcard/registrations/delivery/http"
	_registrationsRepository "gade/srv-goldcard/registrations/repository"
	_registrationsUseCase "gade/srv-goldcard/registrations/usecase"
	_tokenHttpDelivery "gade/srv-goldcard/tokens/delivery/http"
	_tokenRepository "gade/srv-goldcard/tokens/repository"
	_tokenUseCase "gade/srv-goldcard/tokens/usecase"
	_transactionsHttpDelivery "gade/srv-goldcard/transactions/delivery/http"
	_transactionsRepository "gade/srv-goldcard/transactions/repository"
	_transactionsUseCase "gade/srv-goldcard/transactions/usecase"
	_updateLimitHttpDelivery "gade/srv-goldcard/update_limits/delivery/http"
	_updateLimitRepository "gade/srv-goldcard/update_limits/repository"
	_updateLimitUseCase "gade/srv-goldcard/update_limits/usecase"

	_cardsHttpDelivery "gade/srv-goldcard/cards/delivery/http"
	_cardsRepository "gade/srv-goldcard/cards/repository"
	_cardsUseCase "gade/srv-goldcard/cards/usecase"

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
	loc, _ := time.LoadLocation("Asia/Jakarta")
	time.Local = loc
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

	// REPOSITORIES
	tokenRepository := _tokenRepository.NewPsqlTokenRepository(dbConn)
	registrationsRepository := _registrationsRepository.NewPsqlRegistrationsRepository(dbConn, dbpg)
	activationRepository := _activationRepository.NewPsqlActivations(dbConn, dbpg)
	restRegistrationsRepo := _registrationsRepository.NewRestRegistrations(activationRepository)
	restActivationRepository := _activationRepository.NewRestActivations()
	apiRequestsRepository := _apiRequestsRepository.NewPsqlAPIRequestsRepository(dbConn, dbpg)
	transactionsRepository := _transactionsRepository.NewPsqlTransactionsRepository(dbConn, dbpg)
	billingsRepository := _billingsRepository.NewPsqlBillingsRepository(dbConn, dbpg)
	restTransactionsRepo := _transactionsRepository.NewRestTransactions()
	processHandlerRepo := _processHandlerRepository.NewPsqlProcHandlerRepository(dbConn, dbpg)
	updateLimitRepo := _updateLimitRepository.NewPsqlUpdateLimitsRepository(dbConn, dbpg)
	restUpdateLimitRepo := _updateLimitRepository.NewRestUpdateLimits()
	cardsRepository := _cardsRepository.NewPsqlCardsRepository(dbConn, dbpg)
	restCardsRepo := _cardsRepository.NewRestCards()

	// USECASES
	productreqsUseCase := _productreqsUseCase.ProductReqsUseCase()
	processHandlerUseCase := _processHandlerUseCase.ProcessHandUseCase(processHandlerRepo, registrationsRepository)
	tokenUseCase := _tokenUseCase.NewTokenUseCase(tokenRepository, timeoutContext)
	transactionsUseCase := _transactionsUseCase.TransactionsUseCase(transactionsRepository, billingsRepository, restTransactionsRepo, registrationsRepository, restRegistrationsRepo)
	registrationsUseCase := _registrationsUseCase.RegistrationsUseCase(registrationsRepository, restRegistrationsRepo, processHandlerUseCase, transactionsUseCase, restActivationRepository)
	activationUserCase := _activationUseCase.ActivationUseCase(activationRepository, restActivationRepository, registrationsRepository, restRegistrationsRepo, registrationsUseCase, restTransactionsRepo)
	_apiRequestsUseCase.ARUseCase = _apiRequestsUseCase.APIRequestsUseCase(apiRequestsRepository)
	billingsUseCase := _billingsUseCase.BillingsUseCase(billingsRepository, restRegistrationsRepo, transactionsUseCase)
	updateLimitUseCase := _updateLimitUseCase.UpdateLimitUseCase(restActivationRepository, transactionsRepository, restTransactionsRepo, transactionsUseCase, registrationsRepository, restRegistrationsRepo, registrationsUseCase, updateLimitRepo, restUpdateLimitRepo)
	cardsUseCase := _cardsUseCase.CardsUseCase(cardsRepository, restCardsRepo, transactionsUseCase)

	// DELIVERIES
	_productreqsHttpsDelivery.NewProductreqsHandler(echoGroup, productreqsUseCase)
	_tokenHttpDelivery.NewTokensHandler(echoGroup, tokenUseCase)
	_registrationsHttpDelivery.NewRegistrationsHandler(echoGroup, registrationsUseCase)
	_activationHttpDelivery.NewActivationsHandler(echoGroup, activationUserCase)
	_transactionsHttpDelivery.NewTransactionsHandler(echoGroup, transactionsUseCase)
	_billingsHttpDelivery.NewBillingsHandler(echoGroup, billingsUseCase)
	_updateLimitHttpDelivery.NewUpdateLimitHandler(echoGroup, updateLimitUseCase)
	_cardsHttpDelivery.NewCardsHandler(echoGroup, cardsUseCase)

	// PING
	ech.GET("/ping", ping)

	err = ech.Start(":" + os.Getenv(`PORT`))

	if err != nil {
		logger.Make(nil, nil).Fatal(err)
	}

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

	if err != nil {
		logger.Make(nil, nil).Debug(err)
	}

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
}
