package models

import "time"

// Card is a struct to store card data
type Card struct {
	ID          int64     `json:"id"`
	CardName    string    `json:"cardName"`
	CardNumber  string    `json:"cardNumber"`
	CardLimit   string    `json:"cardLimit"`
	ValidUntil  string    `json:"validUntil"`
	PinNumber   string    `json:"pinNumber"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
}
