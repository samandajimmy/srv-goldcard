package usecase

import (
	"errors"
	"srv-goldcard/internal/app/domain/activation"
	"srv-goldcard/internal/app/domain/card"
	"srv-goldcard/internal/app/domain/registration"
	"srv-goldcard/internal/app/domain/transaction"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

type activationsUseCase struct {
	aRepo    activation.Repository
	arRepo   activation.RestRepository
	rRepo    registration.Repository
	rrRepo   registration.RestRepository
	rUsecase registration.UseCase
	trRepo   transaction.RestRepository
	cardRepo card.Repository
}

// ActivationUseCase represent Activation Use Case
func ActivationUseCase(aRepo activation.Repository, arRepo activation.RestRepository,
	rRepo registration.Repository, rrRepo registration.RestRepository, rUsecase registration.UseCase,
	trRepo transaction.RestRepository, cardRepo card.Repository) activation.UseCase {
	return &activationsUseCase{aRepo, arRepo, rRepo, rrRepo, rUsecase, trRepo, cardRepo}
}

func (aUsecase *activationsUseCase) InquiryActivation(c echo.Context, acc model.Account) (model.CardBalance, model.ResponseErrors) {
	var errors model.ResponseErrors
	var cardBal model.CardBalance
	var err error

	// get account and check app number
	if acc.ID == int64(0) {
		appNumber := model.PayloadAppNumber{
			ApplicationNumber: acc.Application.ApplicationNumber,
		}

		acc, err = aUsecase.rUsecase.CheckApplication(c, appNumber)
	}

	if err != nil {
		errors.SetTitle(err.Error())
		return cardBal, errors
	}

	// if account or card has been activated
	if acc.Status == model.AccStatusActive &&
		acc.Application.Status == model.AppStatusActive &&
		acc.Card.Status == model.CardStatusActive {
		errors.SetTitle(model.ErrCardActivated.Error())
		return cardBal, errors
	}

	// if card is being replaced skip activation validation process
	// the condition if there is some record on Card Status table with corresponding card_id then the user has been activated before and request for card replacement
	err = aUsecase.cardRepo.GetCardStatus(c, &acc.Card)

	if err != nil {
		errors.SetTitle(model.ErrInqActivation.Error())
		return cardBal, errors
	}

	if acc.Card.CardStatus.ID != 0 {
		return cardBal, errors
	}

	// validation on inquiry
	// validate application expiry from application_processed_date < 12 months
	// add a year for expiry date
	expDate := acc.Application.ApplicationProcessedDate.AddDate(1, 0, 0)
	now := model.NowUTC()

	if now.After(expDate) {
		errors.SetTitleCode("22", model.ErrAppExpired.Error(), model.ErrAppExpiredDesc.Error())
		return cardBal, errors
	}

	// validate stl price changes
	// compare stl price at applied date and current date
	currStl, deficitStl, isDecreased, err := aUsecase.isStlDecreased(c, acc)

	if err != nil {
		errors.SetTitle(model.ErrGetCurrSTL.Error())
		return cardBal, errors
	}

	// return response if there's no decreasing
	if !isDecreased {
		return cardBal, errors
	}

	// get decreasing percentage
	_ = model.CustomRound("round", float64(deficitStl)/float64(currStl), 10000)
	// get user effective balance
	userDetail, err := aUsecase.arRepo.GetDetailGoldUser(c, acc.Application.SavingAccount)

	if err != nil {
		errors.SetTitle(model.ErrGetUserDetail.Error())
		return cardBal, errors
	}

	if _, ok := userDetail["saldoEfektif"].(string); !ok {
		logger.Make(c, nil).Debug(model.ErrSetVar)
		errors.SetTitle(model.ErrSetVar.Error())
		return cardBal, errors
	}

	goldEffBalance, err := strconv.ParseFloat(userDetail["saldoEfektif"].(string), 64)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		errors.SetTitle(model.ErrGetEffBalance.Error())

		return cardBal, errors
	}

	appliedGoldLimit := acc.Card.GoldLimit
	currGoldLimit := acc.Card.SetGoldLimit(acc.Card.CardLimit, currStl)
	// because we need user to have at least 0.1 effective gold balance
	deficitGoldLimit := model.CustomRound("round", currGoldLimit-appliedGoldLimit, 10000)
	cardBal.CurrGoldLimit = currGoldLimit
	cardBal.CurrStl = currStl
	cardBal.DeficitGoldLimit = deficitGoldLimit
	cardBal.SavingAccount = acc.Application.SavingAccount

	// got not enough effective gold balance
	if goldEffBalance < deficitGoldLimit+model.MinEffBalance {
		errors.SetTitleCode("55", model.ErrDecreasedSTL.Error(), model.ErrDecreasedSTLDesc.Error())
		return cardBal, errors
	}

	// got enough effective gold balance
	errors.SetTitleCode(
		"44",
		model.ErrDecreasedSTL.Error(),
		model.DynamicErr(model.ErrDecreasedSTLOpenDesc, []interface{}{deficitGoldLimit}).Error(),
	)

	return cardBal, errors
}

func (aUsecase *activationsUseCase) ForceActivation(c echo.Context, acc model.Account) (model.RespActivations, error) {
	var respActNil model.RespActivations

	appNumber := model.PayloadAppNumber{
		ApplicationNumber: acc.Application.ApplicationNumber,
	}

	acc, err := aUsecase.rUsecase.CheckApplication(c, appNumber)

	if err != nil {
		return respActNil, err
	}

	err = aUsecase.goldcardActivation(c, &acc, model.PayloadActivations{IsForced: true})

	if err != nil {
		return respActNil, model.ErrPostActivationsFailed
	}

	return respActNil, err
}

