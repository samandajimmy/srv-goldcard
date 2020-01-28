package usecase

import (
	"errors"
	"fmt"
	"gade/srv-goldcard/activations"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"strconv"

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

func (aUsecase *activationsUseCase) InquiryActivation(c echo.Context, pl models.PayloadAppNumber) error {
	// get account
	acc := models.Account{Application: models.Applications{ApplicationNumber: pl.ApplicationNumber}}
	err := aUsecase.rRepo.GetAccountByAppNumber(c, &acc)

	if err != nil {
		return models.ErrAppNumberNotFound
	}

	// validation on inquiry
	// validate application expiry from application_processed_date < 12 months
	// if expired give error message "PENGAJUAN KADALUARSA : Pengajuan harus dibatalkan karena tidak
	// ada aktivitas selama 12 bulan. Saldo emas akan dikembalikan ke saldo efektif." dan Button "Oke, Batalkan Pengajuan"
	// add a year for expiry date
	expDate := acc.Application.ApplicationProcessedDate.AddDate(1, 0, 0)

	if acc.Application.ApplicationProcessedDate.After(expDate) {
		// TODO: change the error message
		return errors.New("Pengajuan harus dibatalkan")
	}

	// validate stl price changes
	// compare stl price at applied date and current date
	currStl, err := aUsecase.rrRepo.GetCurrentGoldSTL(c)

	fmt.Println("currStl")
	fmt.Println(currStl)
	fmt.Println("currStl")

	if err != nil {
		// TODO: change the error message
		return errors.New("something went wrong when trying to get STL")
	}

	appliedStl := acc.Card.CurrentSTL
	deficitStl := appliedStl - currStl

	fmt.Println("appliedStl")
	fmt.Println(appliedStl)
	fmt.Println("appliedStl")

	fmt.Println("deficitStl")
	fmt.Println(deficitStl)
	fmt.Println("deficitStl")

	if deficitStl <= 0 {
		return nil
	}

	// if it decreased
	// if the decrase <= 1,15% then go head
	decreasedPercent := models.CustomRound("round", float64(deficitStl)/float64(currStl), 10000)

	fmt.Println("decreasedPercent")
	fmt.Println(decreasedPercent)
	fmt.Println("decreasedPercent")

	if decreasedPercent <= models.DecreasedLimit {
		return nil
	}

	// if the decrase > 1,15% then
	// get user effective balance
	userDetail, err := aUsecase.arRepo.GetDetailGoldUser(c, acc.Application.SavingAccount)

	if err != nil {
		// TODO: change the error message
		return errors.New("something went wrong when trying to get user detail")
	}

	goldEffBalance, err := strconv.ParseFloat(userDetail["saldoEfektif"], 64)

	if err != nil {
		// TODO: change the error message
		return errors.New("something went wrong when trying to gold effective balance")
	}

	appliedGoldLimit := acc.Card.GoldLimit
	currGoldLimit := acc.Card.ConvertMoneyToGold(acc.Card.CardLimit, currStl)
	deficitGoldLimit := models.CustomRound("round", currGoldLimit-appliedGoldLimit, 10000)

	// gold effective balance is less then 0.1000 gram
	if goldEffBalance < models.EffBalLimit {
		// TODO: change the error message
		return errors.New("effective balance is not enough for balance limit")
	}

	// gold effective balance is less then deficit gold limit
	if goldEffBalance < deficitGoldLimit {
		// TODO: change the error message
		return errors.New("effective balance is not enough for balance limit")
	}

	// update card gold limit and current stl

	fmt.Println("goldEffBal")
	fmt.Println(appliedGoldLimit)
	fmt.Println(currGoldLimit)
	fmt.Println("goldEffBal")

	fmt.Println("deficitGoldLimit")
	fmt.Println(deficitGoldLimit)
	fmt.Println("deficitGoldLimit")

	fmt.Println("acc.Card.CardLimit")
	fmt.Println(acc.Card.CardLimit)
	fmt.Println("acc.Card.CardLimit")

	fmt.Println("goldEffBalance")
	fmt.Println(goldEffBalance)
	fmt.Println("goldEffBalance")

	return nil
}
