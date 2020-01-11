package models

import (
	"gade/srv-goldcard/logger"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var (
	// DefAppDocFileExt is to store var default application document file ext
	DefAppDocFileExt = "jpg"

	// DefAppDocType is to store var default application document type
	DefAppDocType = "D"

	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
	mapStatusDate = map[string]string{
		"application_processed": "ApplicationProcessedDate",
		"card_processed":        "CardProcessedDate",
		"card_send":             "CardSendDate",
		"card_sent":             "CardSentDate",
		"failed":                "FailedDate",
	}
)

// Applications is a struct to store application data
type Applications struct {
	ID                       int64     `json:"id"`
	ApplicationNumber        string    `json:"applicationNumber" validate:"required"`
	Status                   string    `json:"status"`
	KtpImageBase64           string    `json:"ktpImageBase64"`
	NpwpImageBase64          string    `json:"npwpImageBase64"`
	SelfieImageBase64        string    `json:"selfieImageBase64"`
	KtpDocID                 string    `json:"ktpDocId"`
	NpwpDocID                string    `json:"npwpDocId"`
	SelfieDocID              string    `json:"selfieDocId"`
	SavingAccount            string    `json:"savingAccount" validate:"required"`
	ApplicationProcessedDate time.Time `json:"applicationProcessedDate,omitempty"`
	CardProcessedDate        time.Time `json:"cardProcessedDate,omitempty"`
	CardSendDate             time.Time `json:"cardSendDate,omitempty"`
	CardSentDate             time.Time `json:"cardSentDate,omitempty"`
	FailedDate               time.Time `json:"failedDate,omitempty"`
	CreatedAt                time.Time `json:"createdAt"`
	UpdatedAt                time.Time `json:"updatedAt"`
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

func (app *Applications) getStatus(msg string) string {
	switch msg {
	default:
		return "application_processed"
	}
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
