package models

import (
	"gade/srv-goldcard/logger"
	"time"

	"github.com/labstack/echo"
)

var (
	AccStatusActive = "active"

	AccStatusInactive = "inactive"

	appendXCardNumber = "xxxxxx"

	RelationStr = map[int64]string{
		1:  "Suami/Istri",
		2:  "Anak",
		3:  "Adik",
		4:  "Kakak Kandung",
		5:  "Orang Tua",
		6:  "Saudara",
		7:  "HRD",
		8:  "Atasan",
		9:  "Lain-lain",
		10: "Applicant",
	}
)

// Account is a struct to store account data
type Account struct {
	ID                    int64               `json:"id"`
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
	AccountNumber         string              `json:"accountNumber"`
	CreatedAt             time.Time           `json:"createdAt"`
	UpdatedAt             time.Time           `json:"updatedAt"`
	Bank                  Bank                `json:"bank"`
	Card                  Card                `json:"card"`
	Application           Applications        `json:"application"`
	PersonalInformation   PersonalInformation `json:"personalInformation"`
	Occupation            Occupation          `json:"occupation"`
	EmergencyContact      EmergencyContact    `json:"emergencyContact"`
}

// MappingRegistrationData a function to map all data registration
func (acc *Account) MappingRegistrationData(pl PayloadPersonalInformation, addrData AddressData) error {
	acc.Card.CardName = pl.CardName

	acc.PersonalInformation.FirstName = pl.FirstName
	acc.PersonalInformation.LastName = pl.LastName
	acc.PersonalInformation.HandPhoneNumber = pl.HandPhoneNumber
	acc.PersonalInformation.Email = pl.Email
	acc.PersonalInformation.Nik = pl.Nik
	acc.PersonalInformation.BirthPlace = pl.BirthPlace
	acc.PersonalInformation.BirthDate = pl.BirthDate
	acc.PersonalInformation.Nationality = nationalityMap[pl.Nationality]
	acc.PersonalInformation.Sex = acc.PersonalInformation.GetSex(pl.Sex)
	acc.PersonalInformation.Education = pl.Education
	acc.PersonalInformation.MaritalStatus = pl.MaritalStatus
	acc.PersonalInformation.MotherName = pl.MotherName
	acc.PersonalInformation.HomeStatus = pl.HomeStatus
	acc.PersonalInformation.Zipcode = pl.Zipcode
	acc.PersonalInformation.AddressProvince = pl.Province
	acc.PersonalInformation.AddressCity = pl.AddressCity
	acc.PersonalInformation.AddressSubdistrict = pl.Subdistrict
	acc.PersonalInformation.AddressVillage = pl.Village
	acc.PersonalInformation.StayedSince = defStayedSince
	acc.PersonalInformation.Child = defChildNumber
	acc.PersonalInformation.RelativePhoneNumber = pl.RelativePhoneNumber

	// set addressData
	addrData, err := RemappAddress(addrData, 30)

	if err != nil {
		logger.Make(nil, nil).Debug(err)

		return err
	}

	acc.PersonalInformation.AddressLine1 = addrData.AddressLine1
	acc.PersonalInformation.AddressLine2 = addrData.AddressLine2
	acc.PersonalInformation.AddressLine3 = addrData.AddressLine3

	if acc.PersonalInformation.HandPhoneNumber != "" {
		acc.PersonalInformation.HomePhoneArea = pl.HomePhoneArea
	}

	// set home phone number
	acc.PersonalInformation.SetHomePhone()

	// set npwp
	acc.PersonalInformation.SetNPWP(pl.Npwp)

	// set default base64 to NPWP image if empty
	if pl.NpwpImageBase64 == "" {
		pl.NpwpImageBase64 = DefDocBase64
	}

	// application documents
	acc.Application.SetDocument(pl)

	return nil
}

// MappingCardActivationsData a function to map all data activations
func (acc *Account) MappingCardActivationsData(c echo.Context, pa PayloadActivations) error {
	acc.Card.CardNumber = pa.FirstSixDigits + appendXCardNumber + pa.LastFourDigits
	acc.Card.ValidUntil = pa.ExpDate
	acc.Application.Status = AppStatusActive
	acc.Card.Status = cardStatusActive
	acc.Status = AccStatusActive

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
	City         string `json:"city"`
	Province     string `json:"province"`
	Subdistrict  string `json:"subdistrict" pg:"sub_district"`
	Village      string `json:"village"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3"`
	Zipcode      string `json:"zipcode"`
}
