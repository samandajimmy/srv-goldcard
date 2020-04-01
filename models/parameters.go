package models

import (
	"bytes"
	"encoding/base64"
	"math"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/labstack/echo"
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

	// DateFormatDef to store a date format of timestamp
	DateFormatDef = "2006-01-02"

	// DDMMYYYY to strore a date format for DDMMYYYY
	DDMMYYYY = "02-01-2006"

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

	// Default image Base64
	defDocBase64 = "/9j/4AAQSkZJRgABAQABLAEsAAD/4QCMRXhpZgAATU0AKgAAAAgABQESAAMAAAABAAEAAAEaAAUAAAABAAAASgEbAAUAAAABAAAAUgEoAAMAAAABAAIAAIdpAAQAAAABAAAAWgAAAAAAAAEsAAAAAQAAASwAAAABAAOgAQADAAAAAQABAACgAgAEAAAAAQAAAO2gAwAEAAAAAQAAAJYAAAAA/+0AOFBob3Rvc2hvcCAzLjAAOEJJTQQEAAAAAAAAOEJJTQQlAAAAAAAQ1B2M2Y8AsgTpgAmY7PhCfv/CABEIAJYA7QMBEQACEQEDEQH/xAAfAAABBQEBAQEBAQAAAAAAAAADAgQBBQAGBwgJCgv/xADDEAABAwMCBAMEBgQHBgQIBnMBAgADEQQSIQUxEyIQBkFRMhRhcSMHgSCRQhWhUjOxJGIwFsFy0UOSNIII4VNAJWMXNfCTc6JQRLKD8SZUNmSUdMJg0oSjGHDiJ0U3ZbNVdaSVw4Xy00Z2gONHVma0CQoZGigpKjg5OkhJSldYWVpnaGlqd3h5eoaHiImKkJaXmJmaoKWmp6ipqrC1tre4ubrAxMXGx8jJytDU1dbX2Nna4OTl5ufo6erz9PX29/j5+v/EAB8BAAMBAQEBAQEBAQEAAAAAAAECAAMEBQYHCAkKC//EAMMRAAICAQMDAwIDBQIFAgQEhwEAAhEDEBIhBCAxQRMFMCIyURRABjMjYUIVcVI0gVAkkaFDsRYHYjVT8NElYMFE4XLxF4JjNnAmRVSSJ6LSCAkKGBkaKCkqNzg5OkZHSElKVVZXWFlaZGVmZ2hpanN0dXZ3eHl6gIOEhYaHiImKkJOUlZaXmJmaoKOkpaanqKmqsLKztLW2t7i5usDCw8TFxsfIycrQ09TV1tfY2drg4uPk5ebn6Onq8vP09fb3+Pn6/9sAQwABAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAgICAgICAgICAgIDAwMDAwMDAwMD/9sAQwEBAQEBAQEBAQEBAgIBAgIDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMD/9oADAMBAAIRAxEAAAH9/K1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrVq1atWrV//2gAIAQEAAQUC/wDLHX//2gAIAQMRAT8B/wC9HX//2gAIAQIRAT8B/wC9HX//2gAIAQEABj8C/wDLHX//xAAzEAEAAwACAgICAgMBAQAAAgsBEQAhMUFRYXGBkaGxwfDREOHxIDBAUGBwgJCgsMDQ4P/aAAgBAQABPyH/APkdf//aAAwDAQACEQMRAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA//EADMRAQEBAAMAAQIFBQEBAAEBCQEAESExEEFRYSBx8JGBobHRweHxMEBQYHCAkKCwwNDg/9oACAEDEQE/EP8A+R1//9oACAECEQE/EP8A+R1//9oACAEBAAE/EP8A+R1//9k="
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

func GenerateGoldSavingPDF(pl PayloadPersonalInformation) (string, error) {
	var buf bytes.Buffer
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// CellFormat(width, height, text, border, position after, align, fill, link, linkStr)
	pdf.CellFormat(190, 7, "TEST", "0", 0, "CM", false, 0, "")

	err := pdf.Output(&buf)

	if err != nil {
		return "", err
	}

	// Convert to base64
	pdfBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return pdfBase64, nil
}

func GenerateAppFormPDF(pl PayloadPersonalInformation) (string, error) {
	var buf bytes.Buffer
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// CellFormat(width, height, text, border, position after, align, fill, link, linkStr)
	pdf.CellFormat(190, 7, "TEST", "0", 0, "CM", false, 0, "")

	err := pdf.Output(&buf)

	if err != nil {
		return "", err
	}

	// Convert to base64
	pdfBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return pdfBase64, nil
}

// StringCutter to force or hard cut any sentences
func StringCutter(sentence string, length int) string {
	if len(sentence) <= length {
		return sentence
	}

	return sentence[:length]
}

// StringNameFormatter is a function to limit the size of string according to given length
// by shorten the given string
func StringNameFormatter(name string, length int) string {
	var result string
	arrStr := strings.Split(name, " ")
	arrLen := len(arrStr)

	// if length of name is enough then just return
	if len(name) <= length {
		return name
	}

	// if its only one word
	if arrLen == 1 {
		result = arrStr[0]
	}

	// if length of words is enough then just return
	// note that var result could be filled / not
	if len(result) >= length {
		return result[:length]
	}

	// init string result
	result = name

	for _, word := range arrStr {
		// break if length of result is enough
		// why using "<", because we will add a space after the loop
		if len(result) < length {
			break
		}

		// take first char of word and replace itself
		result = strings.ReplaceAll(result, word+" ", strings.ToUpper(word[:1])+".")
	}

	for i := arrLen - 1; i > 0; i-- {
		// put the prev word back if its still too long
		// why using ">=", because we will add a space after the loop
		if len(result) >= length {
			result = strings.ReplaceAll(result, arrStr[i-1][:1]+"."+arrStr[i], arrStr[i-1])
		}
	}

	// get last dot "." index
	lastDotAfterIdx := strings.LastIndex(result, ".") + 1
	// add space
	result = result[:lastDotAfterIdx] + " " + result[lastDotAfterIdx:]

	// the last attempt if its still not enough hard cut it out
	if len(result) > length {
		result = result[:length]
	}

	return result
}

func Contains(arr []string, x string) bool {
	for _, n := range arr {
		if x == n {
			return true
		}
	}
	return false
}

type FuncAfterGC func(c echo.Context, acc *Account, briPl PayloadBriRegister, accChan chan Account, errAppBri, errAppCore chan error) error
