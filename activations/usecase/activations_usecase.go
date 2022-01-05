package usecase

import (
	"errors"
	"gade/srv-goldcard/activations"
	"gade/srv-goldcard/cards"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"gade/srv-goldcard/transactions"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

type activationsUseCase struct {
	aRepo    activations.Repository
	arRepo   activations.RestRepository
	rRepo    registrations.Repository
	rrRepo   registrations.RestRepository
	rUsecase registrations.UseCase
	trRepo   transactions.RestRepository
	cardRepo cards.Repository
}

// ActivationUseCase represent Activation Use Case
func ActivationUseCase(aRepo activations.Repository, arRepo activations.RestRepository,
	rRepo registrations.Repository, rrRepo registrations.RestRepository, rUsecase registrations.UseCase,
	trRepo transactions.RestRepository, cardRepo cards.Repository) activations.UseCase {
	return &activationsUseCase{aRepo, arRepo, rRepo, rrRepo, rUsecase, trRepo, cardRepo}
}

func (aUsecase *activationsUseCase) InquiryActivation(c echo.Context, acc models.Account) (models.CardBalance, models.ResponseErrors) {
	var errors models.ResponseErrors
	var cardBal models.CardBalance
	var err error

	// get account and check app number
	if acc.ID == int64(0) {
		appNumber := models.PayloadAppNumber{
			ApplicationNumber: acc.Application.ApplicationNumber,
		}

		acc, err = aUsecase.rUsecase.CheckApplication(c, appNumber)
	}

	if err != nil {
		errors.SetTitle(err.Error())
		return cardBal, errors
	}

	// if account or card has been activated
	if acc.Status == models.AccStatusActive &&
		acc.Application.Status == models.AppStatusActive &&
		acc.Card.Status == models.CardStatusActive {
		errors.SetTitle(models.ErrCardActivated.Error())
		return cardBal, errors
	}

	// if card is being replaced skip activation validation process
	// the condition if there is some record on Card Status table with corresponding card_id then the user has been activated before and request for card replacement
	if aUsecase.isActivationForReplacedCard(acc.Card) {
		return cardBal, errors
	}

	// validation on inquiry
	// validate application expiry from application_processed_date < 12 months
	// add a year for expiry date
	expDate := acc.Application.ApplicationProcessedDate.AddDate(1, 0, 0)
	now := models.NowUTC()

	if now.After(expDate) {
		errors.SetTitleCode("22", models.ErrAppExpired.Error(), models.ErrAppExpiredDesc.Error())
		return cardBal, errors
	}

	// validate stl price changes
	// compare stl price at applied date and current date
	currStl, deficitStl, isDecreased, err := aUsecase.isStlDecreased(c, acc)

	if err != nil {
		errors.SetTitle(models.ErrGetCurrSTL.Error())
		return cardBal, errors
	}

	// return response if there's no decreasing
	if !isDecreased {
		return cardBal, errors
	}

	// get decreasing percentage
	_ = models.CustomRound("round", float64(deficitStl)/float64(currStl), 10000)
	// get user effective balance
	userDetail, err := aUsecase.arRepo.GetDetailGoldUser(c, acc.Application.SavingAccount)

	if err != nil {
		errors.SetTitle(models.ErrGetUserDetail.Error())
		return cardBal, errors
	}

	if _, ok := userDetail["saldoEfektif"].(string); !ok {
		logger.Make(c, nil).Debug(models.ErrSetVar)
		errors.SetTitle(models.ErrSetVar.Error())
		return cardBal, errors
	}

	goldEffBalance, err := strconv.ParseFloat(userDetail["saldoEfektif"].(string), 64)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		errors.SetTitle(models.ErrGetEffBalance.Error())

		return cardBal, errors
	}

	appliedGoldLimit := acc.Card.GoldLimit
	currGoldLimit := acc.Card.SetGoldLimit(acc.Card.CardLimit, currStl)
	// because we need user to have at least 0.1 effective gold balance
	deficitGoldLimit := models.CustomRound("round", currGoldLimit-appliedGoldLimit, 10000)
	cardBal.CurrGoldLimit = currGoldLimit
	cardBal.CurrStl = currStl
	cardBal.DeficitGoldLimit = deficitGoldLimit
	cardBal.SavingAccount = acc.Application.SavingAccount

	// got not enough effective gold balance
	if goldEffBalance < deficitGoldLimit+models.MinEffBalance {
		errors.SetTitleCode("55", models.ErrDecreasedSTL.Error(), models.ErrDecreasedSTLDesc.Error())
		return cardBal, errors
	}

	// got enough effective gold balance
	errors.SetTitleCode(
		"44",
		models.ErrDecreasedSTL.Error(),
		models.DynamicErr(models.ErrDecreasedSTLOpenDesc, []interface{}{deficitGoldLimit}).Error(),
	)

	return cardBal, errors
}

