package usecase

import (
	"gade/srv-goldcard/cards"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/transactions"

	"github.com/labstack/echo"
)

type cardsUseCase struct {
	cRepo    cards.Repository
	crRepo   cards.RestRepository
	tUseCase transactions.UseCase
}

// cardsUseCase represent cards Use Case
func CardsUseCase(cRepo cards.Repository, crRepo cards.RestRepository, tUseCase transactions.UseCase) cards.UseCase {
	return &cardsUseCase{cRepo, crRepo, tUseCase}
}

func (cus *cardsUseCase) BlockCard(c echo.Context, pl models.PayloadCardBlock) error {
	// Get Account
	acc, err := cus.tUseCase.CheckAccountByAccountNumber(c, pl)
	if err != nil {
		return models.ErrGetAccByAccountNumber
	}

	// Hit BRI based on reasonCode
	briCardBlockStatus, err := cus.crRepo.GetBRICardBlockStatus(c, acc, pl)

	if err != nil {
		return err
	}

	// Update Table Cards status to "inactive"
	acc.Card.Status = models.CardStatusInactive
	err = cus.cRepo.UpdateCardStatus(c, acc.Card)
	if err != nil {
		return models.ErrUpdateCardStatus
	}

	// Mapping Card Statuses
	cardStatuses := models.CardStatuses{}
	err = cardStatuses.MappingBlockCard(briCardBlockStatus, pl, acc.Card)
	if err != nil {
		return models.ErrMappingData
	}

	// Insert Table Card Statuses
	err = cus.cRepo.PostCardStatuses(cardStatuses)
	if err != nil {
		return models.ErrUpdateCardStatus
	}

	return nil
}
