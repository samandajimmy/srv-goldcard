package models

import (
	"time"

	"github.com/labstack/echo"
)

const (
	typeTrxCredit    string = "credit"
	typeTrxDebit     string = "debit"
	statusTrxPending string = "pending"
	statusTrxPosted  string = "pending"
	methodTrxPayment string = "payment"
)

// Transaction is a struct to store transaction data
type Transaction struct {
	ID            int64     `json:"id"`
	AccountId     int64     `json:"accountId"`
	RefTrxPgdn    string    `json:"refTrxPgdn"`
	RefTrx        string    `json:"refTrx"`
	Nominal       int64     `json:"nominal"`
	GoldNominal   float64   `json:"goldNominal"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	Balance       int64     `json:"balance"`
	GoldBalance   float64   `json:"goldBalance"`
	Methods       string    `json:"methods"`
	TrxDate       string    `json:"trxDate"`
	Description   string    `json:"description"`
	CompareID     string    `json:"compareId"`
	TransactionID string    `json:"transactionId"`
	UpdatedAt     time.Time `json:"updatedAt"`
	CreatedAt     time.Time `json:"createdAt"`
	Account       Account   `json:"account"`
}

type BRICardBalance struct {
	CurrentBalance float64 `json:"currentBalance"`
	CreditLimit    float64 `json:"creditLimit"`
}

// MappingTransactions is a struct to mapping transactions data
func (trx *Transaction) MappingTransactions(c echo.Context, pl PayloadBRIPendingTransactions, trans Transaction, refTrxPg string, stl int64) error {
	goldNominal := trx.Account.Card.ConvertMoneyToGold(pl.Amount, stl)
	balance := trx.Account.Card.Balance - pl.Amount
	goldBalance := trx.Account.Card.ConvertMoneyToGold(balance, stl)

	trx.AccountId = trans.Account.ID
	trx.RefTrxPgdn = refTrxPg
	trx.TransactionID = pl.TransactionId
	trx.Nominal = pl.Amount
	trx.GoldNominal = goldNominal
	trx.Type = typeTrxCredit
	trx.Status = statusTrxPending
	trx.Balance = int64(balance)
	trx.GoldBalance = float64(goldBalance)
	trx.TrxDate = pl.TrxDateTime
	trx.Description = pl.TrxDesc
	trx.CompareID = pl.AuthCode

	return nil
}

// MappingTransactions is a struct to mapping payment transactions data
func (trx *Transaction) MappingPaymentTransactions(c echo.Context, pl PayloadBRIPaymentTransactions, trans Transaction, refTrxPg string, stl int64) error {
	goldNominal := trx.Account.Card.ConvertMoneyToGold(pl.PaymentAmount, stl)
	balance := trx.Account.Card.Balance + pl.PaymentAmount
	goldBalance := trx.Account.Card.ConvertMoneyToGold(balance, stl)

	trx.AccountId = trans.Account.ID
	trx.RefTrxPgdn = refTrxPg
	trx.RefTrx = pl.RefID
	trx.Nominal = pl.PaymentAmount
	trx.GoldNominal = goldNominal
	trx.Type = typeTrxDebit
	trx.Status = statusTrxPosted
	trx.Methods = methodTrxPayment
	trx.Balance = int64(balance)
	trx.GoldBalance = float64(goldBalance)
	trx.TrxDate = pl.PaymentDate
	trx.Description = pl.PaymentDesc

	return nil
}
