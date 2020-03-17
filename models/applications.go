package models

import (
	"gade/srv-goldcard/logger"
	"reflect"
	"regexp"
	"strings"
	"time"
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

	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
	mapStatusDate = map[string]string{
		"application_processed": "ApplicationProcessedDate",
		"card_processed":        "CardProcessedDate",
		"card_send":             "CardSendDate",
		"card_sent":             "CardSentDate",
		"failed":                "FailedDate",
	}
	mapDocType = map[string]string{
		"KtpImageBase64":       "ktp",
		"NpwpImageBase64":      "npwp",
		"SelfieImageBase64":    "selfie",
		"GoldSavingSlipBase64": "slip_te",
	}

	mapFileExt = map[string]string{
		"KtpImageBase64":       DefAppDocFileExt,
		"NpwpImageBase64":      DefAppDocFileExt,
		"SelfieImageBase64":    DefAppDocFileExt,
		"GoldSavingSlipBase64": PdfAppDocFileExt,
	}

	// MapBRIDocType to store map values of BRI DOC type
	MapBRIDocType = map[string]string{
		"ktp":    "A",
		"npwp":   "G",
		"selfie": "P",
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
	CurrentStep              int64      `json:"currentStep"`
	ApplicationProcessedDate time.Time  `json:"applicationProcessedDate,omitempty"`
	CardProcessedDate        time.Time  `json:"cardProcessedDate,omitempty"`
	CardSendDate             time.Time  `json:"cardSendDate,omitempty"`
	CardSentDate             time.Time  `json:"cardSentDate,omitempty"`
	FailedDate               time.Time  `json:"failedDate,omitempty"`
	Documents                []Document `json:"documents" pg:"-"`
	CreatedAt                time.Time  `json:"createdAt"`
	UpdatedAt                time.Time  `json:"updatedAt"`
}

// SetStatus as a setter for application status
func (app *Applications) SetStatus(msg string) {
	stat := app.getStatus(msg)
	app.Status = stat
	r := reflect.ValueOf(app)
	rNow := reflect.ValueOf(time.Now())
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
	docNames := []string{"KtpImageBase64", "NpwpImageBase64", "SelfieImageBase64", "GoldSavingSlipBase64"}
	r := reflect.ValueOf(pl)
	currDoc := app.Documents
	app.Documents = emptyDocs

	for _, docName := range docNames {
		base64 := reflect.Indirect(r).FieldByName(docName)

		if base64.IsZero() {
			continue
		}

		doc := app.getCurrentDoc(currDoc, mapDocType[docName])

		if doc.ID == 0 {
			doc = Document{
				FileName:      pl.Nik + "-" + mapDocType[docName],
				FileExtension: mapFileExt[docName],
				Type:          mapDocType[docName],
				ApplicationID: app.ID,
			}
		}

		doc.FileBase64 = base64.String()
		app.Documents = append(app.Documents, doc)
	}
}

func (app *Applications) getCurrentDoc(currDocs []Document, docType string) Document {
	for _, appDoc := range currDocs {
		if appDoc.Type == docType {
			return appDoc
		}
	}

	return Document{}
}

func (app *Applications) getStatus(msg string) string {
	switch strings.ToLower(msg) {
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
