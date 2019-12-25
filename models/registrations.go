package models

import "time"

// Registrations is a struct to store registration data
type Registrations struct {
	ID                    int64      `json:"id,omitempty"`
	FullName              string     `json:"full_name,omitempty"`
	NameOnCard            string     `json:"name_on_card,omitempty"`
	Gender                string     `json:"gender,omitempty"`
	NPWPNumer             string     `json:"npwp_number,omitempty"`
	IdentityNumber        string     `json:"identity_number,omitempty"`
	DOB                   *time.Time `json:"dob,omitempty"`
	POB                   string     `json:"pob,omitempty"`
	Email                 string     `json:"email,omitempty"`
	ResidenceStatus       string     `json:"residence_status,omitempty"`
	ResidenceAddress      string     `json:"residence_address,omitempty" validate:"required"`
	ResidencePhoneNumber  string     `json:"residence_phone_number,omitempty"`
	PhoneNumber           string     `json:"phone_number,omitempty" validate:"required"`
	LatestEducationDegree string     `json:"latest_education_degree,omitempty"`
	MotherName            string     `json:"mother_name,omitempty"`
	CreatedAT             *time.Time `json:"created_at,omitempty"`
	UpdatedAT             *time.Time `json:"updated_at,omitempty"`
}
