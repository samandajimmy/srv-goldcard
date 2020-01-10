package models

import "time"

var (
	// DefAppDocFileExt is to store var default application document file ext
	DefAppDocFileExt = "jpg"

	// DefAppDocType is to store var default application document type
	DefAppDocType = "D"
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
	Status                   string    `json:"status"`
	ApplicationProcessedDate time.Time `json:"applicationProcessedDate,omitempty"`
	CardProcessedDate        time.Time `json:"cardProcessedDate,omitempty"`
	CardSendDate             time.Time `json:"cardSendDate,omitempty"`
	CardSentDate             time.Time `json:"cardSentDate,omitempty"`
	FailedDate               time.Time `json:"failedDate,omitempty"`
}
