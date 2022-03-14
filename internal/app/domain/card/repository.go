package card

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// Repository represent the cards repository contract
type Repository interface {
	UpdateCardStatus(c echo.Context, card model.Card, cs model.CardStatuses) error
	GetCardStatus(c echo.Context, card *model.Card) error
	UpdateOneCardStatus(c echo.Context, cardStatus model.CardStatuses, cols []string) error
	SetInactiveStatus(c echo.Context, acc model.Account) error
}

// RestRepository represent the rest cards repository contract
type RestRepository interface {
	GetBRICardBlockStatus(c echo.Context, acc model.Account, pl model.PayloadCardBlock) (model.BRICardBlockStatus, error)
	PostCardReplaceBRI(c echo.Context, pl model.PayloadBriXkey) error
	CoreBlockaCard(c echo.Context, acc model.Account, cardBlock model.CardBlock) error
	PdsSetNullAppAccNumber(c echo.Context, cif model.PayloadCIF) error
}
