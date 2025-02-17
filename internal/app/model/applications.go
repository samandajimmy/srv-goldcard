package model

import (
	"reflect"
	"regexp"
	"srv-goldcard/internal/pkg/logger"
	"strings"
	"time"

	"github.com/leekchan/accounting"
)

const (
	// DefBriProductRequest is to store default bri application product request
	DefBriProductRequest = "KARTUEMAS"

	// DefBriBillingCycle is to store default bri application billing cycle
	DefBriBillingCycle = 2

	// BriCardDeliverHome is to store default bri application card deliver to home
	BriCardDeliverHome = 1

	// BriCardDeliverOffice is to store default bri application card deliver to office
	BriCardDeliverOffice = 2

	// DefAppDocFileExt is to store var default application document file ext
	DefAppDocFileExt = "jpg"

	// PdfAppDocFileExt is to store var pdf application document file ext
	PdfAppDocFileExt = "pdf"

	// DefAppDocType is to store var default application document type
	DefAppDocType = "D"

	// SlipTeTemplatePath is to store path file template Slip TE
	SlipTeTemplatePath = "template/template_slip.html"

	// ApplicationFormTemplatePath is to store path file Application Form BRI
	ApplicationFormTemplatePath = "template/template_application_form.html"
)

var (
	// AppStatusOngoing is to store var application status ongoing
	AppStatusOngoing = "application_ongoing"

	// AppStatusProcessed is to store var application status proccesed
	AppStatusProcessed = "application_processed"

	// AppStatusCardProcessed is to store var application status card processed
	AppStatusCardProcessed = "card_processed"

	// AppStatusActive is to store var application status active
	AppStatusActive = "active"

	// AppStatusInactive is to store var application status inactive
	AppStatusInactive = "inactive"

	// AppStatusSent is to store var application status card_sent
	AppStatusSent = "card_sent"

	// AppStatusRejected is to store var application status rejected
	AppStatusRejected = "application rejected"

	// AppStatusExpired is to store var application status expired
	AppStatusExpired = "expired"

	// AppStatusForceDeliver is to store var application status force delivery
	AppStatusForceDeliver = "force_deliver"

	// AppStepSavingAcc is to store var application step saving account
	AppStepSavingAcc int64 = 1

	// AppStepCardLimit is to store var application step card limit
	AppStepCardLimit int64 = 2

	// AppStepPersonalInfo is to store var application step personal info
	AppStepPersonalInfo int64 = 3

	// AppStepOccupation is to store var application step post occupation
	AppStepOccupation int64 = 4

	// AppStepAddress is to store var application step post address
	AppStepAddress int64 = 5

	// AppStepCompleted is to store var application step completed
	AppStepCompleted int64 = 99

	// TextFileFound is var to store if file is found in database
	TextFileFound = "Ada"

	// TextFileNotFound is var to store if file is not found in database
	TextFileNotFound = "Tidak ada"

	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
	mapStatusDate = map[string]string{
		"application_processed": "ApplicationProcessedDate",
		"card_processed":        "CardProcessedDate",
		"card_send":             "CardSendDate",
		"card_sent":             "CardSentDate",
		"rejected":              "RejectedDate",
		"card_suspended":        "RejectedDate",
	}
	MapDocType = map[string]string{
		"KtpImageBase64":       "ktp",
		"NpwpImageBase64":      "npwp",
		"SelfieImageBase64":    "selfie",
		"GoldSavingSlipBase64": "slip_te",
		"AppFormBase64":        "app_form",
	}

	mapFileExt = map[string]string{
		"KtpImageBase64":       DefAppDocFileExt,
		"NpwpImageBase64":      DefAppDocFileExt,
		"SelfieImageBase64":    DefAppDocFileExt,
		"GoldSavingSlipBase64": PdfAppDocFileExt,
		"AppFormBase64":        PdfAppDocFileExt,
	}

	// MapBRIDocType to store map values of BRI DOC type
	MapBRIDocType = map[string]string{
		"ktp":      "A",
		"npwp":     "G",
		"selfie":   "P",
		"slip_te":  "D",
		"app_form": "Z",
	}

	// MapBRIExtBase64File to store BRI ext Base64 file
	MapBRIExtBase64File = map[string]string{
		"jpg": "data:image/jpeg;base64,",
		"pdf": "data:application/pdf;base64,",
	}

	DocTypes = []string{
		"ktp",
		"npwp",
		"selfie",
		"slip_te",
		"app_form",
		"undefined",
	}
)

