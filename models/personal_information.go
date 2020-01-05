package models

import "time"

// PersonalInformation is a struct to store personal info data
type PersonalInformation struct {
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	HandPhoneNumber string    `json:"handPhoneNumber"`
	Email           string    `json:"email"`
	Npwp            string    `json:"npwp"`
	Nik             string    `json:"nik"`
	BirthPlace      string    `json:"birthPlace"`
	BirthDate       string    `json:"birthDate"`
	Nationality     string    `json:"nationality"`
	Sex             int64     `json:"sex"`
	Education       int64     `json:"education"`
	MaritalStatus   int64     `json:"maritalStatus"`
	MotherName      string    `json:"motherName"`
	HomePhoneArea   string    `json:"homePhoneArea"`
	HomePhoneNumber string    `json:"homePhoneNumber"`
	HomeStatus      int64     `json:"homeStatus"`
	AddressLine1    string    `json:"addressLine1"`
	AddressLine2    string    `json:"addressLine2"`
	AddressLine3    string    `json:"addressLine3"`
	Zipcode         string    `json:"zipcode"`
	AddressCity     string    `json:"addressCity"`
	StayedSince     string    `json:"stayedSince"`
	Child           int64     `json:"child"`
	CreatedAt       time.Time `json:"createdAt"`
	CpdatedAt       time.Time `json:"updatedAt"`
}
