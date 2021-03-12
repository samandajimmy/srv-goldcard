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
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"gopkg.in/gomail.v2"
)

type updateLimitUseCase struct {
	arRepo     activations.RestRepository
	trxRepo    transactions.Repository
	trResRepo  transactions.RestRepository
	trxUS      transactions.UseCase
	rRepo      registrations.Repository
	rrRepo     registrations.RestRepository
	rUS        registrations.UseCase
	upLimRepo  update_limits.Repository
	rupLimRepo update_limits.RestRepository
}

// UpdateLimitUseCase represent Update Limit Use Case
func UpdateLimitUseCase(arRepo activations.RestRepository, trxRepo transactions.Repository, trResRepo transactions.RestRepository,
	trxUS transactions.UseCase, rRepo registrations.Repository, rrRepo registrations.RestRepository, rUS registrations.UseCase,
	upLimRepo update_limits.Repository, rupLimRepo update_limits.RestRepository) update_limits.UseCase {
	return &updateLimitUseCase{arRepo, trxRepo, trResRepo, trxUS, rRepo, rrRepo, rUS, upLimRepo, rupLimRepo}
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
		// TODO: this function need to be changed later on phase 2
		err := upLimUC.rRepo.UpdateCardLimit(c, acc, nil)

		if err != nil {
			continue
		}

		// Send notification to user in pds
		// TODO: this function need to be changed later on phase 2
		notif.GcDecreasedSTL(acc, oldCard, "")
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

func (upLimUC *updateLimitUseCase) InquiryUpdateLimit(c echo.Context, pl models.PayloadInquiryUpdateLimit) (models.RespUpdateLimitInquiry, models.ResponseErrors) {
	var errors models.ResponseErrors
	var response models.RespUpdateLimitInquiry
	var lastLimitUpdate models.LimitUpdate

	// get acc by account number
	acc, err := upLimUC.trxUS.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		errors.SetTitle(models.ErrGetAccByAccountNumber.Error())
		return response, errors
	}

	// check if there is user limit update status still pending or applied
	lastLimitUpdate, err = upLimUC.upLimRepo.GetLastLimitUpdate(c, acc.ID)

	if err != nil {
		errors.SetTitle(models.ErrGetLastLimitUpdate.Error())
		return response, errors
	}

	if lastLimitUpdate.ID != 0 {
		errors.SetTitleCode("12", models.ErrPendingUpdateLimitAvailable.Error(), "")
		return response, errors
	}

	// validate inquiry update limit closed date
	// the closed date is parameterized
	upLimClosedDate, err := upLimUC.upLimRepo.GetUpdateLimitInquiriesClosedDate(c)

	if err != nil {
		errors.SetTitle(models.ErrGetParameter.Error())
		return response, errors
	}

	if strings.Contains(upLimClosedDate, time.Now().Format("02")) {
		errors.SetTitle(models.ErrClosedUpdateLimitInquiries.Error())
		return response, errors
	}

	// syncronize card limit and balance with provider's data first
	acc.Card, err = upLimUC.trxUS.UpdateAndGetCardBalance(c, acc)

	if err != nil {
		errors.SetTitle(models.ErrGetCardBalance.Error())
		return response, errors
	}

	// validate inquiries
	// do minimum increase limit, is npwp required, effective balance, and minimum effective balance validation
	// get npwp document
	npwp, err := upLimUC.rRepo.GetDocumentByApplicationId(acc.ApplicationID, "npwp")

	if err != nil {
		errors.SetTitle(models.ErrGetDocument.Error())
		return response, errors
	}

	// check minimum increase limit 1 million rupiah
	if (pl.NominalLimit - acc.Card.CardLimit) < models.MinIncreaseLimit {
		errors.SetTitle(models.ErrMinimumIncreaseLimit.Error())
		return response, errors
	}

	errStr := upLimUC.rupLimRepo.CorePostInquiryUpdateLimit(c, acc.CIF, acc.Application.SavingAccount, pl.NominalLimit)

	if errStr == "13" {
		errors.SetTitle(models.ErrMinimumGoldSavingEffBal.Error())
		return response, errors
	}

	if errStr != "00" {
		errors.SetTitleCode("14", models.ErrPostInquiryUpdateLimitToCore.Error(), "")
		return response, errors
	}

	// insert new limit update
	// get current STL
	currStl, err := upLimUC.rrRepo.GetCurrentGoldSTL(c)

	if err != nil {
		errors.SetTitle(models.ErrGetCurrSTL.Error())
		return response, errors
	}

	uuid, _ := uuid.NewRandom()
	refId := uuid.String()
	limitUpdt := models.LimitUpdate{
		RefId:     refId,
		AccountID: acc.ID,
		CardLimit: pl.NominalLimit,
		GoldLimit: acc.Card.SetGoldLimit(pl.NominalLimit, currStl),
		StlLimit:  currStl,
		Status:    models.LimitUpdateStatusInquired,
	}

	// check if new inquired card limit is above 50 millions rupiah, then npwp is required
	if npwp[0].FileBase64 == models.DefDocBase64 && pl.NominalLimit > models.LimitFiftyMillions {
		errors.SetTitleCode("11", models.ErrNPWPRequired.Error(), "")
		limitUpdt.WithNpwp = true
	}

	err = upLimUC.upLimRepo.InsertUpdateCardLimit(c, limitUpdt)

	if err != nil {
		errors.SetTitle(models.ErrUpdateCardLimit.Error())
		return response, errors
	}

	response.RefId = refId
	return response, errors
}

