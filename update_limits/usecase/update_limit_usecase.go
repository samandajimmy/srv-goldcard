package usecase

import (
	"encoding/csv"
	"fmt"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"gade/srv-goldcard/transactions"
	"gade/srv-goldcard/update_limits"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"gopkg.in/gomail.v2"
)

type updateLimitUseCase struct {
	trxRepo transactions.Repository
	rRepo   registrations.Repository
	rrRepo  registrations.RestRepository
	psqlUL  update_limits.Repository
}

// UpdateLimitUseCase represent Update Limit Use Case
func UpdateLimitUseCase(trxRepo transactions.Repository, rRepo registrations.Repository,
	rrRepo registrations.RestRepository, psqlUL update_limits.Repository) update_limits.UseCase {
	return &updateLimitUseCase{trxRepo, rRepo, rrRepo, psqlUL}
}

// DecreasedSTL is a func to recalculate gold card rupiah limit when occurs stl decreased equal or more than 5%
func (ulUS *updateLimitUseCase) DecreasedSTL(c echo.Context, pcds models.PayloadCoreDecreasedSTL) models.ResponseErrors {
	var errors models.ResponseErrors
	var notif models.PdsNotification
	var oldCard models.Card

	// check if payload decreased five percent is false then return
	if pcds.DecreasedFivePercent != "true" {
		return errors
	}

	// Get CurrentStl from Core payload
	currStl := pcds.STL

	// Get All Active Account
	allAccs, err := ulUS.trxRepo.GetAllActiveAccount(c)

	if err != nil {
		errors.SetTitle(models.ErrGetAccByAccountNumber.Error())
		return errors
	}

	for _, acc := range allAccs {
		notif = models.PdsNotification{}
		oldCard = acc.Card

		// set card limit
		err = acc.Card.SetCardLimit(currStl)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			continue
		}

		// update card limit in db
		refId, err := ulUS.rRepo.UpdateCardLimit(c, acc, true)

		if err != nil {
			continue
		}

		// Send notification to user in pds
		notif.GcDecreasedSTL(acc, oldCard, refId)
		_ = ulUS.rrRepo.SendNotification(c, notif, "mobile")

		// Send email notification to user in pds
		ulUS.SendNotificationEmail(c, acc, oldCard)
	}

	return errors
}

func (ulUS *updateLimitUseCase) SendNotificationEmail(c echo.Context, acc models.Account, oldCard models.Card) {
	name := acc.PersonalInformation.FirstName
	ktpNo := acc.PersonalInformation.Nik
	dob := acc.PersonalInformation.BirthDate
	oldLimit := strconv.FormatInt(oldCard.CardLimit, 10)
	newLimit := strconv.FormatInt(acc.Card.CardLimit, 10)

	// create csv file based on data
	file, err := os.Create("./data-stl.csv")
	if err != nil {
		fmt.Println(err)
	}

	writer := csv.NewWriter(file)
	var data = [][]string{
		{"nama", "no_ktp", "tanggal_lahir", "limit_lama", "limit_baru"},
		{name, ktpNo, dob, oldLimit, newLimit}}

	err = writer.WriteAll(data)
	if err != nil {
		logger.Make(c, nil).Debug(err)
	}

	keyParam := "UPDATE_LIMIT_EMAIL_ADDRESS"
	param, err := ulUS.psqlUL.GetParameterByKey(keyParam)
	fmt.Println(param.Value)

	if err != nil {
		logger.Make(nil, nil).Debug(err)
	}

	// call the smtp config from .env
	smtpHost := os.Getenv(`PDS_EMAIL_HOST`)
	smtpPort, _ := strconv.Atoi(os.Getenv(`PDS_EMAIL_PORT`))
	smtpEmail := os.Getenv(`PDS_EMAIL_USERNAME`)
	smtpPass := os.Getenv(`PDS_EMAIL_PASSWORD`)

	// gomail instance to sending an email
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", smtpEmail)
	mailer.SetHeader("To", param.Value)
	mailer.SetHeader("Subject", "Pegadaian Kartu Emas - STL Turun 5%")
	mailer.SetBody("text/plain", "Selamat Pagi \n\n Berikut terlampir file perubahan STL yang turun >5% \n\n Terima Kasih")
	mailer.Attach("./data-stl.csv")

	dialer := gomail.NewDialer(
		smtpHost,
		smtpPort,
		smtpEmail,
		smtpPass,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		logger.Make(c, nil).Debug(err)
	}

	// delete csv file
	os.Remove("./data-stl.csv")
}
