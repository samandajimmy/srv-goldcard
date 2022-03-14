package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"srv-goldcard/internal/app/middleware"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"

	_activationHttpDelivery "srv-goldcard/internal/app/domain/activation/delivery/http"
	_activationRepository "srv-goldcard/internal/app/domain/activation/repository"
	_activationUseCase "srv-goldcard/internal/app/domain/activation/usecase"
	_apiRequestsRepository "srv-goldcard/internal/app/domain/apirequest/repository"
	_apiRequestsUseCase "srv-goldcard/internal/app/domain/apirequest/usecase"
	_billingsHttpDelivery "srv-goldcard/internal/app/domain/billing/delivery/http"
	_billingsRepository "srv-goldcard/internal/app/domain/billing/repository"
	_billingsUseCase "srv-goldcard/internal/app/domain/billing/usecase"
	_cardsHttpDelivery "srv-goldcard/internal/app/domain/card/delivery/http"
	_cardsRepository "srv-goldcard/internal/app/domain/card/repository"
	_cardsUseCase "srv-goldcard/internal/app/domain/card/usecase"
	_healthsHttpDelivery "srv-goldcard/internal/app/domain/health/delivery/http"
	_processHandlerRepository "srv-goldcard/internal/app/domain/process_handler/repository"
	_processHandlerUseCase "srv-goldcard/internal/app/domain/process_handler/usecase"
	_productreqsHttpsDelivery "srv-goldcard/internal/app/domain/productreq/delivery/http"
	_productreqsRepository "srv-goldcard/internal/app/domain/productreq/repository"
	_productreqsUseCase "srv-goldcard/internal/app/domain/productreq/usecase"
	_registrationsHttpDelivery "srv-goldcard/internal/app/domain/registration/delivery/http"
	_registrationsRepository "srv-goldcard/internal/app/domain/registration/repository"
	_registrationsUseCase "srv-goldcard/internal/app/domain/registration/usecase"
	_tokenHttpDelivery "srv-goldcard/internal/app/domain/token/delivery/http"
	_tokenRepository "srv-goldcard/internal/app/domain/token/repository"
	_tokenUseCase "srv-goldcard/internal/app/domain/token/usecase"
	_transactionsHttpDelivery "srv-goldcard/internal/app/domain/transaction/delivery/http"
	_transactionsRepository "srv-goldcard/internal/app/domain/transaction/repository"
	_transactionsUseCase "srv-goldcard/internal/app/domain/transaction/usecase"
	_updateLimitHttpDelivery "srv-goldcard/internal/app/domain/update_limit/delivery/http"
	_updateLimitRepository "srv-goldcard/internal/app/domain/update_limit/repository"
	_updateLimitUseCase "srv-goldcard/internal/app/domain/update_limit/usecase"

	"github.com/go-pg/pg/v9"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
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

	echoGroup := model.EchoGroup{
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
	billingsRestRepository := _billingsRepository.NewRestBillings()
	restTransactionsRepo := _transactionsRepository.NewRestTransactions()
	processHandlerRepo := _processHandlerRepository.NewPsqlProcHandlerRepository(dbConn, dbpg)
	updateLimitRepo := _updateLimitRepository.NewPsqlUpdateLimitsRepository(dbConn, dbpg)
	restUpdateLimitRepo := _updateLimitRepository.NewRestUpdateLimits()
	cardsRepository := _cardsRepository.NewPsqlCardsRepository(dbConn, dbpg)
	restCardsRepo := _cardsRepository.NewRestCards()
	productReqsRepo := _productreqsRepository.NewPsqlProductReqsRepository(dbConn, dbpg)

	// USECASES
	productreqsUseCase := _productreqsUseCase.ProductReqsUseCase(productReqsRepo)
	processHandlerUseCase := _processHandlerUseCase.ProcessHandUseCase(processHandlerRepo, registrationsRepository)
	tokenUseCase := _tokenUseCase.NewTokenUseCase(tokenRepository, timeoutContext)
	transactionsUseCase := _transactionsUseCase.TransactionsUseCase(transactionsRepository, billingsRepository, restTransactionsRepo, registrationsRepository, restRegistrationsRepo)
	registrationsUseCase := _registrationsUseCase.RegistrationsUseCase(registrationsRepository, restRegistrationsRepo, processHandlerUseCase, transactionsUseCase, restActivationRepository)
	activationUserCase := _activationUseCase.ActivationUseCase(activationRepository, restActivationRepository, registrationsRepository, restRegistrationsRepo, registrationsUseCase, restTransactionsRepo, cardsRepository)
	_apiRequestsUseCase.ARUseCase = _apiRequestsUseCase.APIRequestsUseCase(apiRequestsRepository)
	billingsUseCase := _billingsUseCase.BillingsUseCase(billingsRepository, billingsRestRepository, restRegistrationsRepo, transactionsUseCase)
	updateLimitUseCase := _updateLimitUseCase.UpdateLimitUseCase(restActivationRepository, transactionsRepository, restTransactionsRepo, transactionsUseCase, registrationsRepository, restRegistrationsRepo, registrationsUseCase, updateLimitRepo, restUpdateLimitRepo)
	cardsUseCase := _cardsUseCase.CardsUseCase(cardsRepository, restCardsRepo, restTransactionsRepo, transactionsUseCase, registrationsRepository)

	// DELIVERIES
	_productreqsHttpsDelivery.NewProductreqsHandler(echoGroup, productreqsUseCase)
	_tokenHttpDelivery.NewTokensHandler(echoGroup, tokenUseCase)
	_registrationsHttpDelivery.NewRegistrationsHandler(echoGroup, registrationsUseCase)
	_activationHttpDelivery.NewActivationsHandler(echoGroup, activationUserCase)
	_transactionsHttpDelivery.NewTransactionsHandler(echoGroup, transactionsUseCase)
	_billingsHttpDelivery.NewBillingsHandler(echoGroup, billingsUseCase)
	_updateLimitHttpDelivery.NewUpdateLimitHandler(echoGroup, updateLimitUseCase)
	_cardsHttpDelivery.NewCardsHandler(echoGroup, cardsUseCase)
	_healthsHttpDelivery.NewHealthsHandler(ech)

	// PING
	ech.GET("/", ping)
	ech.GET("/ping", ping)

	// run refresh all token
	_ = tokenUseCase.RefreshAllToken()
	go registrationsUseCase.RefreshAppTimeoutJob()

	err = ech.Start(":" + os.Getenv(`PORT`))

	if err != nil {
		logger.Make(nil, nil).Fatal(err)
	}

}

func ping(echTx echo.Context) error {
	response := model.Response{}
	response.Status = model.StatusSuccess
	response.Message = "PONG!!"
	response.Data = map[string]interface{}{
		"appVersion": AppVersion,
		"appHash":    BuildHash,
	}

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
		"file://migration/",
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
