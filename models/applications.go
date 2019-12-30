package models

import "time"

// Registrations is a struct to store registration data
type Applications struct {
	ID                int64      `json:"id"`
	ApplicationNumber string     `json:"applicationNumber" validate:"required"`
	CardLimit         string     `json:"cardLimit"`
	Status            string     `json:"status"`
	KTP               string     `json:"ktp"`
	NPWP              string     `json:"npwp"`
	SavingAccount     string     `json:"savingAccount" validate:"required"`
	CreatedAT         *time.Time `json:"createdAt"`
	UpdatedAT         *time.Time `json:"updatedAt"`
}
