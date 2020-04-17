package models

import (
	"encoding/json"
	"time"
)

// PaymentInquiry is a struct to store paymentInquiry data
type PaymentInquiry struct {
	ID           string          `json:"id"`
	AccountId    int64           `json:"accountId"`
	BillingId    int64           `json:"billingId"`
	RefTrx       string          `json:"refTrx"`
	Nominal      int64           `json:"nominal"`
	Status       string          `json:"status"`
	CoreResponse json.RawMessage `json:"coreResponse"`
	InquiryDate  time.Time       `json:"inquiryDate"`
	UpdatedAt    time.Time       `json:"updatedAt"`
	CreatedAt    time.Time       `json:"createdAt"`
	Billing      Billing         `json:"billing"`
}

// PaymentInquiryNotificationData is a struct to store PaymentInquiryNotificationData data
type PaymentInquiryNotificationData struct {
	ID             string `json:"id"`
	Administration string `json:"administration"`
	ReffSwitching  string `json:"reffSwitching" pg:"reff_switching"`
}
