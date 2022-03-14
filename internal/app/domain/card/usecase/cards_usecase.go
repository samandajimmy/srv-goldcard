package usecase

import (
	"srv-goldcard/internal/app/domain/card"
	"srv-goldcard/internal/app/domain/registration"
	"srv-goldcard/internal/app/domain/transaction"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"
	"time"

	"github.com/labstack/echo"
)

type cardsUseCase struct {
	cRepo    card.Repository
	crRepo   card.RestRepository
	trRepo   transaction.RestRepository
	tUseCase transaction.UseCase
	rRepo    registration.Repository
}

// cardsUseCase represent cards Use Case
func CardsUseCase(cRepo card.Repository, crRepo card.RestRepository, trRepo transaction.RestRepository,
	tUseCase transaction.UseCase, rRepo registration.Repository) card.UseCase {
	return &cardsUseCase{cRepo, crRepo, trRepo, tUseCase, rRepo}
}

func (cus *cardsUseCase) BlockCard(c echo.Context, pl model.PayloadCardBlock) error {
	// Get Account
	acc, err := cus.tUseCase.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		return model.ErrGetAccByAccountNumber
	}

	// Get Cards Status to validate if there user card blocking process still on progress
	err = cus.cRepo.GetCardStatus(c, &acc.Card)

	if err != nil {
		return model.ErrUpdateCardStatus
	}

	if acc.Card.CardStatus.ID != 0 {
		return model.ErrUserBlockProcessStillExisted
	}

	// Hit BRI based on reasonCode
	briCardBlockStatus, err := cus.crRepo.GetBRICardBlockStatus(c, acc, pl)

	if err != nil {
		return err
	}

	// run block a card process
	cardBlock := model.CardBlock{
		Reason:      pl.Reason,
		ReasonCode:  pl.ReasonCode,
		BlockedDate: time.Unix(briCardBlockStatus.ReportingDate/1000, 0).Format(model.DateFormat),
		Description: briCardBlockStatus.ReportDesc,
	}
	err = cus.blockaCard(c, cardBlock, &acc)

	if err != nil {
		return err
	}

	return nil
}

func (cus *cardsUseCase) GetCardStatus(c echo.Context, pl model.PayloadAccNumber) (model.RespCardStatus, error) {
	// Get Account
	acc, err := cus.tUseCase.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		return model.RespCardStatus{}, model.ErrGetAccByAccountNumber
	}

	// get card information
	cardInfo, err := cus.trRepo.GetBRICardInformation(c, acc)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return model.RespCardStatus{}, err
	}

	// check card block status from bri
	// update card status based on bri card info
	// the trigger when card is being blocked by BRI is by blockCode 'F' or 'L' and trfStatus '4'
	cardBlock := model.CardBlock{
		BlockedDate: cardInfo.BlockedDate,
		BlockedCode: cardInfo.BlockCode,
		ReasonCode:  model.ReasonCodeOther,
		TrfStatus:   cardInfo.TrfStatus,
	}

	// run block a card process
	err = cus.blockaCard(c, cardBlock, &acc)

	if err != nil {
		return model.RespCardStatus{}, err
	}

	// prepare response data card status
	err = cus.cRepo.GetCardStatus(c, &acc.Card)

	if err != nil {
		return model.RespCardStatus{}, model.ErrUpdateCardStatus
	}

	// do process replace card to BRI when needed
	err = cus.replaceaCard(c, &acc.Card.CardStatus, acc, cardInfo)

	if err != nil {
		return model.RespCardStatus{}, err
	}

	cardStatus := model.RespCardStatus{Status: acc.Card.Status}

	if acc.Card.CardStatus.CardID != 0 {
		cardStatus.IsReplaced = acc.Card.CardStatus.IsReplaced
	}

	return cardStatus, nil
}

// blockaCard is function to update card status, insert cardStatuses data, and do block goldcard account into core
func (cus *cardsUseCase) blockaCard(c echo.Context, cardBlock model.CardBlock, acc *model.Account) error {
	card := &acc.Card
	// do nothing when card has been blocked
	if card.Status == model.CardStatusBlocked {
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
		return model.ErrMappingData
	}

	// Update Table Cards status to "blocked" and Insert Table Card Statuses
	card.Status = model.CardStatusBlocked
	err = cus.cRepo.UpdateCardStatus(c, *card, cardStatus)

	if err != nil {
		return model.ErrUpdateCardStatus
	}

	return nil
}

// replaceaCard is function to replace card to BRI and update card status
func (cus *cardsUseCase) replaceaCard(c echo.Context, cardStatus *model.CardStatuses,
	acc model.Account, cardInfo model.BRICardBalance) error {
	// do nothing when card is already do replace process
	if cardStatus.IsReplaced == model.BoolYes {
		return nil
	}

	// wait trfStatus equal 4 for replace card to BRI
	if cardInfo.TrfStatus != "4" {
		return nil
	}

	plCardReplace := model.PayloadBriXkey{BriXkey: acc.BrixKey}
	err := cus.crRepo.PostCardReplaceBRI(c, plCardReplace)

	if err != nil {
		return model.ErrReplaceCard
	}

	// update card replaced data in db
	cardStatus.IsReplaced = model.BoolYes
	cardStatus.ReplacedDate = time.Now()
	cardStatus.LastEncryptedCardNumber = cardInfo.BillKey

	// create interface of data to update selected column
	col := []string{"is_replaced", "replaced_date", "last_encrypted_card_number"}
	err = cus.cRepo.UpdateOneCardStatus(c, *cardStatus, col)

	if err != nil {
		return model.ErrUpdateCardStatus
	}

	// reset application status to card processed
	err = cus.rRepo.ResetAppStatusToCardProcessed(acc.ApplicationID)

	if err != nil {
		return model.ErrUpdateCardStatus
	}

	return nil
}

func (cus *cardsUseCase) CloseCard(c echo.Context, pl model.PayloadCIF) error {
	// Get Account
	acc, err := cus.tUseCase.CheckAccountByCIF(c, pl)

	if err != nil {
		return model.ErrGetAccByCIF
	}

	// Close Card (Update status to inactive in cards, application, account)
	err = cus.cRepo.SetInactiveStatus(c, acc)

	if err != nil {
		return err
	}

	err = cus.crRepo.PdsSetNullAppAccNumber(c, pl)

	if err != nil {
		return err
	}

	return nil
}
