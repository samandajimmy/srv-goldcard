package model

import "time"

// GoldPrice is a struct to store goldPrice data
type GoldPrice struct {
	ID        int64     `json:"id"`
	Price     float64   `json:"price"`
	ValidDate time.Time `json:"validDate"`
	CreatedAt time.Time `json:"createdAt"`
}