// Applications is a struct to store application data
type Applications struct {
	// nolint
	tableName struct{} `pg:"applications"`

	ID                       int64      `json:"id"`
	ApplicationNumber        string     `json:"applicationNumber" validate:"required"`
	Status                   string     `json:"status"`
	SavingAccount            string     `json:"savingAccount" validate:"required"`
	SavingAccountOpeningDate string     `json:"savingAccountOpeningDate" validate:"required"`
	CurrentStep              int64      `json:"currentStep"`
	ApplicationProcessedDate time.Time  `json:"applicationProcessedDate,omitempty"`
	CardProcessedDate        time.Time  `json:"cardProcessedDate,omitempty"`
	CardSendDate             time.Time  `json:"cardSendDate,omitempty"`
	CardSentDate             time.Time  `json:"cardSentDate,omitempty"`
	RejectedDate             time.Time  `json:"rejectedDate,omitempty"`
	ExpiredAt                time.Time  `json:"expiredAt"`
	Documents                []Document `json:"documents" pg:"-"`
	CoreOpen                 bool       `json:"coreOpen"`
	CardLimit                int64      `json:"cardLimit"`
	CreatedAt                time.Time  `json:"createdAt"`
	UpdatedAt                time.Time  `json:"updatedAt"`
}

// SetStatus as a setter for application status
func (app *Applications) SetStatus(msg string) {
	stat := app.getStatus(msg)
	mapStat := mapStatusDate[stat]

	if stat == app.Status || mapStat == "" {
		return
	}

	app.Status = stat
	r := reflect.ValueOf(app)
	rNow := reflect.ValueOf(NowDbpg())
	fStatDt := r.Elem().FieldByName(mapStat)
	fStatDt.Set(rNow)
}

