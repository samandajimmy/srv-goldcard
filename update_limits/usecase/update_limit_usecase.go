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
	trxRepo   transactions.Repository
	rRepo     registrations.Repository
	rrRepo    registrations.RestRepository
	upLimRepo update_limits.Repository
}

// UpdateLimitUseCase represent Update Limit Use Case
func UpdateLimitUseCase(trxRepo transactions.Repository, rRepo registrations.Repository,
	rrRepo registrations.RestRepository, upLimRepo update_limits.Repository) update_limits.UseCase {
	return &updateLimitUseCase{trxRepo, rRepo, rrRepo, upLimRepo}
}

// DecreasedSTL is a func to recalculate gold card rupiah limit when occurs stl decreased equal or more than 5%
func (upLimUC *updateLimitUseCase) DecreasedSTL(c echo.Context, pcds models.PayloadCoreDecreasedSTL) models.ResponseErrors {
	var errors models.ResponseErrors
	var notif models.PdsNotification
	var oldCard models.Card
	var cul []models.CardUpdateLimit

	// check if payload decreased five percent is false then return
	if pcds.DecreasedFivePercent != "true" {
		return errors
	}

	// Get CurrentStl from Core payload
	currStl := pcds.STL

	// Get All Active Account
	allAccs, err := upLimUC.trxRepo.GetAllActiveAccount(c)

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
		refId, err := upLimUC.rRepo.UpdateCardLimit(c, acc, true)

		if err != nil {
			continue
		}

		// Send notification to user in pds
		notif.GcDecreasedSTL(acc, oldCard, refId)
		_ = upLimUC.rrRepo.SendNotification(c, notif, "mobile")

		// Insert all STL data that changes to cul struct
		cul = append(cul, models.CardUpdateLimit{OldCard: oldCard, NewCard: acc.Card, Account: acc})
	}

	// Send an notification with its attachment to email
	upLimUC.SendNotificationEmail(c, cul)

	return errors
}

// function to send list of decreased STL to email
func (upLimUC *updateLimitUseCase) SendNotificationEmail(c echo.Context, cul []models.CardUpdateLimit) {

	var data [][]string

	// append all cul data to 2D array
	for _, val := range cul {
		data = append(data, [][]string{{
			val.Account.PersonalInformation.FirstName,
			val.Account.PersonalInformation.Nik,
			val.Account.PersonalInformation.BirthDate,
			strconv.FormatInt(val.OldCard.CardLimit, 10),
			strconv.FormatInt(val.OldCard.CardLimit, 10),
		}}...)
	}

	// create csv file based on data
	file, err := os.Create("./data-stl.csv")
	if err != nil {
		fmt.Println(err)
	}

	writer := csv.NewWriter(file)

	err = writer.WriteAll(data)
	if err != nil {
		logger.Make(c, nil).Debug(err)
	}

	emailAddres, err := upLimUC.upLimRepo.GetEmailByKey(c)

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
	mailer.SetHeader("To", emailAddres)
	mailer.SetHeader("Subject", "Pegadaian Kartu Emas - STL Turun 5%")
	mailer.SetBody("text/plain", "Selamat Pagi \n\nBerikut terlampir file perubahan STL yang turun >5% \n\nTerima Kasih")
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
