package usecase

import (
	"encoding/csv"
	"log"
	"os"
	"srv-goldcard/internal/app/domain/datasync"
	"srv-goldcard/internal/app/domain/transaction"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo"
)

type datasyncUsecase struct {
	dsRepo  datasync.IDataSyncRp
	trxRest transaction.RestRepository
}

func NewDatasyncUsecase(dsRepo datasync.IDataSyncRp, trxRest transaction.RestRepository) datasync.IDataSyncUS {
	return &datasyncUsecase{dsRepo, trxRest}
}

func (dsUs *datasyncUsecase) tempFn(c echo.Context, record model.SyncActivation, idx int) []string {
	var cardInfoErrMsg, appStatusErrMsg string
	var appStatus model.BRIAppStatus

	acc := model.Account{BrixKey: record.Brixkey}
	cardInfo, err := dsUs.trxRest.GetBRICardInformation(c, acc)

	if err != nil {
		cardInfoErrMsg = strings.Replace(err.Error(), "\n", ";", -1)
	}

	appStatus, err = dsUs.trxRest.GetBRIAppStatus(c, record.Brixkey)

	if err != nil {
		appStatusErrMsg = strings.Replace(err.Error(), "\n", ";", -1)
	}

	// if strings.Contains(appStatusErrMsg, "RC- -") {
	// 	appStatusErrMsg +=
	// }

	row := []string{record.Brixkey, record.Cif, record.Name, record.EncryptedCardNumber,
		record.ActivatedDate.Format(model.DateTimeFormatZone), record.AccountStatus, record.CardStatus,
		record.ApplicationStatus, appStatus.AppStatus, cardInfo.BillKey, cardInfo.BlockCode,
		cardInfoErrMsg, appStatusErrMsg}

	return row
}

func (dsUs *datasyncUsecase) tempLoop(c echo.Context, accounts []model.SyncActivation) [][]string {
	// Using WriteAll
	var data [][]string
	var wg sync.WaitGroup
	var batchLimit = 5
	var batchCounter int

	data = append(data, []string{"brixkey", "cif", "name", "billKey", "activatedDate", "accountStatus",
		"cardStatus", "appStatus", "briAppStatus", "briBillKey", "briBlockCode", "cardInfoErrMsg",
		"appStatusErrMsg"})

	for idx, record := range accounts {
		// Increment the WaitGroup counter.
		wg.Add(1)
		batchCounter++
		// Launch a goroutine to fetch the URL.
		go func(record model.SyncActivation, idx int) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			row := dsUs.tempFn(c, record, idx)
			data = append(data, row)
		}(record, idx)

		if batchCounter == batchLimit || idx == len(accounts)-1 {
			// Wait for all HTTP fetches to complete.
			wg.Wait()
			batchCounter = 0
			// delay
			time.Sleep(500 * time.Millisecond)
		}
	}

	return data
}

func (dsUs *datasyncUsecase) UGetAllAccount(c echo.Context) error {
	accounts, _ := dsUs.dsRepo.RGetAllAccount(c)
	file, err := os.Create("application.csv")

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	logger.Dump("cacing", len(accounts))
	// return nil

	w := csv.NewWriter(file)
	defer w.Flush()
	defer file.Close()

	data := dsUs.tempLoop(c, accounts)

	w.WriteAll(data)

	return nil
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
