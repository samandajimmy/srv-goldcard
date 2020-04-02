package models

import (
	"time"
)

const (
	// DecreasedLimit to store a value of decreasing limit
	DecreasedLimit float64 = 0.0115
	// MinEffBalance to store a value of minimum effective balance
	MinEffBalance float64 = 0.1000
	// ReservedLockBalance to store a value of additional lock gold balance when updating gold limit
	ReservedLockLimitBalance float64 = 0.5000

	defMoneyTaken float64 = 0.94

	cardStatusActive string = "active"
)

// Card is a struct to store card data
type Card struct {
	ID                int64     `json:"id"`
	CardName          string    `json:"cardName"`
	CardNumber        string    `json:"cardNumber"`
	CardLimit         int64     `json:"cardLimit"`
	GoldLimit         float64   `json:"goldLimit"`
	ReservedGoldLimit float64   `json:"reservedGoldLimit"`
	StlLimit          int64     `json:"stlLimit"`
	ValidUntil        string    `json:"validUntil"`
	PinNumber         string    `json:"pinNumber"`
	Description       string    `json:"description"`
	Balance           int64     `json:"balance"`
	GoldBalance       float64   `json:"goldBalance"`
	StlBalance        int64     `json:"stlBalance"`
	Status            string    `json:"status"`
	UpdatedAt         time.Time `json:"updatedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}

// ConvertMoneyToGold to convert rupiah into gram
func (c *Card) ConvertMoneyToGold(money int64, stl int64) float64 {
	moneyFloat := float64(money)
	stlFloat := float64(stl)
	gold := (CustomRound("ceil", moneyFloat, 1000) / defMoneyTaken) / stlFloat

	return CustomRound("round", gold, 10000)

}

// SetSubmissionGoldLimit a function to set gold limit in gram
// when updating gold limit, gold limit is added with ReservedLockLimitBalance
func (c *Card) SetSubmissionGoldLimit(money int64, stl int64) float64 {
	return c.ConvertMoneyToGold(money, stl) + ReservedLockLimitBalance
}

// SetTransactionGoldLimit a function to set gold limit in gram
// when updating gold limit, gold limit is added with reserved gold limit value in cards table
func (c *Card) SetTransactionGoldLimit(money int64, stl int64) float64 {
	return c.ConvertMoneyToGold(money, stl) + c.ReservedGoldLimit
}
