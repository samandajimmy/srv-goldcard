package models

import (
	"math"
	"time"
)

var (
	// StatusSuccess to store a status response success
	StatusSuccess = "Success"

	// MessageDataSuccess to store a success message response of data
	MessageDataSuccess = "Data Berhasil Dikirim"

	// MessageUnprocessableEntity to store a message response of unproccessable entity
	MessageUnprocessableEntity = "Entitas Tidak Dapat Diproses"

	// MicroTimeFormat to store a time format of micro timestamp
	MicroTimeFormat = "20060102150405.000000"

	// DateTimeFormat to store a date time format of timestamp
	DateTimeFormat = "2006-01-02 15:04:05"

	// DateTimeFormatZone to store a date time with zone format of timestamp
	DateTimeFormatZone = "2006-01-02T15:04:05.000Z07:00"

	// DateTimeFormatMillisecond to store a date time format of timestamp to millisecond
	DateTimeFormatMillisecond = "2006-01-02 15:04:05.000"

	// DateFormat to store a date format of timestamp
	DateFormat = "2006-01-02"

	// DateFormatRegex to store a regex of dd/mm/yyyy date format
	DateFormatRegex = "(^\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])$)"

	// LetterBytes a string to generate random ID
	LetterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// BriBankCode a string to store BRI bank code
	BriBankCode = "002"

	// EmergencyContactDef a string to store default emergency contact
	EmergencyContactDef = "pegadaian"

	// StarString a string with stars
	StarString = "**********"
)

// Parameter struct is represent a data for parameters model
type Parameter struct {
	ID          int64     `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CustomRound is a function to round the number based on the type
func CustomRound(roundType string, num float64, decimal float64) float64 {
	switch roundType {
	case "ceil":
		return math.Ceil(num*decimal) / decimal
	case "floor":
		return math.Floor(num*decimal) / decimal
	default:
		return math.Round(num*decimal) / decimal
	}
}

// DateIsNotEqual is a function to compare date now and update date
func DateIsNotEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 != y2 || m1 != m2 || d1 != d2
}

// ParseDate is a function parse date from DD-MM-YYYY to YYYY-MM-DD
func ParseDate(data string) string {
	date, _ := time.Parse("02-01-2006", data)

	return date.Format(DateFormat)
}
