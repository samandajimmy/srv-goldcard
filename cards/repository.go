package cards

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the cards repository contract
type Repository interface {
	UpdateCardStatus(c echo.Context, card models.Card, cs models.CardStatuses) error
	GetCardStatus(c echo.Context, card *models.Card) error
	UpdateOneCardStatus(c echo.Context, cardStatus models.CardStatuses, cols []string) error
	SetInactiveStatus(c echo.Context, acc models.Account) error
}

// RestRepository represent the rest cards repository contract
type RestRepository interface {
	GetBRICardBlockStatus(c echo.Context, acc models.Account, pl models.PayloadCardBlock) (models.BRICardBlockStatus, error)
	PostCardReplaceBRI(c echo.Context, pl models.PayloadBriXkey) error
	CoreBlockaCard(c echo.Context, acc models.Account, cardBlock models.CardBlock) error
}
