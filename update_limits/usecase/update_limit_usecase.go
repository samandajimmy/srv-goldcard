package usecase

import (
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"gade/srv-goldcard/transactions"
	"gade/srv-goldcard/update_limits"

	"github.com/labstack/echo"
)

type updateLimitUseCase struct {
	trxRepo transactions.Repository
	rRepo   registrations.Repository
	rrRepo  registrations.RestRepository
}

// UpdateLimitUseCase represent Update Limit Use Case
func UpdateLimitUseCase(trxRepo transactions.Repository, rRepo registrations.Repository,
	rrRepo registrations.RestRepository) update_limits.UseCase {
	return &updateLimitUseCase{trxRepo, rRepo, rrRepo}
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
	}

	return errors
}
