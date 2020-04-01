package usecase

import (
	"errors"
	"gade/srv-goldcard/activations"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
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
}

// ActivationUseCase represent Activation Use Case
func ActivationUseCase(aRepo activations.Repository, arRepo activations.RestRepository,
	rRepo registrations.Repository, rrRepo registrations.RestRepository, rUsecase registrations.UseCase) activations.UseCase {
	return &activationsUseCase{aRepo, arRepo, rRepo, rrRepo, rUsecase}
}

func (aUsecase *activationsUseCase) InquiryActivation(c echo.Context, pl models.PayloadAppNumber) models.ResponseErrors {
	var errors models.ResponseErrors
	// get account and check app number
	acc, err := aUsecase.rUsecase.CheckApplication(c, pl)

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

	appliedStl := acc.Card.StlLimit
	deficitStl := appliedStl - currStl

	if deficitStl <= 0 {
		return errors
	}

	// if it decreased
	// get decreasing percentage
	_ = models.CustomRound("round", float64(deficitStl)/float64(currStl), 10000)
	// get user effective balance
	userDetail, err := aUsecase.arRepo.GetDetailGoldUser(c, acc.Application.SavingAccount)

	if err != nil {
		errors.SetTitle(models.ErrGetUserDetail.Error())
		return errors
	}

	if _, ok := userDetail["saldoEfektif"].(string); !ok {
		errors.SetTitle(models.ErrSetVar.Error())
		return errors
	}

	goldEffBalance, err := strconv.ParseFloat(userDetail["saldoEfektif"].(string), 64)

	if err != nil {
		errors.SetTitle(models.ErrGetEffBalance.Error())
		return errors
	}

	appliedGoldLimit := acc.Card.GoldLimit
	currGoldLimit := acc.Card.SetGoldLimit(acc.Card.CardLimit, currStl)
	// because we need user to have at least 0.1 effective gold balance
	deficitGoldLimit := models.CustomRound("round", currGoldLimit-appliedGoldLimit, 10000) + models.MinEffBalance

	// got not enough effective gold balance
	if goldEffBalance < deficitGoldLimit {
		errors.SetTitleCode("55", models.ErrDecreasedSTL.Error(), models.ErrDecreasedSTLDesc.Error())
		return errors
	}

	acc.Card.GoldLimit = currGoldLimit
	acc.Card.GoldBalance = currGoldLimit
	acc.Card.StlLimit = currStl
	acc.Card.StlBalance = currStl
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

func (aUsecase *activationsUseCase) PostActivations(c echo.Context, pa models.PayloadActivations) (models.RespActivations, error) {
	var respActNil models.RespActivations
	acc, err := aUsecase.rUsecase.CheckApplication(c, pa)

	if err != nil {
		return respActNil, err
	}

	err = acc.MappingCardActivationsData(c, pa)

	if err != nil {
		return respActNil, models.ErrMappingData
	}

	// Inquiry activation
	if models.DateIsNotEqual(acc.Card.UpdatedAt, time.Now()) {
		appNumber := models.PayloadAppNumber{
			ApplicationNumber: acc.Application.ApplicationNumber,
		}

		err := aUsecase.InquiryActivation(c, appNumber)

		if err.Title != "" {
			return respActNil, errors.New(err.Title)
		}
	}

	// init activation channel
	errActCore := make(chan error)

	go func() {
		// Activations to core
		err = aUsecase.arRepo.ActivationsToCore(c, &acc)

		if err != nil {
			errActCore <- err

			return
		}

		errActCore <- nil
	}()

	err = aUsecase.afterActivationGoldcard(c, &acc, pa, errActCore)

	if err != nil {
		return respActNil, models.ErrPostActivationsFailed
	}

	return models.RespActivations{AccountNumber: acc.AccountNumber}, nil
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

func (aUsecase *activationsUseCase) afterActivationGoldcard(c echo.Context, acc *models.Account, pa models.PayloadActivations, errActCore chan error) error {
	var notif models.PdsNotification
	errActBri := make(chan error)
	errActUpdate := make(chan error)
	errActivation := make(chan error)

	// Activations to BRI
	briActivation := func() {
		err := aUsecase.arRepo.ActivationsToBRI(c, *acc, pa)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			errActBri <- err
			return
		}

		errActBri <- nil
	}

	updateActivation := func() {
		err := aUsecase.aRepo.PostActivations(c, *acc)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			errActUpdate <- err
			return
		}

		errActUpdate <- nil
	}

	sendSucceededNotif := func() {
		notif.GcActivation(*acc, "succeeded")
		_ = aUsecase.rrRepo.SendNotification(c, notif, "")
	}

	sendFailedNotif := func() {
		notif.GcActivation(*acc, "failed")
		_ = aUsecase.rrRepo.SendNotification(c, notif, "")
	}

	go func() {
		for {
			select {
			case err := <-errActCore:
				if err == nil {
					go briActivation()
				}

				if err != nil {
					// send notif activation failed
					go sendFailedNotif()
					errActivation <- err
				}
			case err := <-errActBri:
				if err == nil {
					go updateActivation()
				}

				if err != nil {
					// send notif activation failed
					go sendFailedNotif()
					errActivation <- err
				}
			case err := <-errActUpdate:
				if err == nil {
					// send notif activation succeeded
					go sendSucceededNotif()
					errActivation <- nil
				}

				if err != nil {
					// send notif activation failed
					go sendFailedNotif()
					errActivation <- err
				}
			}
		}
	}()

	if err := <-errActivation; err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}

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