func (aUsecase *activationsUseCase) PostActivations(c echo.Context, pa models.PayloadActivations) (models.RespActivations, error) {
	var respActNil models.RespActivations
	acc, err := aUsecase.rUsecase.CheckApplication(c, pa)

	if err != nil {
		return respActNil, err
	}

	return aUsecase.doActivation(c, &acc, pa)
}

func (aUsecase *activationsUseCase) PostReactivations(c echo.Context, pa models.PayloadActivations) (models.RespActivations, error) {
	var respActNil models.RespActivations
	acc, err := aUsecase.rUsecase.CheckApplication(c, pa)

	if err != nil {
		return respActNil, err
	}

	err = aUsecase.cardRepo.GetCardStatus(c, &acc.Card)

	if err != nil {
		return respActNil, models.ErrUpdateCardStatus
	}

	// only do reactivation when card blocked and is reactivation no
	cardStatus := acc.Card.CardStatus
	cardStatus.LastEncryptedCardNumber = acc.Card.EncryptedCardNumber

	if !acc.Card.IsReactivationAvail() {
		return respActNil, models.ErrCannotReactivation
	}

	// make reactivation available
	acc.Card.EncryptedCardNumber = ""
	acc.AccountNumber = ""
	response, err := aUsecase.doActivation(c, &acc, pa)

	if err != nil {
		return response, err
	}

	go func() {
		cardStatus.IsReactivated = models.BoolYes
		cardStatus.ReactivatedDate = models.NowDbpg()
		cols := []string{"is_reactivated", "last_encrypted_card_number", "reactivated_date"}
		err = aUsecase.cardRepo.UpdateOneCardStatus(c, cardStatus, cols)

		if err != nil {
			logger.Make(c, nil).Debug(models.ErrUpdateCardStatus)
		}
	}()

	return response, nil
}

func (aUsecase *activationsUseCase) ValidateActivation(c echo.Context, pa models.PayloadActivations) models.ResponseErrors {
	var errors models.ResponseErrors
	// get account and check app number
	acc, err := aUsecase.rUsecase.CheckApplication(c, pa)

	if err != nil {
		errors.SetTitle(err.Error())
		return errors
	}

	// validate birth date if not equal
	err = aUsecase.validateBirthDate(acc, pa)

	if err != nil {
		errors.SetTitleCode("11", err.Error(), models.ErrPostActivationsFailed.Error())
		return errors
	}

	return errors
}

func (aUsecase *activationsUseCase) goldcardActivation(c echo.Context, acc *models.Account, pa models.PayloadActivations) error {
	var notif models.PdsNotification
	errActCore := make(chan error)
	errActBri := make(chan error)

	// activation to core
	go func() {
		if acc.AccountNumber != "" {
			errActCore <- nil

			return
		}

		err := aUsecase.arRepo.ActivationsToCore(c, acc)

		if err != nil {
			errActCore <- err

			return
		}

		errActCore <- nil
	}()

	// activation to BRI
	go func() {
		err := aUsecase.briActivation(c, acc, pa)

		if err != nil {
			errActBri <- err

			return
		}

		errActBri <- nil
	}()

	errCore := <-errActCore
	errBri := <-errActBri

	defer func() {
		_ = aUsecase.aRepo.PostActivations(c, *acc)
	}()

	if errCore != nil {
		return errCore
	}

	if errBri != nil {
		return errBri
	}

	acc.Application.Status = models.AppStatusActive
	acc.Card.Status = models.CardStatusActive
	acc.Status = models.AccStatusActive

	go func() {
		notif.GcActivation(*acc)
		_ = aUsecase.rrRepo.SendNotification(c, notif, "")
	}()

	return nil
}

