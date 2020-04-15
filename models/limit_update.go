package models

import (
	"time"
)

// LimitUpdate is a struct to store historical card limit update data
type LimitUpdate struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"accountId"`
	RefId     string    `json:"refId"`
	LimitDate time.Time `json:"limitDate"`
	CardLimit int64     `json:"cardLimit"`
	GoldLimit float64   `json:"goldLimit"`
	StlLimit  int64     `json:"stlLimit"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
	Account   Account   `json:"account"`
}

// CardUpdateLimit is a struct to store oldcard & newcard data limit
type CardUpdateLimit struct {
	OldCard Card
	NewCard Card
	Account Account
}
