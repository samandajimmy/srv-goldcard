package models

import (
	"encoding/base64"
	"gade/srv-goldcard/logger"
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"

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

	// DateFormat to store a format of date yyyy-mm-dd
	DateFormat = "2006-01-02"

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

	// DMYSLASH to strore a date format for DMYSLASH
	DMYSLASH = "02/01/2006"

	// DateFormatRegex to store a regex of dd/mm/yyyy date format
	DateFormatRegex = "(^\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])$)"

	// LetterBytes a string to generate random ID
	LetterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// BriBankCode a string to store BRI bank code
	BriBankCode = "002"

	// EmergencyContactDef a string to store default emergency contact
	EmergencyContactDef = "pegadaian"

	// Default image Base64
	DefDocBase64 = "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"
)

type FuncAfterGC func(c echo.Context, acc *Account, briPl PayloadBriRegister, accChan chan Account) error

// Parameter struct is represent a data for parameters model
type Parameter struct {
	ID          int64     `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NowDbpg this function only because db pg when insert or update conver time format with UTC
// so when you INSERT/UPDATE using DBPG then you need this for to get time now
func NowDbpg() time.Time {
	return time.Now().Add(7 * time.Hour)
}

// NowUTC to get real current datetime but UTC format
func NowUTC() time.Time {
	return time.Now().UTC().Add(7 * time.Hour)
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

func GenerateApplicationFormPDF(af ApplicationForm, templatePath string) (string, error) {
	generatePdf := RequestPdf{}
	requestPdf := generatePdf.NewRequestPdf("")

	// Mapping Application Form Data to Html
	err := requestPdf.ParseTemplate(templatePath, af)
	if err != nil {
		logger.Make(nil, nil).Debug(err)
		return "", err
	}

	// Generate Pdf File
	bufPdf, err := requestPdf.GeneratePDF()
	if err != nil {
		logger.Make(nil, nil).Debug(err)
		return "", err
	}

	// Convert to base64
	pdfBase64 := base64.StdEncoding.EncodeToString(bufPdf)
	return pdfBase64, nil
}

func GenerateSlipTePDF(st SlipTE, templatePath string) (string, error) {
	generatePdf := RequestPdf{}
	requestPdf := generatePdf.NewRequestPdf("")

	// Mapping Application Form Data to Html
	err := requestPdf.ParseTemplate(templatePath, st)
	if err != nil {
		logger.Make(nil, nil).Debug(err)
		return "", err
	}

	// Generate Pdf File
	bufPdf, err := requestPdf.GeneratePDF()
	if err != nil {
		logger.Make(nil, nil).Debug(err)
		return "", err
	}

	// Convert to base64
	pdfBase64 := base64.StdEncoding.EncodeToString(bufPdf)
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
func StringNameFormatter(name string, length int, withSpace bool) string {
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
	for i := arrLen - 1; i >= 0; i-- {
		// break if length of result is enough
		if len(result) < length {
			break
		}

		// take first char of word and replace itself
		if withSpace {
			result = strings.ReplaceAll(result, arrStr[i], strings.ToUpper(arrStr[i][:1]))
			continue
		}

		result = strings.ReplaceAll(result, " "+arrStr[i], "")
	}

	// the last attempt if its still not enough hard cut it out
	if len(result) > length {
		result = result[:length]
	}

	return strings.TrimSpace(result)
}

// GetInterfaceValue to get interface value dynamicly
func GetInterfaceValue(r reflect.Value, key string) interface{} {
	val := r.FieldByName(key)

	if !val.IsValid() {
		switch key {
		case "PaymentAmount", "Amount":
			return int64(0)
		default:
			return ""
		}
	}

	if val.IsZero() {
		switch r.FieldByName(key).Kind() {
		case reflect.Int64:
			return int64(0)
		default:
			return ""
		}
	}

	switch r.FieldByName(key).Kind() {
	case reflect.Int64:
		return val.Interface().(int64)
	default:
		return val.Interface().(string)
	}
}

func ReverseArray(data []ListTrx) []ListTrx {
	dataLength := len(data) - 1
	reverseDataLength := dataLength / 2

	for i := 0; i < reverseDataLength; i++ {
		temp := data[dataLength]
		data[dataLength] = data[i]
		data[i] = temp
		dataLength--
	}

	return data
}

func RemappAddress(addr AddressData, length int) (AddressData, error) {
	keyAddr := []string{"", "", ""}

	// Make a Regex to say we only want letters and numbers
	specialChar := regexp.MustCompile("[^a-zA-Z0-9]+")
	space := regexp.MustCompile(`\s+`)
	// remove special characters
	processedString := specialChar.ReplaceAllString(addr.AddressLine1, " ")
	// remove double spaces
	processedString = space.ReplaceAllString(processedString, " ")
	arrStr := strings.Split(processedString, " ")
	mapIdx := 0
	newStr := arrStr[0]

	for idx, str := range arrStr {
		if idx == 0 {
			continue
		}

		if len(str) >= (length - 1) {
			return AddressData{}, ErrAddrNotGood
		}

		if mapIdx > 2 {
			break
		}

		newStr += " " + str

		if len(newStr) > length {
			mapIdx += 1
			newStr = str
			continue
		}

		keyAddr[mapIdx] = newStr
	}

	addr.AddressLine1 = keyAddr[0]
	addr.AddressLine2 = keyAddr[1]
	addr.AddressLine3 = keyAddr[2]

	return addr, nil
}
