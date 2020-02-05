package models

import "time"

var (
	defWNI         string = "1"
	defWNA         string = "2"
	defStayedSince string = "00/00"
	defNPWP        string = "00000000000"
	defChildNumber int64  = 0
	briSexInt             = map[int64]string{
		1: "male",
		2: "female",
	}
	nationalityMap = map[string]string{
		defWNI: "WNI",
		defWNA: "WNA",
	}
)

// PersonalInformation is a struct to store personal info data
type PersonalInformation struct {
	ID                  int64     `json:"id"`
	FirstName           string    `json:"firstName"`
	LastName            string    `json:"lastName"`
	HandPhoneNumber     string    `json:"handPhoneNumber"`
	Email               string    `json:"email"`
	Npwp                string    `json:"npwp"`
	Nik                 string    `json:"nik"`
	BirthPlace          string    `json:"birthPlace"`
	BirthDate           string    `json:"birthDate"`
	Nationality         string    `json:"nationality"`
	Sex                 string    `json:"sex"`
	Education           int64     `json:"education"`
	MaritalStatus       int64     `json:"maritalStatus"`
	MotherName          string    `json:"motherName"`
	HomePhoneArea       string    `json:"homePhoneArea"`
	HomePhoneNumber     string    `json:"homePhoneNumber"`
	HomeStatus          int64     `json:"homeStatus"`
	AddressLine1        string    `json:"addressLine1" pg:"address_line_1"`
	AddressLine2        string    `json:"addressLine2" pg:"address_line_2"`
	AddressLine3        string    `json:"addressLine3" pg:"address_line_3"`
	Zipcode             string    `json:"zipcode"`
	AddressCity         string    `json:"addressCity"`
	StayedSince         string    `json:"stayedSince"`
	Child               int64     `json:"child"`
	RelativePhoneNumber string    `json:"relativePhoneNumber"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

// GetSex to get goldcard sex status
func (pi *PersonalInformation) GetSex(sex int64) string {
	return briSexInt[sex]
}

// GetSexInt to get bri sex status
func (pi *PersonalInformation) GetSexInt(sex string) int64 {
	var res int64

	for k, v := range briSexInt {
		if v == sex {
			return k
		}
	}

	return res
}

// SetHomePhone to set home phone number value
func (pi *PersonalInformation) SetHomePhone() {
	if pi.HomePhoneNumber != "" {
		return
	}

	if pi.RelativePhoneNumber == "" {
		return
	}

	// string into slices character
	runes := []rune(pi.RelativePhoneNumber)
	pi.HomePhoneArea = string(runes[0:4])
	pi.HomePhoneNumber = string(runes[4:])
}

// SetNPWP to get npwp value
func (pi *PersonalInformation) SetNPWP(npwp string) {
	if npwp != "" {
		pi.Npwp = npwp
		return
	}

	pi.Npwp = defNPWP
}
