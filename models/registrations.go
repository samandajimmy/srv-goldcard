package models

import "time"

// Registrations is a struct to store registration data
type Registrations struct {
	ID                    int64      `json:"id"`
	FullName              string     `json:"fullName"`
	NameOnCard            string     `json:"nameOnCard"`
	Gender                string     `json:"gender"`
	NPWPNumer             string     `json:"npwpNumber"`
	IdentityNumber        string     `json:"identityNumber"`
	DOB                   *time.Time `json:"dob"`
	POB                   string     `json:"pob"`
	Email                 string     `json:"email"`
	ResidenceStatus       string     `json:"residenceStatus"`
	ResidenceAddress      string     `json:"residenceAddress" validate:"required"`
	ResidencePhoneNumber  string     `json:"residencePhoneNumber"`
	PhoneNumber           string     `json:"phoneNumber" validate:"required"`
	LatestEducationDegree string     `json:"latestEducationDegree"`
	MotherName            string     `json:"motherName"`
	CreatedAT             *time.Time `json:"createdAT"`
	UpdatedAT             *time.Time `json:"updatedAT"`
}
