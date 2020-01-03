package models

// PayloadList a struct to store all payload for a list response
type PayloadList struct {
	CIF             string `json:"cif,omitempty"`
	Status          string `json:"status,omitempty"`
	StartDate       string `json:"startDate,omitempty" validate:"dateString"`
	EndDate         string `json:"endDate,omitempty"`
	ProductCode     string `json:"productCode,omitempty"`
	TransactionType string `json:"transactionType,omitempty"`
	Page            int64  `json:"page,omitempty" validate:"min=1"`
	Limit           int64  `json:"limit,omitempty" validate:"min=1"`
}

// PayloadGetAddress a struct to store all payload for a get address
type PayloadGetAddress struct {
	PhoneNumber string `json:"phoneNumber,omitempty" validate:"required"`
}

// PayloadRegistration a struct to store all payload for registration
type PayloadRegistration struct {
	CIF         string `json:"cif" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
}

// PayloadPersonalInformation a struct to store all payload for a payload personal information
type PayloadPersonalInformation struct {
	FirstName            string `json:"firstName" validate:"required"`
	LastName             string `json:"lastName" validate:"required"`
	CardName             string `json:"cardName" validate:"required"`
	Npwp                 string `json:"npwp" validate:"required"`
	Nik                  string `json:"nik" validate:"required"`
	BirthPlace           string `json:"birthPlace" validate:"required"`
	BirthDate            string `json:"birthDate" validate:"required"`
	AddressLine1         string `json:"addressLine1" validate:"required"`
	AddressLine2         string `json:"addressLine2" validate:"required"`
	AddressLine3         string `json:"addressLine3" validate:"required"`
	Sex                  string `json:"sex" validate:"required"`
	HomeStatus           string `json:"homeStatus" validate:"required"`
	AddressCity          string `json:"addressCity" validate:"required"`
	Nationality          string `json:"nationality" validate:"required"`
	StayedSince          string `json:"stayedSince" validate:"required"`
	Education            string `json:"education" validate:"required"`
	Zipcode              string `json:"zipcode" validate:"required"`
	MaritalStatus        string `json:"maritalStatus" validate:"required"`
	MotherName           string `json:"motherName" validate:"required"`
	HandPhoneNumber      string `json:"handPhoneNumber" validate:"required"`
	HomePhoneArea        string `json:"homePhoneArea" validate:"required"`
	HomePhoneNumber      string `json:"homePhoneNumber" validate:"required"`
	Email                string `json:"email" validate:"required"`
	JobBidangUsaha       string `json:"jobBidangUsaha" validate:"required"`
	JobSubBidangUsaha    string `json:"jobSubBidangUsaha" validate:"required"`
	JobCategory          string `json:"jobCategory" validate:"required"`
	JobStatus            string `json:"jobStatus" validate:"required"`
	TotalEmployee        string `json:"totalEmployee" validate:"required"`
	Company              string `json:"company" validate:"required"`
	JobTitle             string `json:"jobTitle" validate:"required"`
	WorkSince            string `json:"workSince" validate:"required"`
	OfficeAddress1       string `json:"officeAddress1" validate:"required"`
	OfficeAddress2       string `json:"officeAddress2" validate:"required"`
	OfficeAddress3       string `json:"officeAddress3" validate:"required"`
	OfficeZipcode        string `json:"officeZipcode" validate:"required"`
	OfficeCity           string `json:"officeCity" validate:"required"`
	OfficePhone          string `json:"officePhone" validate:"required"`
	Income               string `json:"income" validate:"required"`
	Child                string `json:"child" validate:"required"`
	EmergencyName        string `json:"emergencyName" validate:"required"`
	EmergencyRelation    string `json:"emergencyRelation" validate:"required"`
	EmergencyAddress1    string `json:"emergencyAddress1" validate:"required"`
	EmergencyAddress2    string `json:"emergencyAddress2" validate:"required"`
	EmergencyAddress3    string `json:"emergencyAddress3" validate:"required"`
	EmergencyCity        string `json:"emergencyCity" validate:"required"`
	EmergencyPhoneNumber string `json:"emergencyPhoneNumber" validate:"required"`
	ProductRequest       string `json:"productRequest" validate:"required"`
	BillingCycle         string `json:"billingCycle" validate:"required"`
	FotoKtpBase64        string `json:"fotoKtpBase64" validate:"required"`
	FotoNPWPBase64       string `json:"fotoNPWPBase64" validate:"required"`
	FotoDiriBase64       string `json:"fotoDiriBase64" validate:"required"`
}
