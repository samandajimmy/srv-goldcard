package models

import (
	"math"
	"time"
)

var IsReactivatedNo = "no"

const (
	// DecreasedLimit to store a value of decreasing limit
	DecreasedLimit float64 = 0.0115
	// MinEffBalance to store a value of minimum effective balance
	MinEffBalance float64 = 0.1000
	// ReservedLockBalance to store a value of additional lock gold balance when updating gold limit
	ReservedLockLimitBalance float64 = 0.5000

	defMoneyTaken float64 = 0.94

	cardStatusActive string = "active"
	// RequestPathCardBlock to store path BRI endpoint for block card
	RequestPathCardBlock string = "/v1/cobranding/card/block"
	// RequestPathCardStolen to store path BRI endpoint for stolen card
	RequestPathCardStolen string = "/v1/cobranding/card/stolen"
	// ReasonCodeStolen reason code stolen
	ReasonCodeStolen string = "stolen"
	// CardStatusInactive card status inactive
	CardStatusInactive string = "inactive"
)

// Card is a struct to store card data
type Card struct {
	ID                  int64     `json:"id"`
	CardName            string    `json:"cardName"`
	CardNumber          string    `json:"cardNumber"`
	CardLimit           int64     `json:"cardLimit"`
	GoldLimit           float64   `json:"goldLimit"`
	StlLimit            int64     `json:"stlLimit"`
	ValidUntil          string    `json:"validUntil"`
	PinNumber           string    `json:"pinNumber"`
	Description         string    `json:"description"`
	Balance             int64     `json:"balance"`
	GoldBalance         float64   `json:"goldBalance"`
	StlBalance          int64     `json:"stlBalance"`
	Status              string    `json:"status"`
	EncryptedCardNumber string    `json:"encryptedCardNumber"`
	UpdatedAt           time.Time `json:"updatedAt"`
	CreatedAt           time.Time `json:"createdAt"`
}

// ConvertMoneyToGold to convert rupiah into gram
func (c *Card) ConvertMoneyToGold(money int64, stl int64) float64 {
	moneyFloat := float64(money)
	stlFloat := float64(stl)
	gold := (CustomRound("ceil", moneyFloat, 1000) / defMoneyTaken) / stlFloat

	return CustomRound("round", gold, 10000)

}

// SetGoldLimit a function to set gold limit in gram
// when updating gold limit, gold limit is added with ReservedLockLimitBalance
func (c *Card) SetGoldLimit(money int64, stl int64) float64 {
	return c.ConvertMoneyToGold(money, stl) + ReservedLockLimitBalance
}

// SetCardLimit a function to set card limit in rupiah
// when updateing card limit, gold limit is subtracted by ReservedLockLimitBalance first, then convert to rupiah
func (c *Card) SetCardLimit(stl int64) error {
	// round down to nearest 10.000s
	c.CardLimit = int64(math.Floor((c.GoldLimit-ReservedLockLimitBalance)*float64(stl)*defMoneyTaken/10000)) * 10000
	c.StlLimit = stl

	return nil
}

// BRICardBlockStatus to store response BRI card block status
type BRICardBlockStatus struct {
	ReportingDate string `json:"reportingDate"`
	ReportDesc    string `json:"reportDesc"`
}

// CardStatuses is a struct to store card statuses
type CardStatuses struct {
	ID              int64     `json:"id"`
	Reason          string    `json:"reason"`
	ReasonCode      string    `json:"reasonCode"`
	IsReactivated   string    `json:"isReactivated"`
	CardID          int64     `json:"cardId"`
	BlockedDate     time.Time `json:"blockedDate"`
	ReactivatedDate time.Time `json:"reactivatedDate"`
	UpdatedAt       time.Time `json:"updatedAt"`
	CreatedAt       time.Time `json:"createdAt"`
}

func (cs *CardStatuses) MappingBlockCard(briCardBlockStatus BRICardBlockStatus, pl PayloadCardBlock, card Card) error {
	cs.Reason = briCardBlockStatus.ReportDesc
	cs.ReasonCode = pl.ReasonCode
	cs.CardID = card.ID
	cs.IsReactivated = IsReactivatedNo
	cs.BlockedDate, _ = time.Parse(DateTimeFormat, briCardBlockStatus.ReportingDate)

	if pl.Reason != "" {
		cs.Reason = pl.Reason
	}

	return nil
}
