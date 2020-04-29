package update_limits

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the transactions Repository
type Repository interface {
	GetEmailByKey(c echo.Context) (string, error)
	GetDocumentByTypeAndApplicationId(appId int64, docType string) (models.Document, error)
}

// Repository represent the update limits Rest Repository
type RestRepository interface {
	CorePostUpdateLimit(c echo.Context, savingAccNum string, card models.Card) error
	BRIPostUpdateLimit(c echo.Context, acc models.Account, doc models.Document) error
}
