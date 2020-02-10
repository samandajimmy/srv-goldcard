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

// Billing is a struct to store billing data
type Billing struct {
	AccountId   int64     `json:"accountId"`
	Amount      int64     `json:"amount"`
	GoldAmount  float64   `json:"goldAmount"`
	BillingDate time.Time `json:"billingDate"`
	DepthAmount int64     `json:"depthAmount"`
	DepthGold   float64   `json:"depthGold"`
	STL         int64     `json:"stl"`
	DepthSTL    int64     `json:"depthStl"`
	CreatedAt   time.Time `json:"createdAt"`
	Account     Account   `json:"account"`
}

// BillingTransaction is a struct to store billing transaction data
type BillingTransaction struct {
	TrxId       int64       `json:"trxId"`
	BillId      int64       `json:"billId"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	CreatedAt   time.Time   `json:"createdAt"`
	Transaction Transaction `json:"transaction"`
	Billing     Billing     `json:"billing"`
}

// BillingPayment is a struct to store billing payment data
type BillingPayment struct {
	TrxId       int64       `json:"trxId"`
	BillId      int64       `json:"billId"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	CreatedAt   time.Time   `json:"createdAt"`
	Transaction Transaction `json:"transaction"`
	Billing     Billing     `json:"billing"`
}

// MappingTransactions is a struct to mapping transactions data
func (trx *Transaction) MappingTransactions(c echo.Context, pl PayloadBRIPendingTransactions, trans Transaction, refTrxPg string, stl int64) error {
	goldNominal := trx.Account.Card.ConvertMoneyToGold(pl.Amount, stl)
	balance := trx.Account.Card.CardLimit - pl.Amount
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
