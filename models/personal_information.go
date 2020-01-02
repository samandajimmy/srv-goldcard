package models

import "time"

// PersonalInformation is a struct to store personal info data
type PersonalInformation struct {
	ID                    int64     `json:"id"`
	FullName              string    `json:"fullName"`
	NameOnCard            string    `json:"nameOnCard"`
	Gender                string    `json:"gender"`
	NpwpNumber            string    `json:"npwpNumber"`
	IdentityNumber        string    `json:"identityNumber"`
	DOB                   time.Time `json:"dob"`
	POB                   string    `json:"pob"`
	Email                 string    `json:"email"`
	ResidenceDtatus       string    `json:"residenceStatus"`
	ResidenceSddress      string    `json:"residenceAddress"`
	ResidencePhoneNumber  string    `json:"residencePhoneNumber"`
	PhoneNumber           string    `json:"phoneNumber"`
	LatestEducationDegree string    `json:"latestEducationDegree"`
	MotherName            string    `json:"motherName"`
	CreatedAt             time.Time `json:"createdAt"`
	CpdatedAt             time.Time `json:"updatedAt"`
}
