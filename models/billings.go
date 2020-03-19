package models

import (
	"time"

	"github.com/labstack/echo"
)

var (
	// BillUnpaid is to store var billing status unpaid
	BillUnpaid = "unpaid"
)

// Billing is a struct to store billing data
type Billing struct {
	ID             int64     `json:"id"`
	AccountId      int64     `json:"accountId"`
	Amount         int64     `json:"amount"`
	GoldAmount     float64   `json:"goldAmount"`
	BillingDate    time.Time `json:"billingDate"`
	BillingDueDate time.Time `json:"billingDueDate"`
	DebtAmount     int64     `json:"debtAmount"`
	DebtGold       float64   `json:"debtGold"`
	MinimumPayment float64   `json:"minimum_payment"`
	STL            int64     `json:"stl"`
	DebtSTL        int64     `json:"debtStl"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Account        Account   `json:"account"`
}

// BillingTransaction is a struct to store billing transaction data
type BillingTransaction struct {
	ID          int64       `json:"id"`
	TrxId       int64       `json:"trxId"`
	BillId      int64       `json:"billId"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	CreatedAt   time.Time   `json:"createdAt"`
	Transaction Transaction `json:"transaction"`
	Billing     Billing     `json:"billing"`
}

// BillingPayment is a struct to store billing payment data
type BillingPayment struct {
	ID        int64     `json:"id"`
	TrxId     int64     `json:"trxId"`
	BillId    int64     `json:"billId"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// BillingStatement is a struct to store response for billing inquiry
type BillingStatement struct {
	BillingPrintDate  string `json:"billingPrintDate"`
	BillingDueDate    string `json:"billingDueDate"`
	BillingMinPayment int64  `json:"billingMinPayment"`
	BillingAmount     int64  `json:"billingAmount"`
}

// PegadaianBilling is a struct to store pegadaian billings data
type PegadaianBilling struct {
	ID            int64     `json:"id"`
	RefID         string    `json:"refID"`
	FileName      string    `json:"fileName"`
	BillingDate   string    `json:"billingDate"`
	FileBase64    string    `json:"fileBase64"`
	FileExtension string    `json:"fileExtension"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// MappingBillingAccount is a struct to mapping billing account data
func (bill *Billing) MappingAccountNumberToBilling(c echo.Context, pl PayloadAccNumber) error {
	bill.Account.AccountNumber = pl.AccountNumber

	return nil
}

// MappingPegadaianBilling is a function to mapping pegadaian billing
func (pgdBil *PegadaianBilling) MappingPegadaianBilling(c echo.Context, pl PayloadBRIPegadaianBillings) error {
	pgdBil.RefID = pl.RefID
	pgdBil.BillingDate = pl.BillingDate
	pgdBil.FileName = pl.FileName
	pgdBil.FileBase64 = pl.FileBase64
	pgdBil.FileExtension = pl.FileExtension
	pgdBil.CreatedAt = time.Now()

	return nil
}
