package usecase

import (
	"gade/srv-goldcard/cards"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/transactions"

	"github.com/labstack/echo"
)

type cardsUseCase struct {
	cRepo    cards.Repository
	crRepo   cards.RestRepository
	trRepo   transactions.RestRepository
	tUseCase transactions.UseCase
}

// cardsUseCase represent cards Use Case
func CardsUseCase(cRepo cards.Repository, crRepo cards.RestRepository, trRepo transactions.RestRepository,
	tUseCase transactions.UseCase) cards.UseCase {
	return &cardsUseCase{cRepo, crRepo, trRepo, tUseCase}
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

	// run block a card process
	cardBlock := models.CardBlock{
		Reason:      pl.Reason,
		ReasonCode:  pl.ReasonCode,
		BlockedDate: briCardBlockStatus.ReportingDate,
	}
	err = cus.blockaCard(c, cardBlock, &acc.Card)

	if err != nil {
		return err
	}

	return nil
}

func (cus *cardsUseCase) GetCardStatus(c echo.Context, pl models.PayloadAccNumber) (models.RespCardStatus, error) {
	// Get Account
	acc, err := cus.tUseCase.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		return models.RespCardStatus{}, models.ErrGetAccByAccountNumber
	}

	// get card information
	cardInfo, err := cus.trRepo.GetBRICardInformation(c, acc)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return models.RespCardStatus{}, err
	}

	// check card block status from bri
	// update card status based on bri card info
	cardBlock := models.CardBlock{
		BlockedDate: cardInfo.BlockedDate,
		BlockedCode: cardInfo.BlockCode,
		ReasonCode:  models.ReasonCodeOther,
	}

	// run block a card process
	err = cus.blockaCard(c, cardBlock, &acc.Card)

	if err != nil {
		return models.RespCardStatus{}, err
	}

	// TODO: check trfStatus from bri and do process

	// prepare response data card status
	err = cus.cRepo.GetCardStatus(c, &acc.Card)

	if err != nil {
		return models.RespCardStatus{}, models.ErrUpdateCardStatus
	}

	cardStatus := models.RespCardStatus{Status: acc.Card.Status}

	if acc.Card.CardStatus.CardID != 0 {
		cardStatus.IsReplaced = acc.Card.CardStatus.IsReplaced
	}

	return cardStatus, nil
}

func (cus *cardsUseCase) blockaCard(c echo.Context, cardBlock models.CardBlock, card *models.Card) error {
	// do nothing when card has been blocked
	if card.Status == models.CardStatusBlocked {
		return nil
	}

	if !cardBlock.IsCardBlockedBri() {
		return nil
	}

	// TODO: hit endpoint core to block a card

	// Mapping Card Statuses
	cardStatus, err := card.MappingBlockCard(cardBlock)

	if err != nil {
		return models.ErrMappingData
	}

	// Update Table Cards status to "blocked" and Insert Table Card Statuses
	card.Status = models.CardStatusBlocked
	err = cus.cRepo.UpdateCardStatus(c, *card, cardStatus)

	if err != nil {
		return models.ErrUpdateCardStatus
	}

	return nil
}
