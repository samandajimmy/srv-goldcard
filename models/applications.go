package models

import "time"

// Applications is a struct to store application data
type Applications struct {
	ID                int64      `json:"id"`
	ApplicationNumber string     `json:"applicationNumber" validate:"required"`
	CardLimit         string     `json:"cardLimit"`
	Status            string     `json:"status"`
	KTP               string     `json:"ktp"`
	NPWP              string     `json:"npwp"`
	SavingAccount     string     `json:"savingAccount" validate:"required"`
	CreatedAt         *time.Time `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
}
