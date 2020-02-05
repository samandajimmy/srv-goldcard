package models

import "time"

type Transactions struct {
	ID              int64     `json:"id"`
	AccountId       string    `json:"accountId"`
	RefTrxPegadaian string    `json:"refTrxPegadaian"`
	RefTrx          int64     `json:"refTrx"`
	Nominal         float64   `json:"nominal"`
	GoldNominal     int64     `json:"goldNominal"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	Balance         string    `json:"balance"`
	GoldBalance     string    `json:"goldBalance"`
	Methods         string    `json:"methods"`
	UpdatedAt       time.Time `json:"updatedAt"`
	CreatedAt       time.Time `json:"createdAt"`
	Account         Account   `json:"account"`
}
