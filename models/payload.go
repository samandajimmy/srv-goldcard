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
	CardDeliver       int64  `json:"cardDeliver" validate:"required"`
}

// PayloadPersonalInformation a struct to store all payload for a payload personal information
type PayloadPersonalInformation struct {
	ApplicationNumber    string `json:"applicationNumber,omitempty" validate:"required"`
	FirstName            string `json:"firstName" validate:"required"`
	LastName             string `json:"lastName"`
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
	AppFormBase64        string `json:"appFormBase64,omitempty"`
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
	LastName             string `json:"lastName"`
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
	Income               int64  `json:"income" validate:"required,max=9999999999999"`
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
	Limit                int64  `json:"limit" validate:"required" pg:"card_limit"`
}

// PayloadOccupation to store response occupation
type PayloadOccupation struct {
	ApplicationNumber string `json:"applicationNumber,omitempty" validate:"required"`
	JobBidangUsaha    int64  `json:"jobBidangUsaha" validate:"required"`
	JobSubBidangUsaha int64  `json:"jobSubBidangUsaha" validate:"required"`
	JobCategory       int64  `json:"jobCategory" validate:"required"`
	JobStatus         int64  `json:"jobStatus" validate:"required"`
	TotalEmployee     int64  `json:"totalEmployee" validate:"required"`
	Company           string `json:"company" validate:"required,max=30"`
	JobTitle          string `json:"jobTitle"`
	WorkSince         string `json:"workSince" validate:"required"`
	OfficeAddress1    string `json:"officeAddress1" validate:"required"`
	OfficeAddress2    string `json:"officeAddress2"`
	OfficeAddress3    string `json:"officeAddress3"`
	OfficeCity        string `json:"officeCity" validate:"required"`
	OfficeProvince    string `json:"officeProvince" validate:"required"`
	OfficeSubdistrict string `json:"officeSubdistrict" validate:"required"`
	OfficeVillage     string `json:"officeVillage" validate:"required"`
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

// PayloadCardBlock a struct to store all payload for block a card
type PayloadCardBlock struct {
	AccountNumber string `json:"accountNumber" validate:"required"`
	Reason        string `json:"reason" validate:"required"`
	ReasonCode    string `json:"reasonCode" validate:"required"`
	BlockedDate   string `json:"blockedDate"`
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

// PlPaymentInquiry a struct to store all payload for payment inquiry
type PlPaymentInquiry struct {
	AccountNumber string `json:"accountNumber" validate:"required"`
	PaymentAmount int64  `json:"paymentAmount" validate:"required"`
	RefTrx        string `json:"refTrx"`
}

// RespPaymentInquiry a struct to store all response for payment inquiry
type RespPaymentInquiry struct {
	AccountNumber string `json:"accountNumber"`
	PaymentAmount int64  `json:"paymentAmount"`
	RefTrx        string `json:"refTrx"`
}

// RespLimitUpdateInquiry a struct to store all response for update limit inquiry
type RespUpdateLimitInquiry struct {
	RefId string `json:"refId"`
}

// PlPaymentTrxCore a struct to store all payload for payment transactions for core
type PlPaymentTrxCore struct {
	Source        string `json:"source"`
	AccountNumber string `json:"accountNumber" validate:"required"`
	RefTrx        string `json:"refTrx" validate:"required"`
	PaymentAmount int64  `json:"paymentAmount"`
}

// PayloadPaymentTransactions a struct to store all payload for payment transactions
type PayloadPaymentTransactions struct {
	Source               string `json:"source" validate:"oneof=bri"`
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
	Page  int64 `json:"page" validate:"required"`
}

// PayloadListTrx struct to store request history transactions
type PayloadListTrx struct {
	AccountNumber string `json:"accountNumber" validate:"required"`
}

// PayloadAccNumber a struct to store all payload for transactions
type PayloadAccNumber struct {
	AccountNumber string `json:"accountNumber" validate:"required"`
}

// PayloadBRIPegadaianBillings a struct to store all payload for post pegadaian billing from BRI
type PayloadBRIPegadaianBillings struct {
	BillingDate   string `json:"billingDate" validate:"required"`
	FileBase64    string `json:"fileBase64" validate:"required,base64"`
	FileExtension string `json:"fileExtension" validate:"required"`
	FileName      string `json:"fileName" validate:"required"`
	RefID         string `json:"refID" validate:"required"`
}

// PayloadCoreDecreasedSTL a struct store all payload for new gold price when stl decreased bigger than 5 %
type PayloadCoreDecreasedSTL struct {
	STL                  int64  `json:"stl" validate:"required"`
	DecreasedFivePercent string `json:"decreasedFivePercent" validate:"required"`
}

// PayloadInquiryUpdateLimit a struct to store all payload for inquiry update limit
type PayloadInquiryUpdateLimit struct {
	AccountNumber string `json:"accountNumber" validate:"required"`
	NominalLimit  int64  `json:"nominalLimit" validate:"required"`
}

// PayloadUpdateLimit a struct to store all payload for submit update limit
type PayloadUpdateLimit struct {
	RefId           string `json:"refId" validate:"required"`
	NpwpImageBase64 string `json:"npwpImageBase64"`
}

// PayloadCoreGtePayment is a struct to store all payload for get
type PayloadCoreGtePayment struct {
	SavingAccount      string `json:"savingAccount" validate:"required"`
	AvailableGram      string `json:"availableGram" validate:"required"`
	NominalTransaction int64  `json:"nominalTransaction" validate:"required"`
	TrxId              string `json:"trxId" validate:"required"`
}

// PayloadInsertPublicHoliday is a struct to store all payload for insert public holiday date
type PayloadInsertPublicHoliday struct {
	PublicHolidayDate []string `json:"publicHolidayDate" validate:"required,gte=1"`
}

// PayloadBRICardReplace is a struct to store payload calling api card replace to BRI
type PayloadBRICardReplace struct {
	BriXkey string `json:"briXkey" validate:"required"`
}

// ValidateBRIRegisterSpecification a function to validate registration specification to BRI
func (plBRIReg *PayloadBriRegister) ValidateBRIRegisterSpecification() error {
	plBRIReg.FirstName = StringNameFormatter(plBRIReg.FirstName, 15, false)
	plBRIReg.LastName = StringNameFormatter(plBRIReg.LastName, 14, false)
	plBRIReg.CardName = StringNameFormatter(plBRIReg.CardName, 19, true)
	plBRIReg.Nik = StringCutter(plBRIReg.Nik, 30)
	plBRIReg.Npwp = StringCutter(plBRIReg.Npwp, 15)
	plBRIReg.BirthPlace = StringCutter(plBRIReg.BirthPlace, 20)
	plBRIReg.BirthDate = StringCutter(plBRIReg.BirthDate, 30)
	plBRIReg.AddressLine1 = StringCutter(plBRIReg.AddressLine1, 30)
	plBRIReg.AddressLine2 = StringCutter(plBRIReg.AddressLine2, 30)
	plBRIReg.AddressLine3 = StringCutter(plBRIReg.AddressLine3, 30)
	plBRIReg.AddressCity = StringCutter(plBRIReg.AddressCity, 28)
	plBRIReg.Nationality = StringCutter(plBRIReg.Nationality, 3)
	plBRIReg.MotherName = StringNameFormatter(plBRIReg.MotherName, 30, true)
	plBRIReg.HandPhoneNumber = StringCutter(plBRIReg.HandPhoneNumber, 13)
	plBRIReg.HomePhoneArea = StringCutter(plBRIReg.HomePhoneArea, 5)
	plBRIReg.HomePhoneNumber = StringCutter(plBRIReg.HomePhoneNumber, 10)
	plBRIReg.Email = StringCutter(plBRIReg.Email, 50)
	plBRIReg.Company = StringCutter(plBRIReg.Company, 25)
	plBRIReg.JobTitle = StringCutter(plBRIReg.JobTitle, 30)
	plBRIReg.OfficeAddress1 = StringCutter(plBRIReg.OfficeAddress1, 30)
	plBRIReg.OfficeAddress2 = StringCutter(plBRIReg.OfficeAddress2, 30)
	plBRIReg.OfficeAddress3 = StringCutter(plBRIReg.OfficeAddress3, 30)
	plBRIReg.OfficeCity = StringCutter(plBRIReg.OfficeCity, 30)
	plBRIReg.OfficePhone = StringCutter(plBRIReg.OfficePhone, 13)
	plBRIReg.EmergencyName = StringNameFormatter(plBRIReg.EmergencyName, 30, true)
	plBRIReg.EmergencyAddress1 = StringCutter(plBRIReg.EmergencyAddress1, 100)
	plBRIReg.EmergencyAddress2 = StringCutter(plBRIReg.EmergencyAddress2, 100)
	plBRIReg.EmergencyAddress3 = StringCutter(plBRIReg.EmergencyAddress3, 100)
	plBRIReg.EmergencyCity = StringCutter(plBRIReg.EmergencyCity, 50)
	plBRIReg.EmergencyPhoneNumber = StringCutter(plBRIReg.EmergencyPhoneNumber, 13)
	plBRIReg.ProductRequest = StringCutter(plBRIReg.ProductRequest, 30)

	return nil
}

// PayloadCoreGtePayment is a struct to store all payload for get
type RespGetAddress struct {
	CardDeliver int64       `json:"cardDeliver"`
	Office      AddressData `json:"office"`
	Domicile    AddressData `json:"domicile"`
}

// PayloadCoreGtePayment is a struct to store all payload for get
type RespCardStatus struct {
	Status     string `json:"status"`
	IsReplaced string `json:"isReplaced"`
}
