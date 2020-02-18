package models

import (
	"math"
	"time"

	"github.com/labstack/echo"
)

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

// BillingStatement is a struct to store response for billing inquiry
type BillingStatement struct {
	BillingPrintDate  string `json:"billingPrintDate"`
	BillingDueDate    string `json:"billingDueDate"`
	BillingMinPayment int64  `json:"billingMinPayment"`
	BillingAmount     int64  `json:"billingAmount"`
}

// MappingBillingAccount is a struct to mapping billing account data
func (bill *Billing) MappingAccountNumberToBilling(c echo.Context, pl PayloadAccNumber) error {
	bill.Account.AccountNumber = pl.AccountNumber

	return nil
}

func (bill *Billing) MapBillingStatementResponse(c echo.Context, dueDate int, minPay float64, billStmt *BillingStatement) error {

	billStmt.BillingAmount = bill.Amount
	billStmt.BillingDueDate = bill.BillingDate.AddDate(0, 0, dueDate).Format("2006-01-02")
	billStmt.BillingMinPayment = int64(math.Ceil(float64(bill.Amount) * minPay))
	billStmt.BillingPrintDate = bill.BillingDate.Format("2006-01-02")

	return nil
}
