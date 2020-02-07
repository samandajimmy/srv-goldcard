package models

import (
	"time"

	"github.com/labstack/echo"
)

const (
	typeTrxCredit string = "credit"

	statusTrxPending string = "pending"
)

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
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	Account     Account   `json:"account"`
}

type Billings struct {
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

type BillingTransactions struct {
	TrxId       int64       `json:"trxId"`
	BillId      int64       `json:"billId"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	CreatedAt   time.Time   `json:"createdAt"`
	Transaction Transaction `json:"transaction"`
	Billings    Billings    `json:"billings"`
}

type BillingPayments struct {
	TrxId       int64       `json:"trxId"`
	BillId      int64       `json:"billId"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	CreatedAt   time.Time   `json:"createdAt"`
	Transaction Transaction `json:"transaction"`
	Billings    Billings    `json:"billings"`
}

func (trx *Transaction) MappingTransactions(c echo.Context, pl PayloadBRIPendingTransactions, trans Transaction, refTrxPg string, stl int64) error {
	goldNominal := trans.Account.Card.ConvertMoneyToGold(pl.Amount, stl)

	trx.AccountId = trans.Account.ID
	trx.RefTrxPgdn = refTrxPg
	trx.RefTrx = pl.TransactionId
	trx.Nominal = pl.Amount
	trx.GoldNominal = goldNominal
	trx.Type = typeTrxCredit
	trx.Status = statusTrxPending
	trx.Balance = int64(trans.Account.Card.CardLimit - pl.Amount)
	trx.GoldBalance = float64(trans.Account.Card.GoldLimit - goldNominal)
	trx.TrxDate = pl.TrxDateTime

	return nil
}
