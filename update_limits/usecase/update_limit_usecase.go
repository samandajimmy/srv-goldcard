package usecase

import (
	"encoding/csv"
	"gade/srv-goldcard/activations"
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
	arRepo    activations.RestRepository
	trxRepo   transactions.Repository
	trxUS     transactions.UseCase
	rRepo     registrations.Repository
	rrRepo    registrations.RestRepository
	upLimRepo update_limits.Repository
}

// UpdateLimitUseCase represent Update Limit Use Case
func UpdateLimitUseCase(arRepo activations.RestRepository, trxRepo transactions.Repository,
	trxUS transactions.UseCase, rRepo registrations.Repository, rrRepo registrations.RestRepository,
	upLimRepo update_limits.Repository) update_limits.UseCase {
	return &updateLimitUseCase{arRepo, trxRepo, trxUS, rRepo, rrRepo, upLimRepo}
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
		logger.Make(c, nil).Debug(err)
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

func (upLimUC *updateLimitUseCase) InquiryUpdateLimit(c echo.Context, pl models.PayloadInquiryUpdateLimit) models.ResponseErrors {
	var errors models.ResponseErrors

	// get acc by account number
	acc, err := upLimUC.trxUS.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		errors.SetTitle(models.ErrGetAccByAccountNumber.Error())
		return errors
	}

	// get user gold effective balance
	userGoldDetail, err := upLimUC.arRepo.GetDetailGoldUser(c, acc.Application.SavingAccount)

	if err != nil {
		errors.SetTitle(models.ErrGetUserDetail.Error())
		return errors
	}

	if _, ok := userGoldDetail["saldoEfektif"].(string); !ok {
		errors.SetTitle(models.ErrSetVar.Error())
		return errors
	}

	goldEffBalance, err := strconv.ParseFloat(userGoldDetail["saldoEfektif"].(string), 64)

	if err != nil {
		errors.SetTitle(models.ErrGetEffBalance.Error())
		return errors
	}

	// get current STL first
	currStl, err := upLimUC.rrRepo.GetCurrentGoldSTL(c)

	if err != nil {
		errors.SetTitle(models.ErrGetCurrSTL.Error())
		return errors
	}

	// check minimum increase limit 1 million rupiah
	if pl.NominalLimit-acc.Card.CardLimit < models.MinIncreaseLimit {
		errors.SetTitle(models.ErrMinimumIncreaseLimit.Error())
		return errors
	}

	// check if gold effective balance is sufficient
	err = upLimUC.checkGoldEffBalanceSufficient(pl.NominalLimit, acc.Card, currStl, goldEffBalance)
	if err != nil {
		errors.SetTitle(err.Error())
		return errors
	}

	// check if new inquired card limit is above 50 millions rupiah, then npwp is required
	docs, err := upLimUC.rRepo.GetDocumentByApplicationId(acc.ApplicationID)
	if err != nil {
		errors.SetTitle(models.ErrGetDocument.Error())
		return errors
	}

	npwp := acc.Application.GetCurrentDoc(docs, models.MapDocType["NpwpImageBase64"])
	if npwp.ID == 0 && pl.NominalLimit > models.LimitFiftyMillions {
		errors.SetTitleCode("11", models.ErrNPWPRequired.Error(), "")
		return errors
	}

	return errors
}

// checkGoldEffBalanceSufficient is a function to check whether remaining effective gold balance is sufficient when trying to increase card limit
func (upLimUC *updateLimitUseCase) checkGoldEffBalanceSufficient(newLimit int64, currentCard models.Card, currStl int64, goldEffBalance float64) error {
	appliedGoldLimit := currentCard.GoldLimit
	newGoldLimit := currentCard.SetGoldLimit(newLimit, currStl)
	deficitGoldLimit := models.CustomRound("round", newGoldLimit-appliedGoldLimit, 10000)

	// got not enough effective gold balance
	if goldEffBalance < deficitGoldLimit {
		return models.ErrInsufGoldSavingEffBalance
	}

	// got not enough minimum effective balance 0.1 gram
	if goldEffBalance < deficitGoldLimit+models.MinEffBalance {
		return models.ErrMinimumGoldSavingEffBal
	}

	return nil
}
