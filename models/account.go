package models

import "time"

// Account is a struct to store account data
type Account struct {
	ID                    int64     `json:"id"`
	CIF                   string    `json:"cif"`
	AccountNumber         string    `json:"accountNumber"`
	BrixKey               string    `json:"brixKey"`
	CardLimit             int64     `json:"cardLimit"`
	Status                string    `json:"status"`
	BankID                int64     `json:"bankId"`
	CardID                int64     `json:"cardId"`
	ApplicationID         int64     `json:"applicationId"`
	PersonalInformationID int64     `json:"personalInformationId"`
	OccupationID          int64     `json:"occupationId"`
	EmergencyContactID    int64     `json:"emergencyContactId"`
	CorrespondenceID      int64     `json:"correspondenceId"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}
