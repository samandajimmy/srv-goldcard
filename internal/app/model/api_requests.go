package model

import (
	"encoding/json"
	"time"
)

// APIRequest is a struct to store apiRequest data
type APIRequest struct {
	ID           int64           `json:"id"`
	RequestID    string          `json:"requestId"`
	HostName     string          `json:"hostName"`
	Endpoint     string          `json:"endpoint"`
	Status       string          `json:"status"`
	RequestData  json.RawMessage `json:"requestData"`
	ResponseData json.RawMessage `json:"responseData"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}
