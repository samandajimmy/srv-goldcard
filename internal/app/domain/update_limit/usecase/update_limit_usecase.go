package usecase

import (
	"encoding/csv"
	"os"
	"reflect"
	"srv-goldcard/internal/app/domain/activation"
	"srv-goldcard/internal/app/domain/registration"
	"srv-goldcard/internal/app/domain/transaction"
	update_limit "srv-goldcard/internal/app/domain/update_limit"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"gopkg.in/gomail.v2"
)

type updateLimitUseCase struct {
	arRepo     activation.RestRepository
	trxRepo    transaction.Repository
	trResRepo  transaction.RestRepository
	trxUS      transaction.UseCase
	rRepo      registration.Repository
	rrRepo     registration.RestRepository
	rUS        registration.UseCase
	upLimRepo  update_limit.Repository
	rupLimRepo update_limit.RestRepository
}

// UpdateLimitUseCase represent Update Limit Use Case
func UpdateLimitUseCase(arRepo activation.RestRepository, trxRepo transaction.Repository, trResRepo transaction.RestRepository,
	trxUS transaction.UseCase, rRepo registration.Repository, rrRepo registration.RestRepository, rUS registration.UseCase,
	upLimRepo update_limit.Repository, rupLimRepo update_limit.RestRepository) update_limit.UseCase {
	return &updateLimitUseCase{arRepo, trxRepo, trResRepo, trxUS, rRepo, rrRepo, rUS, upLimRepo, rupLimRepo}
}

// DecreasedSTL is a func to recalculate gold card rupiah limit when occurs stl decreased equal or more than 5%
func (upLimUC *updateLimitUseCase) DecreasedSTL(c echo.Context, pcds model.PayloadCoreDecreasedSTL) model.ResponseErrors {
	var errors model.ResponseErrors
	var notif model.PdsNotification
	var oldCard model.Card
	var cul []model.CardUpdateLimit

	// check if payload decreased five percent is false then return
	if pcds.DecreasedFivePercent != "true" {
		return errors
	}

	// Get CurrentStl from Core payload
	currStl := pcds.STL

	// Get All Active Account
	allAccs, err := upLimUC.trxRepo.GetAllActiveAccount(c)

	if err != nil {
		errors.SetTitle(model.ErrGetAccByAccountNumber.Error())
		return errors
	}

	for _, acc := range allAccs {
		notif = model.PdsNotification{}
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
		cul = append(cul, model.CardUpdateLimit{OldCard: oldCard, NewCard: acc.Card, Account: acc})
	}

	// Send an notification with its attachment to email
	upLimUC.SendNotificationEmail(c, cul)

	return errors
}

// function to send list of decreased STL to email
func (upLimUC *updateLimitUseCase) SendNotificationEmail(c echo.Context, cul []model.CardUpdateLimit) {

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

func (upLimUC *updateLimitUseCase) InquiryUpdateLimit(c echo.Context, pl model.PayloadInquiryUpdateLimit) (model.RespUpdateLimitInquiry, model.ResponseErrors) {
	var errors model.ResponseErrors
	var response model.RespUpdateLimitInquiry
	var lastLimitUpdate model.LimitUpdate

	// get acc by account number
	acc, err := upLimUC.trxUS.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		errors.SetTitle(model.ErrGetAccByAccountNumber.Error())
		return response, errors
	}

	// check if there is user limit update status still pending or applied
	lastLimitUpdate, err = upLimUC.upLimRepo.GetLastLimitUpdate(c, acc.ID)

	if err != nil {
		errors.SetTitle(model.ErrGetLastLimitUpdate.Error())
		return response, errors
	}

	if lastLimitUpdate.ID != 0 {
		errors.SetTitleCode("12", model.ErrPendingUpdateLimitAvailable.Error(), "")
		return response, errors
	}

	// validate inquiry update limit closed date
	// the closed date is parameterized
	upLimClosedDate, err := upLimUC.upLimRepo.GetUpdateLimitInquiriesClosedDate(c)

	if err != nil {
		errors.SetTitle(model.ErrGetParameter.Error())
		return response, errors
	}

	if strings.Contains(upLimClosedDate, time.Now().Format("02")) {
		errors.SetTitle(model.ErrClosedUpdateLimitInquiries.Error())
		return response, errors
	}

	// syncronize card limit and balance with provider's data first
	acc.Card, err = upLimUC.trxUS.UpdateAndGetCardBalance(c, acc)

	if err != nil {
		errors.SetTitle(model.ErrGetCardBalance.Error())
		return response, errors
	}

	// validate inquiries
	// do minimum increase limit, is npwp required, effective balance, and minimum effective balance validation
	// get npwp document
	npwp, err := upLimUC.rRepo.GetDocumentByApplicationId(acc.ApplicationID, "npwp")

	if err != nil {
		errors.SetTitle(model.ErrGetDocument.Error())
		return response, errors
	}

	// check minimum increase limit 1 million rupiah
	if (pl.NominalLimit - acc.Card.CardLimit) < model.MinIncreaseLimit {
		errors.SetTitle(model.ErrMinimumIncreaseLimit.Error())
		return response, errors
	}

	errStr := upLimUC.rupLimRepo.CorePostInquiryUpdateLimit(c, acc.CIF, acc.Application.SavingAccount, pl.NominalLimit)

	if errStr == "13" {
		errors.SetTitle(model.ErrMinimumGoldSavingEffBal.Error())
		return response, errors
	}

	if errStr != "00" {
		errors.SetTitleCode("14", model.ErrPostInquiryUpdateLimitToCore.Error(), "")
		return response, errors
	}

	// insert new limit update
	// get current STL
	currStl, err := upLimUC.rrRepo.GetCurrentGoldSTL(c)

	if err != nil {
		errors.SetTitle(model.ErrGetCurrSTL.Error())
		return response, errors
	}

	uuid, _ := uuid.NewRandom()
	refId := uuid.String()
	limitUpdt := model.LimitUpdate{
		RefId:     refId,
		AccountID: acc.ID,
		CardLimit: pl.NominalLimit,
		GoldLimit: acc.Card.SetGoldLimit(pl.NominalLimit, currStl),
		StlLimit:  currStl,
		Status:    model.LimitUpdateStatusInquired,
	}

	// check if new inquired card limit is above 50 millions rupiah, then npwp is required
	if npwp[0].FileBase64 == model.DefDocBase64 && pl.NominalLimit > model.LimitFiftyMillions {
		errors.SetTitleCode("11", model.ErrNPWPRequired.Error(), "")
		limitUpdt.WithNpwp = true
	}

	err = upLimUC.upLimRepo.InsertUpdateCardLimit(c, limitUpdt)

	if err != nil {
		errors.SetTitle(model.ErrUpdateCardLimit.Error())
		return response, errors
	}

	response.RefId = refId
	return response, errors
}

