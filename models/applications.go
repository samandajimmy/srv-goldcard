package models

import "time"

// Applications is a struct to store application data
type Applications struct {
	ID                int64      `json:"id"`
	ApplicationNumber string     `json:"applicationNumber" validate:"required"`
	Status            string     `json:"status"`
	KtpImageBase64    string     `json:"ktpImageBase64"`
	NpwpImageBase64   string     `json:"npwpImageBase64"`
	SelfieImageBase64 string     `json:"selfieImageBase64"`
	SavingAccount     string     `json:"savingAccount" validate:"required"`
	CreatedAt         *time.Time `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
}