// GetStatusDateKey to get status date struct key
func (app *Applications) GetStatusDateKey() string {
	if app.Status == "" {
		return ""
	}

	if mapStatusDate[app.Status] == "" {
		return ""
	}

	snake := matchFirstCap.ReplaceAllString(mapStatusDate[app.Status], "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// SetDocument to set application document array
func (app *Applications) SetDocument(pl PayloadPersonalInformation) {
	var emptyDocs []Document
	var doc Document
	var base64 reflect.Value
	var base64Str string
	docNames := []string{"KtpImageBase64", "NpwpImageBase64", "SelfieImageBase64", "GoldSavingSlipBase64", "AppFormBase64"}
	r := reflect.ValueOf(pl)
	currDoc := app.Documents
	app.Documents = emptyDocs

	for _, docName := range docNames {
		base64 = reflect.Indirect(r).FieldByName(docName)
		base64Str = base64.String()
		doc = app.GetCurrentDoc(currDoc, MapDocType[docName])

		if base64.IsZero() {
			base64Str = doc.FileBase64
		}

		if doc.ID == 0 {
			doc = Document{
				FileName:      pl.Nik + "-" + MapDocType[docName],
				FileExtension: mapFileExt[docName],
				Type:          MapDocType[docName],
				ApplicationID: app.ID,
			}
		}

		doc.FileBase64 = base64Str
		app.Documents = append(app.Documents, doc)
	}
}

func (app *Applications) GetCurrentDoc(currDocs []Document, docType string) Document {
	for _, appDoc := range currDocs {
		if appDoc.Type == docType {
			return appDoc
		}
	}

	return Document{}
}

func (app *Applications) getStatus(msg string) string {
	switch strings.ToLower(msg) {
	case "application on review", "application final approval", "rescheduling delivery":
		return "application_processed"
	case "application approved", "on printing", "ready to deliver":
		return "card_processed"
	case "on deliver":
		return "card_send"
	case "application rejected":
		return "rejected"
	case "card suspended":
		return "card_suspended"
	case "delivered":
		return "card_sent"
	default:
		return "application_processed"
	}
}

// Document is a struct to store document data
type Document struct {
	ID            int64     `json:"id"`
	FileName      string    `json:"fileName"`
	FileBase64    string    `json:"fileBase64"`
	FileExtension string    `json:"fileExtension"`
	Type          string    `json:"type"`
	DocID         string    `json:"docId"`
	ApplicationID int64     `json:"applicationId"`
	UpdatedAt     time.Time `json:"updatedAt"`
	CreatedAt     time.Time `json:"createdAt"`
}

// AppDocument is a struct to store application document data
type AppDocument struct {
	BriXkey    string `json:"briXkey"`
	DocType    string `json:"docType"`
	FileName   string `json:"fileName"`
	FileExt    string `json:"fileExt"`
	Base64file string `json:"base64file"`
}

// AppStatus is a struct to store application status data
type AppStatus struct {
	Status                   string     `json:"status"`
	ApplicationProcessedDate *time.Time `json:"applicationProcessedDate,omitempty"`
	CardProcessedDate        *time.Time `json:"cardProcessedDate,omitempty"`
	CardSendDate             *time.Time `json:"cardSendDate,omitempty"`
	CardSentDate             *time.Time `json:"cardSentDate,omitempty"`
	RejectedDate             *time.Time `json:"rejectedDate,omitempty"`
}

// ApplicationForm a struct to store all payload for Application Form BRI
type ApplicationForm struct {
	Account            Account `json:"account"`
	Date               string  `json:"date"`
	TimeStamp          string  `json:"timeStamp"`
	TextHomeStatus     string  `json:"textHomeStatus"`
	TextEducation      string  `json:"textEducation"`
	TextMaritalStatus  string  `json:"textMaritalStatus"`
	TextJobBidangUsaha string  `json:"textJobBidangUsaha"`
	TextJobCategory    string  `json:"textJobCategory"`
	TextRelation       string  `json:"textrelation"`
	FileKtp            string  `json:"fileKtp"`
	FileSelfie         string  `json:"fileSelfie"`
	FileNpwp           string  `json:"fileNpwp"`
	FileAppForm        string  `json:"fileAppForm"`
	FileSlipTe         string  `json:"fileSlipTe"`
	ShippingAddress1   string  `json:"shippingAddress1"`
	ShippingAddress2   string  `json:"shippingAddress2"`
	ShippingAddress3   string  `json:"shippingAddress3"`
}

// SlipTE a struct to store all payload for Slip TE Document
type SlipTE struct {
	Account         Account `json:"account"`
	Date            string  `json:"date"`
	TimeStamp       string  `json:"timeStamp"`
	CardLimitFormat string  `json:"cardLimitFormat"`
	SignatoryName   string  `json:"signatoryName"`
	SignatoryNip    string  `json:"signatoryNip"`
	GoldEffBalance  float64 `json:"goldEffBalance"`
	OpeningDate     string  `json:"openingDate"`
}

// MappingApplicationForm a function to mapping application form BRI and slip te data
func (af *ApplicationForm) MappingApplicationForm() error {
	now := time.Now()
	acc := af.Account
	docs := af.Account.Application.Documents

	af.TimeStamp = now.Format(DateTimeFormat)
	af.Date = now.Format(DDMMYYYY)
	af.TextHomeStatus = HomeStatusStr[acc.PersonalInformation.HomeStatus]
	af.TextEducation = EducationStr[acc.PersonalInformation.Education]
	af.TextMaritalStatus = MaritalStatusStr[acc.PersonalInformation.MaritalStatus]
	af.TextJobBidangUsaha = JobBidangUsahaStr[acc.Occupation.JobBidangUsaha]
	af.TextJobCategory = JobCategoryStr[acc.Occupation.JobCategory]
	af.TextRelation = RelationStr[acc.EmergencyContact.Relation]
	af.FileKtp = TextFileNotFound
	af.FileNpwp = TextFileNotFound
	af.FileSelfie = TextFileNotFound
	af.FileAppForm = TextFileFound
	af.FileSlipTe = TextFileFound
	af.ShippingAddress1 = acc.PersonalInformation.AddressLine1
	af.ShippingAddress2 = acc.PersonalInformation.AddressLine2
	af.ShippingAddress3 = acc.PersonalInformation.AddressLine3

	if acc.CardDeliver == BriCardDeliverOffice {
		af.ShippingAddress1 = acc.Occupation.OfficeAddress1
		af.ShippingAddress2 = acc.Occupation.OfficeAddress2
		af.ShippingAddress3 = acc.Occupation.OfficeAddress3
	}

	for _, document := range docs {
		switch document.Type {
		case "ktp":
			af.FileKtp = TextFileFound
		case "npwp":
			af.FileNpwp = TextFileFound
		case "selfie":
			af.FileSelfie = TextFileFound
		}
	}

	// Set App Form Base64
	appFormBase64, err := GenerateApplicationFormPDF(*af, ApplicationFormTemplatePath)

	if err != nil {
		logger.Make(nil, nil).Debug(err)
		return err
	}

	// Set Application Form and Slip TE
	personalInformation := PayloadPersonalInformation{}
	personalInformation.Nik = acc.PersonalInformation.Nik
	personalInformation.AppFormBase64 = appFormBase64

	// Set Application Document
	af.Account.Application.SetDocument(personalInformation)

	return nil
}

// MappingSlipTe a function to mapping slip te data
func (st *SlipTE) MappingSlipTe(params map[string]interface{}) error {
	now := time.Now()
	acc := st.Account
	ac := accounting.Accounting{Symbol: "Rp ", Thousand: "."}
	opDate, _ := time.Parse(DateTimeFormatZone, acc.Application.SavingAccountOpeningDate)

	st.Account = acc
	st.TimeStamp = now.Format(DateTimeFormat)
	st.Date = now.Format(DMYSLASH)
	st.SignatoryName = params["signatoryName"].(string)
	st.SignatoryNip = params["signatoryNip"].(string)
	st.CardLimitFormat = ac.FormatMoney(acc.Card.CardLimit)
	st.GoldEffBalance = params["goldEffBalance"].(float64)
	st.OpeningDate = opDate.Format(DMYSLASH)

	// Set gold saving slip Base64
	slipBase64, err := GenerateSlipTePDF(*st, SlipTeTemplatePath)

	if err != nil {
		return err
	}

	// Set Application Form and Slip TE
	personalInformation := PayloadPersonalInformation{}
	personalInformation.Nik = acc.PersonalInformation.Nik
	personalInformation.GoldSavingSlipBase64 = slipBase64

	// Set Application Document
	st.Account.Application.SetDocument(personalInformation)

	return nil
}

// SavingAccount a struct to store saving account number
type SavingAccount struct {
	SavingAccount string `json:"savingAccount"`
}