// PostUpdateLimit is a func to submit update limit after inquiry update limit
func (upLimUC *updateLimitUseCase) PostUpdateLimit(c echo.Context, pl model.PayloadUpdateLimit) model.ResponseErrors {
	var errors model.ResponseErrors
	var notif model.PdsNotification

	// get limit update with account
	limitUpdt, err := upLimUC.upLimRepo.GetLimitUpdate(c, pl.RefId)

	if err != nil {
		errors.SetTitle(model.ErrUpdateLimitNF.Error())
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
		errors.SetTitle(model.ErrGetDocument.Error())
		return errors
	}

	acc.Application.Documents = docs
	limitUpdt.Account = acc

	// insert npwp document if any
	if pl.NpwpImageBase64 != "" {
		acc.Application.SetDocument(model.PayloadPersonalInformation{NpwpImageBase64: pl.NpwpImageBase64})
	}

	// insert updated/latest slip TE and npwp
	err = upLimUC.rUS.GenerateSlipTEDocument(c, &acc)

	if err != nil {
		errors.SetTitle(model.ErrGenerateSlipTE.Error())
		return errors
	}

	// post update limit to core
	err = upLimUC.rupLimRepo.CorePostUpdateLimit(c, acc.Application.SavingAccount, acc.Card, acc.CIF)

	if err != nil {
		errors.SetTitle(model.ErrPostUpdateLimitToCore.Error())
		return errors
	}

	// save updated limit updates data table
	limitUpdt.Status = model.LimitUpdateStatusPending
	err = upLimUC.upLimRepo.UpdateCardLimitData(c, limitUpdt)

	if err != nil {
		errors.SetTitle(model.ErrUpdateCardLimit.Error())
		return errors
	}

	// Send notification to user in pds and email
	notif = model.PdsNotification{}

	notif.GcSla2Days(acc)
	_ = upLimUC.rrRepo.SendNotification(c, notif, "")

	return errors
}

// check if core already pass the payload for endpoint
func (upLimUC *updateLimitUseCase) CoreGtePayment(c echo.Context, pcgp model.PayloadCoreGtePayment) model.ResponseErrors {
	var errors model.ResponseErrors
	gtePayment, err := upLimUC.upLimRepo.GetsertGtePayment(c, pcgp)

	// if account is closed or payment notif already succeess then set response code to 22
	if err == model.ErrGtePaymenTrxIdExist || err == model.ErrSavingAccNotFound {
		errors.SetTitleCode("22", model.ErrGtePaymenTrxIdExistOrAccountClosed.Error(), model.ErrGtePaymenTrxIdExistOrAccountClosed.Error())
		return errors
	}

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
		errors.SetTitle(model.ErrPostPaymentBRI.Error())
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
		notif := model.PdsNotification{}
		notif.GtePayment(acc, pcgp)
		return upLimUC.rrRepo.SendNotification(c, notif, "")
	}()

	if err != nil {
		errors.SetTitle(model.ErrSendNotification.Error())
		return errors
	}

	// update gte payment pds notified status
	gtePayment.PdsNotified = true
	_ = upLimUC.upLimRepo.UpdateGtePayment(c, gtePayment, []string{"pds_notified"})

	return errors
}

func (upLimUC *updateLimitUseCase) GetSavingAccount(c echo.Context, plAcc model.PayloadAccNumber) (interface{}, error) {
	acc := model.Account{AccountNumber: plAcc.AccountNumber}
	err := upLimUC.trxRepo.GetAccountByAccountNumber(c, &acc)

	if err != nil {
		return nil, err
	}

	return model.SavingAccount{
		SavingAccount: acc.Application.SavingAccount,
	}, err
}

func (upLimUC *updateLimitUseCase) CheckAccountBySavingAccount(c echo.Context, pl interface{}) (model.Account, error) {
	r := reflect.ValueOf(pl)
	savingAcc := r.FieldByName("SavingAccount")

	// Get Account by saving account
	acc, err := upLimUC.upLimRepo.GetAccountBySavingAccount(c, savingAcc.String())

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return model.Account{}, model.ErrGetAccBySavingAcc
	}

	return acc, nil
}
