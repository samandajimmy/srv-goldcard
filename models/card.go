package models

import (
	"math"
	"strconv"
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
	// defMoneyTaken to store a value of convertion gram to rupiah 94%
	defMoneyTaken float64 = 0.94
	// defCeilingLimit to store a value of ceiling limit 80% of defMoneyTaken
	defCeilingLimit float64 = 0.8
	// CardStatusActive to store a value of card status active
	CardStatusActive string = "active"
	// RequestPathCardBlock to store path BRI endpoint for block card
	RequestPathCardBlock string = "/card/block"
	// RequestPathCardStolen to store path BRI endpoint for stolen card
	RequestPathCardStolen string = "/card/stolen"
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
	ActivatedDate       time.Time `json:"activatedDate"`
}

// ConvertMoneyToGold to convert rupiah into gram
func (c *Card) ConvertMoneyToGold(money int64, stl int64) float64 {
	moneyFloat := float64(money)
	stlFloat := float64(stl)
	gold := (CustomRound("ceil", moneyFloat, 1000) / defMoneyTaken) / stlFloat

	return CustomRound("round", gold, 10000)

}

// SetGoldLimit a function to set gold limit in gram when open gte and change limit process
// by dividing nominal and gold price ceiling, gold price ceiling is 80% of appraised gold price at 94%
func (c *Card) SetGoldLimit(money int64, stl int64) float64 {
	moneyFloat := float64(money)
	stlFloat := float64(stl)
	goldPriceCeiling := defMoneyTaken * (defCeilingLimit * stlFloat)
	gold := moneyFloat / goldPriceCeiling

	return CustomRound("ceil", gold, 10000)
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
	ReportingDate int64  `json:"reportingDate"`
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
	reportDt := strconv.Itoa(int(briCardBlockStatus.ReportingDate))
	cs.Reason = pl.Reason
	cs.ReasonCode = pl.ReasonCode
	cs.CardID = card.ID
	cs.IsReactivated = IsReactivatedNo
	cs.BlockedDate, _ = time.Parse(DateTimeFormat, reportDt)

	return nil
}

// CardBalance is a struct to store card balance detail
type CardBalance struct {
	CurrGoldLimit float64
	CurrStl       int64
}
