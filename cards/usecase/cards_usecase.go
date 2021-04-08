package usecase

import (
	"gade/srv-goldcard/cards"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/transactions"
	"time"

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
		BlockedDate: time.Unix(briCardBlockStatus.ReportingDate/1000, 0).Format(models.DateTimeFormat),
		Description: briCardBlockStatus.ReportDesc,
	}
	err = cus.blockaCard(c, cardBlock, &acc)

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
	err = cus.blockaCard(c, cardBlock, &acc)

	if err != nil {
		return models.RespCardStatus{}, err
	}

	// prepare response data card status
	err = cus.cRepo.GetCardStatus(c, &acc.Card)

	if err != nil {
		return models.RespCardStatus{}, models.ErrUpdateCardStatus
	}

	// do process replace card to BRI when needed
	err = cus.replaceaCard(c, &acc.Card.CardStatus, acc, cardInfo)

	if err != nil {
		return models.RespCardStatus{}, err
	}

	cardStatus := models.RespCardStatus{Status: acc.Card.Status}

	if acc.Card.CardStatus.CardID != 0 {
		cardStatus.IsReplaced = acc.Card.CardStatus.IsReplaced
	}

	return cardStatus, nil
}

// blockaCard is function to update card status, insert cardStatuses data, and do block goldcard account into core
func (cus *cardsUseCase) blockaCard(c echo.Context, cardBlock models.CardBlock, acc *models.Account) error {
	card := &acc.Card
	// do nothing when card has been blocked
	if card.Status == models.CardStatusBlocked {
		return nil
	}

	if !cardBlock.IsCardBlockedBri() {
		return nil
	}

	// hit endpoint core to block a card
	err := cus.crRepo.CoreBlockaCard(c, *acc, cardBlock)

	if err != nil {
		return err
	}

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

// replaceaCard is function to replace card to BRI and update card status
func (cus *cardsUseCase) replaceaCard(c echo.Context, cardStatus *models.CardStatuses,
	acc models.Account, cardInfo models.BRICardBalance) error {
	// do nothing when card is already do replace process
	if cardStatus.IsReplaced == models.BoolYes {
		return nil
	}

	// wait trfStatus equal 4 for replace card to BRI
	if cardInfo.TrfStatus != "4" {
		return nil
	}

	plCardReplace := models.PayloadBriXkey{BriXkey: acc.BrixKey}
	_, err := cus.crRepo.PostCardReplaceBRI(c, plCardReplace)

	if err != nil {
		return models.ErrReplaceCard
	}

	// update card replaced data in db
	cardStatus.IsReplaced = models.BoolYes
	cardStatus.ReplacedDate = time.Now()
	cardStatus.LastEncryptedCardNumber = cardInfo.BillKey

	// create interface of data to update selected column
	col := []string{"is_replaced", "replaced_date", "last_encrypted_card_number"}
	err = cus.cRepo.UpdateOneCardStatus(c, *cardStatus, col)

	if err != nil {
		return models.ErrUpdateCardStatus
	}

	return nil
}
