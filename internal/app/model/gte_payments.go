package model

import "time"

type GtePayment struct {
	ID          int64     `json:"id"`
	AccountId   int64     `json:"accountId"`
	TrxId       string    `json:"trxId"`
	GoldAmount  float64   `json:"goldAmount"`
	TrxAmount   int64     `json:"trxAmount"`
	BriUpdated  bool      `json:"briUpdated"`
	PdsNotified bool      `json:"pdsNotified"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	Account     Account
}
