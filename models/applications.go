package models

import (
	"gade/srv-goldcard/logger"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/leekchan/accounting"
)

const (
	// DefBriProductRequest is to store default bri application product request
	DefBriProductRequest = "PAYLATER"

	// DefBriBillingCycle is to store default bri application billing cycle
	DefBriBillingCycle = 3

	// DefBriCardDeliver is to store default bri application card deliver
	DefBriCardDeliver = 1

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

	// AppStatusActive is to store var application status active
	AppStatusActive = "active"

	// AppStatusSent is to store var application status card_sent
	AppStatusSent = "card_sent"

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
		"failed":                "FailedDate",
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
	FailedDate               time.Time  `json:"failedDate,omitempty"`
	Documents                []Document `json:"documents" pg:"-"`
	CoreOpen                 bool       `json:"coreOpen"`
	CreatedAt                time.Time  `json:"createdAt"`
	UpdatedAt                time.Time  `json:"updatedAt"`
}

// SetStatus as a setter for application status
func (app *Applications) SetStatus(msg string) {
	stat := app.getStatus(msg)
	app.Status = stat
	r := reflect.ValueOf(app)
	rNow := reflect.ValueOf(NowDbpg())
	fStatDt := r.Elem().FieldByName(mapStatusDate[stat])
	fStatDt.Set(rNow)
}

// GetStatusDateKey to get status date struct key
func (app *Applications) GetStatusDateKey() string {
	if app.Status == "" {
		logger.Make(nil, nil).Fatal("Application status cannot be nil")
	}

	snake := matchFirstCap.ReplaceAllString(mapStatusDate[app.Status], "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// SetDocument to set application document array
func (app *Applications) SetDocument(pl PayloadPersonalInformation) {
	var emptyDocs []Document
	docNames := []string{"KtpImageBase64", "NpwpImageBase64", "SelfieImageBase64", "GoldSavingSlipBase64", "AppFormBase64"}
	r := reflect.ValueOf(pl)
	currDoc := app.Documents
	app.Documents = emptyDocs

	for _, docName := range docNames {
		base64 := reflect.Indirect(r).FieldByName(docName)

		if base64.IsZero() {
			continue
		}

		doc := app.GetCurrentDoc(currDoc, MapDocType[docName])

		if doc.ID == 0 {
			doc = Document{
				FileName:      pl.Nik + "-" + MapDocType[docName],
				FileExtension: mapFileExt[docName],
				Type:          MapDocType[docName],
				ApplicationID: app.ID,
			}
		}

		doc.FileBase64 = base64.String()
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
	case "application approved", "application on review", "application final approval", "rescheduling delivery":
		return "card_processed"
	case "on printing", "ready to deliver", "on deliver":
		return "card_send"
	case "application rejected":
		return "failed" // TODO: not related naming status
	case "card suspended":
		return "inactive" // TODO: not related naming status
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
	FailedDate               *time.Time `json:"failedDate,omitempty"`
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
}

// MappingApplicationForm a function to mapping application form BRI and slip te data
func (af *ApplicationForm) MappingApplicationForm(params map[string]interface{}) error {
	time := time.Now().UTC()
	acc := params["acc"].(Account)
	docs := params["docs"].([]Document)

	af.Account = acc
	af.TimeStamp = time.Format(DateTimeFormat)
	af.Date = time.Format(DDMMYYYY)
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
	time := time.Now().UTC()
	acc := params["acc"].(Account)
	ac := accounting.Accounting{Symbol: "Rp ", Thousand: "."}

	st.Account = acc
	st.TimeStamp = time.Format(DateTimeFormat)
	st.Date = time.Format(DDMMYYYY)
	st.SignatoryName = params["signatoryName"].(string)
	st.SignatoryNip = params["signatoryNip"].(string)
	st.CardLimitFormat = ac.FormatMoney(acc.Card.CardLimit)
	st.GoldEffBalance = params["goldEffBalance"].(float64)

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
