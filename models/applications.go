package models

import "time"

// Registrations is a struct to store registration data
type Applications struct {
	ID                int64      `json:"id,omitempty"`
	ApplicationNumber string     `json:"application_number,omitempty"`
	CardLimit         string     `json:"card_limit,omitempty"`
	Status            string     `json:"status,omitempty"`
	KTP               string     `json:"ktp,omitempty"`
	NPWP              string     `json:"npwp,omitempty"`
	SavingAccount     string     `json:"saving_account,omitempty"`
	CreatedAT         *time.Time `json:"created_at,omitempty"`
	UpdatedAT         *time.Time `json:"updated_at,omitempty"`
}
