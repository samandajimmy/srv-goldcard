package usecase

import (
	"errors"
	"gade/srv-goldcard/activations"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"reflect"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

type activationsUseCase struct {
	aRepo  activations.Repository
	arRepo activations.RestRepository
	rRepo  registrations.Repository
	rrRepo registrations.RestRepository
}

// ActivationUseCase represent Activation Use Case
func ActivationUseCase(aRepo activations.Repository, arRepo activations.RestRepository,
	rRepo registrations.Repository, rrRepo registrations.RestRepository) activations.UseCase {
	return &activationsUseCase{aRepo, arRepo, rRepo, rrRepo}
}

func (aUsecase *activationsUseCase) InquiryActivation(c echo.Context, pl models.PayloadAppNumber) models.ResponseErrors {
	var errors models.ResponseErrors
	// get account and check app number
	acc, err := aUsecase.checkApplication(c, pl)

	if err != nil {
		errors.SetTitle(err.Error())
		return errors
	}

	// validation on inquiry
	// validate application expiry from application_processed_date < 12 months
	// add a year for expiry date
	expDate := acc.Application.ApplicationProcessedDate.AddDate(1, 0, 0)

	if acc.Application.ApplicationProcessedDate.After(expDate) {
		errors.SetTitleCode("22", models.ErrAppExpired.Error(), models.ErrAppExpiredDesc.Error())
		return errors
	}

	// validate stl price changes
	// compare stl price at applied date and current date
	currStl, err := aUsecase.rrRepo.GetCurrentGoldSTL(c)

	if err != nil {
		errors.SetTitle(models.ErrGetCurrSTL.Error())
		return errors
	}

	appliedStl := acc.Card.CurrentSTL
	deficitStl := appliedStl - currStl

	if deficitStl <= 0 {
		return errors
	}

	// if it decreased
	// if the decrase <= 1,15% then go head
	decreasedPercent := models.CustomRound("round", float64(deficitStl)/float64(currStl), 10000)

	if decreasedPercent <= models.DecreasedLimit {
		return errors
	}

	// if the decrase > 1,15% then
	// get user effective balance
	userDetail, err := aUsecase.arRepo.GetDetailGoldUser(c, acc.Application.SavingAccount)

	if err != nil {
		errors.SetTitle(models.ErrGetUserDetail.Error())
		return errors
	}

	goldEffBalance, err := strconv.ParseFloat(userDetail["saldoEfektif"], 64)

	if err != nil {
		errors.SetTitle(models.ErrGetEffBalance.Error())
		return errors
	}

	appliedGoldLimit := acc.Card.GoldLimit
	currGoldLimit := acc.Card.ConvertMoneyToGold(acc.Card.CardLimit, currStl)
	// because we need user to have at least 0.1 effective gold balance
	deficitGoldLimit := models.CustomRound("round", currGoldLimit-appliedGoldLimit, 10000) + models.MinEffBalance

	// got not enough effective gold balance
	if goldEffBalance < deficitGoldLimit {
		errors.SetTitleCode("55", models.ErrDecreasedSTL.Error(), models.ErrDecreasedSTLDesc.Error())
		return errors
	}

	acc.Card.GoldLimit = currGoldLimit
	acc.Card.CurrentSTL = currStl
	// recalculate open goldcard registrations
	err = aUsecase.rrRepo.OpenGoldcard(c, acc, true)

	if err != nil {
		errors.SetTitle(err.Error())
		return errors
	}

	// update card gold limit and current stl
	err = aUsecase.aRepo.UpdateGoldLimit(c, acc.Card)

	if err != nil {
		errors.SetTitle(models.ErrUpdateCardLimit.Error())
		return errors
	}

	return errors
}

func (aUsecase *activationsUseCase) PostActivations(c echo.Context, pa models.PayloadActivations) error {
	acc, err := aUsecase.checkApplication(c, pa)

	if err != nil {
		return err
	}

	err = acc.MappingCardActivationsData(c, pa)

	if err != nil {
		return models.ErrMappingData
	}

	// Inquiry activation
	if models.DateIsNotEqual(acc.Card.UpdatedAt, time.Now()) {
		appNumber := models.PayloadAppNumber{
			ApplicationNumber: acc.Application.ApplicationNumber,
		}

		err := aUsecase.InquiryActivation(c, appNumber)

		if err.Title != "" {
			return errors.New(err.Title)
		}
	}

	// Activations to BRI
	errBri := aUsecase.arRepo.ActivationsToBRI(c, acc, pa)

	if errBri != nil {
		return errBri
	}

	// Activations to core
	errSwitching := aUsecase.arRepo.ActivationsToCore(c, acc)

	if errSwitching != nil {
		return errSwitching
	}

	errUpdateAct := aUsecase.aRepo.PostActivations(c, acc)

	if errUpdateAct != nil {
		return models.ErrPostActivationsFailed
	}

	return nil
}

func (aUsecase *activationsUseCase) checkApplication(c echo.Context, pl interface{}) (models.Account, error) {
	r := reflect.ValueOf(pl)
	appNumber := r.FieldByName("ApplicationNumber")

	if appNumber.IsZero() {
		return models.Account{}, nil
	}

	acc := models.Account{Application: models.Applications{ApplicationNumber: appNumber.String()}}
	err := aUsecase.aRepo.GetAccountByAppNumber(c, &acc)

	if err != nil {
		return models.Account{}, models.ErrAppNumberNotFound
	}

	return acc, nil
}
