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
