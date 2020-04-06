package models

import (
	"reflect"
	"time"

	"github.com/google/uuid"
)

const (
	// TypeTrxCredit trx type credit
	TypeTrxCredit string = "credit"
	// TypeTrxDebit trx type debit
	TypeTrxDebit string = "debit"
	// BillTrxUnpaid billing status unpaid
	BillTrxUnpaid string = "unpaid"
	// BillTrxPaid billing status paid
	BillTrxPaid string = "paid"
	// SourceCore source for core
	SourceCore       string = "core"
	statusTrxPending string = "pending"
	statusTrxPosted  string = "posted"
	methodTrxPayment string = "payment"
)

// Transaction is a struct to store transaction data
type Transaction struct {
	ID              int64            `json:"id"`
	AccountId       int64            `json:"accountId"`
	RefTrxPgdn      string           `json:"refTrxPgdn"`
	RefTrx          string           `json:"refTrx"`
	Nominal         int64            `json:"nominal"`
	GoldNominal     float64          `json:"goldNominal"`
	Type            string           `json:"type"`
	Status          string           `json:"status"`
	Balance         int64            `json:"balance"`
	GoldBalance     float64          `json:"goldBalance"`
	Methods         string           `json:"methods"`
	TrxDate         string           `json:"trxDate"`
	Description     string           `json:"description"`
	CompareID       string           `json:"compareId"`
	TransactionID   string           `json:"transactionId"`
	UpdatedAt       time.Time        `json:"updatedAt"`
	CreatedAt       time.Time        `json:"createdAt"`
	Account         Account          `json:"account"`
	BillingPayments []BillingPayment `json:"billingPayments" pg:"-"`
	CurrStl         int64            `json:"currStl" pg:"-"`
}

// ListTrx struct to store list history transactions
type ListTrx struct {
	// nolint
	tableName struct{} `pg:"transactions"`

	RefTrx      string `json:"refTrx"`
	Nominal     int64  `json:"nominal"`
	TrxDate     string `json:"trxDate"`
	Description string `json:"description"`
}

// ResponseListTrx struct to store response history transactions
type ResponseListTrx struct {
	IsLastPage bool      `json:"isLastPage"`
	ListTrx    []ListTrx `json:"listHistoryTransactions"`
}

type BRICardBalance struct {
	CurrentBalance  float64 `json:"currentBalance"`
	CreditLimit     float64 `json:"creditLimit"`
	AvailableCredit int64   `json:"availableCredit"`
}

// MappingTrx is a struct to mapping trx data
func (trx *Transaction) MappingTrx(pl interface{}, trxType string, isTrx bool) error {
	// Generate ref transactions pegadaian
	refTrxPgdn, _ := uuid.NewRandom()
	// reflect payload interface
	r := reflect.ValueOf(pl)
	// init variables inside trx struct
	trx.Nominal = GetInterfaceValue(r, "PaymentAmount").(int64)
	trx.TransactionID = GetInterfaceValue(r, "TransactionId").(string)
	trx.RefTrx = GetInterfaceValue(r, "RefTrx").(string)
	trx.TrxDate = GetInterfaceValue(r, "PaymentDate").(string)
	trx.CompareID = GetInterfaceValue(r, "AuthCode").(string)
	trx.Description = GetInterfaceValue(r, "PaymentDesc").(string)

	// if no value on pl->RefTrx
	if trx.RefTrx == "" {
		trx.RefTrx = GetInterfaceValue(r, "RefID").(string)
	}

	// if no value on pl->PaymentDate
	if trx.TrxDate == "" {
		trx.TrxDate = GetInterfaceValue(r, "TrxDateTime").(string)
	}

	// if no value on pl->TrxDateTime
	if trx.TrxDate == "" {
		trx.TrxDate = time.Now().Format(time.RFC3339)
	}

	// if no value on pl->PaymentDesc
	if trx.Description == "" {
		// TODO: we should have default payment description
		trx.Description = ""
	}

	// if no value on pl->PaymentAmount
	if trx.Nominal == 0 {
		trx.Nominal = GetInterfaceValue(r, "Amount").(int64)
	}

	// if its transaction data
	if isTrx {
		trx.Balance = trx.Account.Card.Balance - trx.Nominal
		trx.Status = statusTrxPending
	}

	// if its payment transaction data
	if !isTrx {
		trx.Balance = trx.Account.Card.Balance + trx.Nominal
		trx.Status = statusTrxPosted
		trx.Methods = methodTrxPayment
		trx.BillingPayments = []BillingPayment{
			BillingPayment{Source: GetInterfaceValue(r, "Source").(string)},
		}
	}

	trx.GoldNominal = trx.Account.Card.ConvertMoneyToGold(trx.Nominal, trx.CurrStl)
	trx.GoldBalance = trx.Account.Card.ConvertMoneyToGold(trx.Balance, trx.CurrStl)
	trx.RefTrxPgdn = refTrxPgdn.String()
	trx.Type = trxType

	return nil
}