// PostUpdateLimit is a func to submit update limit after inquiry update limit
func (upLimUC *updateLimitUseCase) PostUpdateLimit(c echo.Context, pl models.PayloadUpdateLimit) models.ResponseErrors {
	var errors models.ResponseErrors
	var notif models.PdsNotification

	// get limit update with account
	limitUpdt, err := upLimUC.upLimRepo.GetLimitUpdate(c, pl.RefId)

	if err != nil {
		errors.SetTitle(models.ErrUpdateLimitNF.Error())
		return errors
	}

	// set card limit along with gold limit
	acc := limitUpdt.Account
	acc.Card.CardLimit = limitUpdt.CardLimit
	acc.Card.GoldLimit = acc.Card.SetGoldLimit(acc.Card.CardLimit, limitUpdt.StlLimit)
	acc.Card.StlLimit = limitUpdt.StlLimit
	// Get Document (ktp, npwp, selfie, slip_te, and app_form)
	docs, err := upLimUC.rRepo.GetDocumentByApplicationId(acc.ApplicationID, "")

	if err != nil {
		errors.SetTitle(models.ErrGetDocument.Error())
		return errors
	}

	acc.Application.Documents = docs
	limitUpdt.Account = acc

	// insert npwp document if any
	if pl.NpwpImageBase64 != "" {
		acc.Application.SetDocument(models.PayloadPersonalInformation{NpwpImageBase64: pl.NpwpImageBase64})
	}

	// insert updated/latest slip TE and npwp
	err = upLimUC.rUS.GenerateSlipTEDocument(c, &acc)

	if err != nil {
		errors.SetTitle(models.ErrGenerateSlipTE.Error())
		return errors
	}

	// post update limit to core
	err = upLimUC.rupLimRepo.CorePostUpdateLimit(c, acc.Application.SavingAccount, acc.Card, acc.CIF)

	if err != nil {
		errors.SetTitle(models.ErrPostUpdateLimitToCore.Error())
		return errors
	}

	// save updated limit updates data table
	limitUpdt.Status = models.LimitUpdateStatusPending
	err = upLimUC.upLimRepo.UpdateCardLimitData(c, limitUpdt)

	if err != nil {
		errors.SetTitle(models.ErrUpdateCardLimit.Error())
		return errors
	}

	// Send notification to user in pds and email
	notif = models.PdsNotification{}

	notif.GcSla2Days(acc)
	_ = upLimUC.rrRepo.SendNotification(c, notif, "")

	return errors
}

// check if core already pass the payload for endpoint
func (upLimUC *updateLimitUseCase) CoreGtePayment(c echo.Context, pcgp models.PayloadCoreGtePayment) models.ResponseErrors {
	var errors models.ResponseErrors
	gtePayment, err := upLimUC.upLimRepo.GetsertGtePayment(c, pcgp)

	if err != nil {
		errors.SetTitle(err.Error())
		return errors
	}

	acc := gtePayment.Account
	err = func() error {
		if gtePayment.BriUpdated {
			return nil
		}

		// send information to BRI after GTE already paid from core
		return upLimUC.trResRepo.PostPaymentBRI(c, acc, pcgp.NominalTransaction)
	}()

	if err != nil {
		errors.SetTitle(models.ErrPostPaymentBRI.Error())
		return errors
	}

	// update gte payment bri updated status
	gtePayment.BriUpdated = true
	_ = upLimUC.upLimRepo.UpdateGtePayment(c, gtePayment, []string{"bri_updated"})
	err = func() error {
		if gtePayment.PdsNotified {
			return nil
		}

		// update and get card balance by account
		acc.Card, err = upLimUC.trxUS.UpdateAndGetCardBalance(c, acc)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			return err
		}

		// Send notification to user in pds and email
		notif := models.PdsNotification{}
		notif.GtePayment(acc, pcgp)
		return upLimUC.rrRepo.SendNotification(c, notif, "")
	}()

	if err != nil {
		errors.SetTitle(models.ErrSendNotification.Error())
		return errors
	}

	// update gte payment pds notified status
	gtePayment.PdsNotified = true
	_ = upLimUC.upLimRepo.UpdateGtePayment(c, gtePayment, []string{"pds_notified"})

	return errors
}

func (upLimUC *updateLimitUseCase) GetSavingAccount(c echo.Context, plAcc models.PayloadAccNumber) (interface{}, error) {
	acc := models.Account{AccountNumber: plAcc.AccountNumber}
	err := upLimUC.trxRepo.GetAccountByAccountNumber(c, &acc)

	if err != nil {
		return nil, err
	}

	return models.SavingAccount{
		SavingAccount: acc.Application.SavingAccount,
	}, err
}

func (upLimUC *updateLimitUseCase) CheckAccountBySavingAccount(c echo.Context, pl interface{}) (models.Account, error) {
	r := reflect.ValueOf(pl)
	savingAcc := r.FieldByName("SavingAccount")

	// Get Account by saving account
	acc, err := upLimUC.upLimRepo.GetAccountBySavingAccount(c, savingAcc.String())

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return models.Account{}, models.ErrGetAccBySavingAcc
	}

	return acc, nil
}
