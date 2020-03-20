package models

var (
	// UseExistingAddress is var to store use existing address status
	UseExistingAddress int64
)

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

// PayloadRegistration a struct to store all payload for registration
type PayloadRegistration struct {
	CIF               string `json:"cif" validate:"required"`
	HandPhoneNumber   string `json:"handPhoneNumber" validate:"required"`
	BranchCode        string `json:"branchCode" validate:"required"`
	ApplicationNumber string `json:"applicationNumber"`
}

// PayloadAppNumber a struct to store all payload for final registration
type PayloadAppNumber struct {
	ApplicationNumber string `json:"applicationNumber" validate:"required"`
}

// PayloadSavingAccount a struct to store all payload for saving account
type PayloadSavingAccount struct {
	ApplicationNumber string `json:"applicationNumber" validate:"required"`
	AccountNumber     string `json:"accountNumber" validate:"required"`
}

// PayloadCardLimit a struct to store all payload for card limit
type PayloadCardLimit struct {
	ApplicationNumber string `json:"applicationNumber" validate:"required"`
	CardLimit         int64  `json:"cardLimit" validate:"required"`
}

// PayloadAddress a struct to store all payload for user address
type PayloadAddress struct {
	ApplicationNumber string `json:"applicationNumber" validate:"required"`
	IsNew             int64  `json:"isNew" validate:"min=0,max=1"`
	AddressLine1      string `json:"addressLine1" validate:"required_with=IsNew"`
	AddressLine2      string `json:"addressLine2"`
	AddressLine3      string `json:"addressLine3"`
	AddressCity       string `json:"addressCity" validate:"required_with=IsNew"`
	Province          string `json:"province" validate:"required_with=IsNew"`
	Subdistrict       string `json:"subdistrict" validate:"required_with=IsNew"`
	Village           string `json:"village" validate:"required_with=IsNew"`
}

// PayloadPersonalInformation a struct to store all payload for a payload personal information
type PayloadPersonalInformation struct {
	ApplicationNumber    string `json:"applicationNumber,omitempty" validate:"required"`
	FirstName            string `json:"firstName" validate:"required"`
	LastName             string `json:"lastName" validate:"required"`
	CardName             string `json:"cardName" validate:"required"`
	Npwp                 string `json:"npwp"`
	Nik                  string `json:"nik" validate:"required"`
	BirthPlace           string `json:"birthPlace" validate:"required"`
	BirthDate            string `json:"birthDate" validate:"required"`
	AddressLine1         string `json:"addressLine1" validate:"required" pg:"address_line_1"`
	AddressLine2         string `json:"addressLine2" pg:"address_line_2"`
	AddressLine3         string `json:"addressLine3" pg:"address_line_3"`
	Sex                  int64  `json:"sex" validate:"required" pg:"-"`
	SexString            string `json:"sexString,omitempty" pg:"sex"`
	HomeStatus           int64  `json:"homeStatus" validate:"required"`
	AddressCity          string `json:"addressCity" validate:"required"`
	Nationality          string `json:"nationality" validate:"required"`
	StayedSince          string `json:"stayedSince"`
	Education            int64  `json:"education" validate:"required"`
	Zipcode              string `json:"zipcode"`
	MaritalStatus        int64  `json:"maritalStatus" validate:"required"`
	MotherName           string `json:"motherName" validate:"required"`
	HandPhoneNumber      string `json:"handPhoneNumber"`
	HomePhoneArea        string `json:"homePhoneArea"`
	Email                string `json:"email" validate:"required"`
	KtpImageBase64       string `json:"ktpImageBase64,omitempty" validate:"required"`
	NpwpImageBase64      string `json:"npwpImageBase64,omitempty"`
	SelfieImageBase64    string `json:"selfieImageBase64,omitempty" validate:"required"`
	GoldSavingSlipBase64 string `json:"goldSavingSlipBase64,omitempty"`
	Province             string `json:"province,omitempty" validate:"required"`
	Subdistrict          string `json:"subdistrict,omitempty" validate:"required"`
	Village              string `json:"village,omitempty" validate:"required"`
	Child                int64  `json:"child" validate:"min=0"`
	EmergencyName        string `json:"emergencyName"`
	EmergencyRelation    int64  `json:"emergencyRelation"`
	EmergencyAddress1    string `json:"emergencyAddress1" pg:"emergency_address_1"`
	EmergencyAddress2    string `json:"emergencyAddress2" pg:"emergency_address_2"`
	EmergencyAddress3    string `json:"emergencyAddress3" pg:"emergency_address_3"`
	EmergencyCity        string `json:"emergencyCity"`
	EmergencyZipcode     string `json:"emergencyZipcode"`
	EmergencyPhoneNumber string `json:"emergencyPhoneNumber"`
	RelativePhoneNumber  string `json:"relativePhoneNumber" validate:"required"`
}

