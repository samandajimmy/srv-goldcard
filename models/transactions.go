package models

import (
	"time"

	"github.com/labstack/echo"
)

const (
	typeTrxCredit string = "credit"

	statusTrxPending string = "pending"
)

// Transaction is a struct to store transaction data
type Transaction struct {
	ID          int64     `json:"id"`
	AccountId   int64     `json:"accountId"`
	RefTrxPgdn  string    `json:"refTrxPgdn"`
	RefTrx      string    `json:"refTrx"`
	Nominal     int64     `json:"nominal"`
	GoldNominal float64   `json:"goldNominal"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Balance     int64     `json:"balance"`
	GoldBalance float64   `json:"goldBalance"`
	Methods     string    `json:"methods"`
	TrxDate     string    `json:"trxDate"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	Account     Account   `json:"account"`
}

// MappingTransactions is a struct to mapping transactions data
func (trx *Transaction) MappingTransactions(c echo.Context, pl PayloadBRIPendingTransactions, trans Transaction, refTrxPg string, stl int64) error {
	goldNominal := trx.Account.Card.ConvertMoneyToGold(pl.Amount, stl)
	balance := trx.Account.Card.Balance - pl.Amount
	goldBalance := trx.Account.Card.ConvertMoneyToGold(balance, stl)

	trx.AccountId = trans.Account.ID
	trx.RefTrxPgdn = refTrxPg
	trx.RefTrx = pl.TransactionId
	trx.Nominal = pl.Amount
	trx.GoldNominal = goldNominal
	trx.Type = typeTrxCredit
	trx.Status = statusTrxPending
	trx.Balance = int64(balance)
	trx.GoldBalance = float64(goldBalance)
	trx.TrxDate = pl.TrxDateTime
	trx.Description = pl.TrxDesc

	return nil
}
