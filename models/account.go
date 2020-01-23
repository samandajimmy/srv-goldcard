package models

import (
	"time"

	"github.com/labstack/echo"
)

// Account is a struct to store account data
type Account struct {
	ID                    string              `json:"id"`
	CIF                   string              `json:"cif"`
	ProductRequest        string              `json:"productRequest"`
	BillingCycle          int64               `json:"billingCycle"`
	CardDeliver           int64               `json:"cardDeliver"`
	BrixKey               string              `json:"brixkey" pg:"brixkey"`
	BranchCode            string              `json:"branchCode"`
	Status                string              `json:"status"`
	BankID                int64               `json:"bankId"`
	CardID                int64               `json:"cardId"`
	ApplicationID         int64               `json:"applicationId"`
	PersonalInformationID int64               `json:"personalInformationId"`
	OccupationID          int64               `json:"occupationId"`
	EmergencyContactID    int64               `json:"emergencyContactId"`
	CorrespondenceID      int64               `json:"correspondenceId"`
	CreatedAt             time.Time           `json:"createdAt"`
	UpdatedAt             time.Time           `json:"updatedAt"`
	Bank                  Bank                `json:"bank"`
	Card                  Card                `json:"card"`
	Application           Applications        `json:"application"`
	PersonalInformation   PersonalInformation `json:"personalInformation"`
	Occupation            Occupation          `json:"occupation"`
	EmergencyContact      EmergencyContact    `json:"emergencyContact"`
	Correspondence        Correspondence      `json:"correspondence"`
}

// MappingRegistrationData a function to map all data registration
func (acc *Account) MappingRegistrationData(c echo.Context, pl PayloadPersonalInformation) error {
	acc.Card.CardName = pl.CardName

	acc.Occupation.JobBidangUsaha = pl.JobBidangUsaha
	acc.Occupation.JobSubBidangUsaha = pl.JobSubBidangUsaha
	acc.Occupation.JobCategory = pl.JobCategory
	acc.Occupation.JobStatus = pl.JobStatus
	acc.Occupation.TotalEmployee = pl.TotalEmployee
	acc.Occupation.Company = pl.Company
	acc.Occupation.JobTitle = pl.JobTitle
	acc.Occupation.WorkSince = pl.WorkSince
	acc.Occupation.OfficeAddress1 = pl.OfficeAddress1
	acc.Occupation.OfficeAddress2 = pl.OfficeAddress2
	acc.Occupation.OfficeAddress3 = pl.OfficeAddress3
	acc.Occupation.OfficeZipcode = pl.OfficeZipcode
	acc.Occupation.OfficeCity = pl.OfficeCity
	acc.Occupation.OfficePhone = pl.OfficePhone
	acc.Occupation.Income = pl.Income

	acc.Application.KtpImageBase64 = pl.KtpImageBase64
	acc.Application.NpwpImageBase64 = pl.NpwpImageBase64
	acc.Application.SelfieImageBase64 = pl.SelfieImageBase64

	acc.PersonalInformation.FirstName = pl.FirstName
	acc.PersonalInformation.LastName = pl.LastName
	acc.PersonalInformation.HandPhoneNumber = pl.HandPhoneNumber
	acc.PersonalInformation.Email = pl.Email
	acc.PersonalInformation.Npwp = pl.Npwp
	acc.PersonalInformation.Nik = pl.Nik
	acc.PersonalInformation.BirthPlace = pl.BirthPlace
	acc.PersonalInformation.BirthDate = pl.BirthDate
	acc.PersonalInformation.Nationality = pl.Nationality
	acc.PersonalInformation.Sex = acc.PersonalInformation.GetSex(pl.Sex)
	acc.PersonalInformation.Education = pl.Education
	acc.PersonalInformation.MaritalStatus = pl.MaritalStatus
	acc.PersonalInformation.MotherName = pl.MotherName
	acc.PersonalInformation.HomePhoneArea = pl.HomePhoneArea
	acc.PersonalInformation.HomePhoneNumber = pl.HomePhoneNumber
	acc.PersonalInformation.HomeStatus = pl.HomeStatus
	acc.PersonalInformation.AddressLine1 = pl.AddressLine1
	acc.PersonalInformation.AddressLine2 = pl.AddressLine2
	acc.PersonalInformation.AddressLine3 = pl.AddressLine3
	acc.PersonalInformation.Zipcode = pl.Zipcode
	acc.PersonalInformation.AddressCity = pl.AddressCity
	acc.PersonalInformation.StayedSince = pl.StayedSince
	acc.PersonalInformation.Child = pl.Child
	acc.PersonalInformation.RelativePhoneNumber = pl.RelativePhoneNumber

	acc.ProductRequest = pl.ProductRequest
	acc.BillingCycle = pl.BillingCycle
	acc.CardDeliver = pl.CardDeliver

	return nil
}

// MappingAddressData a function to map all data addresses
func (acc *Account) MappingAddressData(c echo.Context, pl PayloadAddress) error {
	acc.Correspondence.AddressLine1 = pl.AddressLine1
	acc.Correspondence.AddressLine2 = pl.AddressLine2
	acc.Correspondence.AddressLine3 = pl.AddressLine3
	acc.Correspondence.AddressCity = pl.AddressCity

	return nil
}

// Bank is a struct to strore bank data
type Bank struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Correspondence is a struct to store correspondence data
type Correspondence struct {
	ID           int64     `json:"id"`
	AddressLine1 string    `json:"addressLine1" pg:"address_line_1"`
	AddressLine2 string    `json:"addressLine2" pg:"address_line_2"`
	AddressLine3 string    `json:"addressLine3" pg:"address_line_3"`
	AddressCity  string    `json:"addressCity"`
	Zipcode      string    `json:"zipcode"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// EmergencyContact is a struct to store emergencyContact data
type EmergencyContact struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Relation     int64     `json:"relation"`
	PhoneNumber  string    `json:"phoneNumber"`
	AddressLine1 string    `json:"addressLine1" pg:"address_line_1"`
	AddressLine2 string    `json:"addressLine2" pg:"address_line_2"`
	AddressLine3 string    `json:"addressLine3" pg:"address_line_3"`
	AddressCity  string    `json:"addressCity"`
	Zipcode      string    `json:"zipcode"`
	Type         string    `json:"type"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// AddressData is a struct to store address data
type AddressData struct {
	City        string `json:"city"`
	Province    string `json:"province"`
	Subdistrict string `json:"subdistrict"`
	Village     string `json:"village"`
}
