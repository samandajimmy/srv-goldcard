package model

import (
	"time"
)

var (
	// LimitFiftyMillions is a param to store value if increase limit inquiry is above 50 millions rupiah
	LimitFiftyMillions int64 = 50000000

	// MinIncreaseLimit is a param to store value of minimum increase limit goldcard
	MinIncreaseLimit int64 = 1000000

	// LimitUpdateStatusInquired to store a inquired update limit status
	LimitUpdateStatusInquired = "inquired"

	// LimitUpdateStatusPending to store a pending update limit status
	LimitUpdateStatusPending = "pending"

	// LimitUpdateStatusApplied to store a apllied update limit status
	LimitUpdateStatusApplied = "applied"
)

// LimitUpdate is a struct to store historical card limit update data
type LimitUpdate struct {
	ID               int64     `json:"id"`
	AccountID        int64     `json:"accountId"`
	RefId            string    `json:"refId"`
	AppliedLimitDate time.Time `json:"appliedLimitDate"`
	CardLimit        int64     `json:"cardLimit"`
	GoldLimit        float64   `json:"goldLimit"`
	StlLimit         int64     `json:"stlLimit"`
	Status           string    `json:"status"`
	WithNpwp         bool      `json:"withNpwp"`
	UpdatedAt        time.Time `json:"updatedAt"`
	CreatedAt        time.Time `json:"createdAt"`
	Account          Account   `json:"account"`
}

// CardUpdateLimit is a struct to store oldcard & newcard data limit
type CardUpdateLimit struct {
	OldCard Card
	NewCard Card
	Account Account
}