func (aUsecase *activationsUseCase) PostActivations(c echo.Context, pa model.PayloadActivations) (model.RespActivations, error) {
	var respActNil model.RespActivations
	acc, err := aUsecase.rUsecase.CheckApplication(c, pa)

	if err != nil {
		return respActNil, err
	}

	return aUsecase.doActivation(c, &acc, pa)
}

func (aUsecase *activationsUseCase) PostReactivations(c echo.Context, pa model.PayloadActivations) (model.RespActivations, error) {
	var respActNil model.RespActivations
	acc, err := aUsecase.rUsecase.CheckApplication(c, pa)

	if err != nil {
		return respActNil, err
	}

	err = aUsecase.cardRepo.GetCardStatus(c, &acc.Card)

	if err != nil {
		return respActNil, model.ErrUpdateCardStatus
	}

	// only do reactivation when card blocked and is reactivation no
	cardStatus := acc.Card.CardStatus
	cardStatus.LastEncryptedCardNumber = acc.Card.EncryptedCardNumber

	if !acc.Card.IsReactivationAvail() {
		return respActNil, model.ErrCannotReactivation
	}

	defer func() {
		if err != nil {
			return
		}

		cardStatus.IsReactivated = model.BoolYes
		cardStatus.ReactivatedDate = model.NowDbpg()
		cols := []string{"is_reactivated", "last_encrypted_card_number", "reactivated_date"}
		err = aUsecase.cardRepo.UpdateOneCardStatus(c, cardStatus, cols)

		if err != nil {
			logger.Make(c, nil).Debug(model.ErrUpdateCardStatus)
		}
	}()

	if pa.IsForced {
		err = aUsecase.goldcardActivation(c, &acc, pa)

		return respActNil, err
	}

	// make reactivation available
	acc.Card.EncryptedCardNumber = ""
	acc.AccountNumber = ""
	response, err := aUsecase.doActivation(c, &acc, pa)

	if err != nil {
		return response, err
	}

	return response, nil
}

func (aUsecase *activationsUseCase) ValidateActivation(c echo.Context, pa model.PayloadActivations) model.ResponseErrors {
	var errors model.ResponseErrors

	// skip validation when its forced
	if pa.IsForced {
		return errors
	}

	// get account and check app number
	acc, err := aUsecase.rUsecase.CheckApplication(c, pa)

	if err != nil {
		errors.SetTitle(err.Error())
		return errors
	}

	// validate birth date if not equal
	err = aUsecase.validateBirthDate(acc, pa)

	if err != nil {
		errors.SetTitleCode("11", err.Error(), model.ErrPostActivationsFailed.Error())
		return errors
	}

	return errors
}

func (aUsecase *activationsUseCase) goldcardActivation(c echo.Context, acc *model.Account, pa model.PayloadActivations) error {
	var notif model.PdsNotification
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

	acc.Application.Status = model.AppStatusActive
	acc.Card.Status = model.CardStatusActive
	acc.Status = model.AccStatusActive

	go func() {
		notif.GcActivation(*acc)
		_ = aUsecase.rrRepo.SendNotification(c, notif, "")
	}()

	return nil
}

func (aUsecase *activationsUseCase) validateBirthDate(acc model.Account, pa model.PayloadActivations) error {
	date, err := time.Parse(model.DDMMYYYY, pa.BirthDate)

	if err != nil {
		return err
	}

	birthDate := date.Format(model.DateFormatDef)

	if acc.PersonalInformation.BirthDate != birthDate {
		return model.ErrBirthDateNotMatch
	}

	return nil
}

func (aUsecase *activationsUseCase) isStlDecreased(c echo.Context, acc model.Account) (int64, int64, bool, error) {
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

func (aUsecase *activationsUseCase) reRegistration(c echo.Context, acc model.Account, cardBal model.CardBalance) error {
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

func (aUsecase *activationsUseCase) briActivation(c echo.Context, acc *model.Account, pa model.PayloadActivations) error {
	fnMapBillKey := func(c echo.Context, acc *model.Account) error {
		// get card information
		cardInformation, err := aUsecase.trRepo.GetBRICardInformation(c, *acc)

		if err != nil {
			return err
		}

		acc.Card.EncryptedCardNumber = cardInformation.BillKey
		acc.Card.ActivatedDate = time.Now()

		return nil
	}

	if pa.IsForced {
		return fnMapBillKey(c, acc)
	}

	if acc.Card.EncryptedCardNumber != "" {
		return nil
	}

	err := aUsecase.arRepo.ActivationsToBRI(c, *acc, pa)

	if err != nil {
		return err
	}

	err = fnMapBillKey(c, acc)

	if err != nil {
		return err
	}

	return nil
}

func (aUsecase *activationsUseCase) doActivation(c echo.Context, acc *model.Account, pa model.PayloadActivations) (model.RespActivations, error) {
	var respActNil model.RespActivations
	var errs model.ResponseErrors

	acc.Card.CardNumber = pa.FirstSixDigits + model.AppendXCardNumber + pa.LastFourDigits
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
			return respActNil, model.ErrPostActivationsFailed
		}
	}

	err := aUsecase.goldcardActivation(c, acc, pa)

	if err != nil {
		return respActNil, model.ErrPostActivationsFailed
	}

	return model.RespActivations{AccountNumber: acc.AccountNumber}, nil
}