func (aUsecase *activationsUseCase) validateBirthDate(acc models.Account, pa models.PayloadActivations) error {
	date, err := time.Parse(models.DDMMYYYY, pa.BirthDate)

	if err != nil {
		return err
	}

	birthDate := date.Format(models.DateFormatDef)

	if acc.PersonalInformation.BirthDate != birthDate {
		return models.ErrBirthDateNotMatch
	}

	return nil
}

func (aUsecase *activationsUseCase) isStlDecreased(c echo.Context, acc models.Account) (int64, int64, bool, error) {
	currStl, err := aUsecase.rrRepo.GetCurrentGoldSTL(c)

	if err != nil {
		return 0, 0, false, err
	}

	appliedStl := acc.Card.StlLimit
	deficitStl := appliedStl - currStl

	if deficitStl <= 0 {
		return 0, 0, false, nil
	}

	return currStl, deficitStl, true, nil
}

func (aUsecase *activationsUseCase) reRegistration(c echo.Context, acc models.Account, cardBal models.CardBalance) error {
	acc.Card.GoldLimit = cardBal.CurrGoldLimit
	acc.Card.GoldBalance = cardBal.CurrGoldLimit
	acc.Card.StlLimit = cardBal.CurrStl
	acc.Card.StlBalance = cardBal.CurrStl

	// recalculate open goldcard registrations
	err := aUsecase.rrRepo.OpenGoldcard(c, acc, true)

	if err != nil {
		return err
	}

	// update card gold limit and current stl
	err = aUsecase.aRepo.UpdateGoldLimit(c, acc.Card)

	if err != nil {
		return err
	}

	return nil
}

func (aUsecase *activationsUseCase) briActivation(c echo.Context, acc *models.Account, pa models.PayloadActivations) error {
	if acc.Card.EncryptedCardNumber != "" {
		return nil
	}

	err := aUsecase.arRepo.ActivationsToBRI(c, *acc, pa)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	// get card information
	cardInformation, err := aUsecase.trRepo.GetBRICardInformation(c, *acc)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	acc.Card.EncryptedCardNumber = cardInformation.BillKey
	acc.Card.ActivatedDate = time.Now()

	return nil
}

func (aUsecase *activationsUseCase) doActivation(c echo.Context, acc *models.Account, pa models.PayloadActivations) (models.RespActivations, error) {
	var respActNil models.RespActivations
	var errs models.ResponseErrors

	acc.Card.CardNumber = pa.FirstSixDigits + models.AppendXCardNumber + pa.LastFourDigits
	acc.Card.ValidUntil = pa.ExpDate

	// Inquiry activation
	cardBal, errs := aUsecase.InquiryActivation(c, *acc)

	if errs.Title != "" && errs.Code != "44" {
		return respActNil, errors.New(errs.Title)
	}

	// do re registration of the inquiry code 44
	if errs.Code == "44" && acc.AccountNumber == "" {
		err := aUsecase.reRegistration(c, *acc, cardBal)

		if err != nil {
			return respActNil, models.ErrPostActivationsFailed
		}
	}

	err := aUsecase.goldcardActivation(c, acc, pa)

	if err != nil {
		return respActNil, models.ErrPostActivationsFailed
	}

	return models.RespActivations{AccountNumber: acc.AccountNumber}, nil
}

func (aUsecase *activationsUseCase) isActivationForReplacedCard(card models.Card) bool {
	_ = aUsecase.cardRepo.GetCardStatus(nil, &card)
	return card.CardStatus.ID != 0
}