// PayloadToken a struct to store all payload for token
type PayloadToken struct {
	UserName string `json:"userName" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RespRegistration to store response registration
type RespRegistration struct {
	ApplicationNumber string `json:"applicationNumber"`
	ApplicationStatus string `json:"applicationStatus,omitempty"`
	CurrentStep       int64  `json:"currentStep"`
}

// PayloadBriRegister a struct to store all payload for a payload bri register
type PayloadBriRegister struct {
	FirstName            string `json:"firstName" validate:"required"`
	LastName             string `json:"lastName" validate:"required"`
	CardName             string `json:"cardName" validate:"required"`
	Npwp                 string `json:"npwp" validate:"required"`
	Nik                  string `json:"nik" validate:"required"`
	BirthPlace           string `json:"birthPlace" validate:"required"`
	BirthDate            string `json:"birthDate" validate:"required"`
	AddressLine1         string `json:"addressLine1" validate:"required" pg:"address_line_1"`
	AddressLine2         string `json:"addressLine2" pg:"address_line_2"`
	AddressLine3         string `json:"addressLine3" pg:"address_line_3"`
	Sex                  int64  `json:"sex" validate:"required" pg:"-"`
	SexString            string `json:"sexString,omitempty" pg:"sex"`
	HomeStatus           int64  `json:"homeStatus" validate:"required"`
	AddressCity          string `json:"addressCity" validate:"required"`
	Nationality          string `json:"nationality" validate:"required"`
	StayedSince          string `json:"stayedSince" validate:"required"`
	Education            int64  `json:"education" validate:"required"`
	Zipcode              string `json:"zipcode" validate:"required"`
	MaritalStatus        int64  `json:"maritalStatus" validate:"required"`
	MotherName           string `json:"motherName" validate:"required"`
	HandPhoneNumber      string `json:"handPhoneNumber" validate:"required"`
	HomePhoneArea        string `json:"homePhoneArea" validate:"required"`
	HomePhoneNumber      string `json:"homePhoneNumber" validate:"required"`
	Email                string `json:"email" validate:"required"`
	Income               int64  `json:"income" validate:"required"`
	JobBidangUsaha       int64  `json:"jobBidangUsaha" validate:"required"`
	JobSubBidangUsaha    int64  `json:"jobSubBidangUsaha" validate:"required"`
	JobCategory          int64  `json:"jobCategory" validate:"required"`
	JobStatus            int64  `json:"jobStatus" validate:"required"`
	TotalEmployee        int64  `json:"totalEmployee" validate:"required"`
	Company              string `json:"company" validate:"required"`
	JobTitle             string `json:"jobTitle" validate:"required"`
	WorkSince            string `json:"workSince" validate:"required"`
	OfficeAddress1       string `json:"officeAddress1" validate:"required" pg:"office_address_1"`
	OfficeAddress2       string `json:"officeAddress2" pg:"office_address_2"`
	OfficeAddress3       string `json:"officeAddress3" pg:"office_address_3"`
	OfficeZipcode        string `json:"officeZipcode" validate:"required"`
	OfficeCity           string `json:"officeCity" validate:"required"`
	OfficePhone          string `json:"officePhone" validate:"required"`
	Child                int64  `json:"child" validate:"min=0"`
	EmergencyName        string `json:"emergencyName" validate:"required"`
	EmergencyRelation    int64  `json:"emergencyRelation" validate:"required"`
	EmergencyAddress1    string `json:"emergencyAddress1" validate:"required" pg:"emergency_address_1"`
	EmergencyAddress2    string `json:"emergencyAddress2" pg:"emergency_address_2"`
	EmergencyAddress3    string `json:"emergencyAddress3" pg:"emergency_address_3"`
	EmergencyCity        string `json:"emergencyCity" validate:"required"`
	EmergencyZipcode     string `json:"emergencyZipcode" validate:"required"`
	EmergencyPhoneNumber string `json:"emergencyPhoneNumber" validate:"required"`
	ProductRequest       string `json:"productRequest" validate:"required"`
	BillingCycle         int64  `json:"billingCycle" validate:"required"`
	CardDeliver          int64  `json:"cardDeliver" validate:"required"`
}

// PayloadOccupation to store response occupation
type PayloadOccupation struct {
	ApplicationNumber string `json:"applicationNumber,omitempty" validate:"required"`
	JobBidangUsaha    int64  `json:"jobBidangUsaha" validate:"required"`
	JobSubBidangUsaha int64  `json:"jobSubBidangUsaha" validate:"required"`
	JobCategory       int64  `json:"jobCategory" validate:"required"`
	JobStatus         int64  `json:"jobStatus" validate:"required"`
	TotalEmployee     int64  `json:"totalEmployee" validate:"required"`
	Company           string `json:"company" validate:"required"`
	JobTitle          string `json:"jobTitle"`
	WorkSince         string `json:"workSince" validate:"required"`
	OfficeAddress1    string `json:"officeAddress1" validate:"required"`
	OfficeAddress2    string `json:"officeAddress2"`
	OfficeAddress3    string `json:"officeAddress3"`
	OfficeZipcode     string `json:"officeZipcode"`
	OfficeCity        string `json:"officeCity"`
	OfficePhone       string `json:"officePhone" validate:"required"`
	Income            int64  `json:"income" validate:"required"`
}

// PayloadActivations a struct to store all payload for activations
type PayloadActivations struct {
	ExpDate           string `json:"expDate" validate:"required"`
	FirstSixDigits    string `json:"firstSixDigits" validate:"required"`
	LastFourDigits    string `json:"lastFourDigits" validate:"required"`
	BirthDate         string `json:"birthDate" validate:"required"`
	ApplicationNumber string `json:"applicationNumber" validate:"required"`
}

// PayloadBRIPendingTransactions a struct to store all payload for transactions pending from BRI
type PayloadBRIPendingTransactions struct {
	TransactionId  string `json:"transactionId" validate:"required"`
	BrixKey        string `json:"brixKey" validate:"required"`
	TrxDateTime    string `json:"trxDateTime" validate:"required"`
	Amount         int64  `json:"amount" validate:"required"`
	CurrencyAmount string `json:"currencyAmount" validate:"required"`
	TrxDesc        string `json:"trxDesc" validate:"required"`
	AuthCode       string `json:"authCode" validate:"required"`
}

// PayloadBRIPaymentTransactions a struct to store all payload for payment transactions from BRI
type PayloadPaymentTransactions struct {
	Source               string `json:"source" validate:"oneof=bri pgdn-core"`
	BillingStatementDate string `json:"billingStatementDate" validate:"required"`
	PaymentAmount        int64  `json:"paymentAmount" validate:"required"`
	RefID                string `json:"refID" validate:"required"`
	BrixKey              string `json:"brixKey" validate:"required"`
	PaymentDate          string `json:"paymentDate" validate:"required"`
	PaymentDesc          string `json:"paymentDesc" validate:"required"`
}

// Payload ResActivations atruct to store all payload for success response activation
type RespActivations struct {
	AccountNumber string `json:"accountNumber"`
}

// Payload PayloadBriGetCardInformation a struct to store all payload for card information inquiry to BRI
type PayloadBriGetCardInformation struct {
	BriXkey string `json:"briXkey" validate:"required"`
}

// PaginationPayload struct to store pagination payload
type PaginationPayload struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

// PayloadHistoryTransactions struct to store request history transactions
type PayloadHistoryTransactions struct {
	AccountNumber string            `json:"accountNumber" validate:"required"`
	Pagination    PaginationPayload `json:"pagination" validate:"required"`
}

// ListHistoryTransactions struct to store list history transactions
type ListHistoryTransactions struct {
	RefTrx      string `json:"refTrx"`
	Nominal     int64  `json:"nominal"`
	TrxDate     string `json:"trxDate"`
	Description string `json:"description"`
}

// ResponseHistoryTransactions struct to store response history transactions
type ResponseHistoryTransactions struct {
	IsLastPage              bool                      `json:"isLastPage"`
	ListHistoryTransactions []ListHistoryTransactions `json:"listHistoryTransactions"`
}

// PayloadAccNumber a struct to store all payload for transactions
type PayloadAccNumber struct {
	AccountNumber string `json:"accountNumber" validate:"required"`
}

// PayloadPaymentInquiry a struct to store all payload for payment inquiry
type PayloadPaymentInquiry struct {
	AccountNumber string `json:"accountNumber" validate:"required"`
	PaymentAmount int64  `json:"paymentAmount" validate:"required"`
}

// PayloadBRIPegadaianBillings a struct to store all payload for post pegadaian billing from BRI
type PayloadBRIPegadaianBillings struct {
	BillingDate   string `json:"billingDate" validate:"required"`
	FileBase64    string `json:"fileBase64" validate:"required,base64"`
	FileExtension string `json:"fileExtension" validate:"required"`
	FileName      string `json:"fileName" validate:"required"`
	RefID         string `json:"refID" validate:"required"`
}
