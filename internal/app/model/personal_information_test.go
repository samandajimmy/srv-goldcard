package model_test

import (
	"srv-goldcard/internal/app/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

// personal information dummy data
var pi = model.PersonalInformation{
	ID:                  1,
	FirstName:           "AIRLANGGA",
	LastName:            "RAMADHON",
	HandPhoneNumber:     "",
	Email:               "m.nurfaizal96@gmail.com",
	Npwp:                "",
	Nik:                 "3173082903900004",
	BirthPlace:          "TANGERANG",
	BirthDate:           "1990-03-29",
	Nationality:         "WNI",
	Sex:                 "male",
	Education:           3,
	MaritalStatus:       1,
	MotherName:          "SRI REJEKI",
	HomePhoneArea:       "",
	HomePhoneNumber:     "",
	HomeStatus:          1,
	AddressLine1:        "Jalan Laksamana 92 A",
	AddressLine2:        "",
	AddressLine3:        "",
	Zipcode:             "",
	AddressCity:         "MALANG",
	StayedSince:         "",
	Child:               0,
	RelativePhoneNumber: "082245492497",
}

func TestSetHomePhone(t *testing.T) {
	// when homephoneNumber is not ""
	pi.HomePhoneNumber = "22113344"
	assert.Equal(t, "22113344", pi.HomePhoneNumber)

	// when RelativePhoneNumber is ""
	pi.RelativePhoneNumber = ""
	assert.Equal(t, "22113344", pi.HomePhoneNumber)

	// when homephoneNumber is ""
	// and RelativePhoneNumber is not ""
	pi.RelativePhoneNumber = "082245492497"
	pi.HomePhoneNumber = ""
	pi.SetHomePhone()
	assert.Equal(t, "0822", pi.HomePhoneArea)
	assert.Equal(t, "45492497", pi.HomePhoneNumber)
}

func TestSetNPWP(t *testing.T) {
	// when npwp is ""
	pi.SetNPWP("")
	assert.Equal(t, "00000000000", pi.Npwp)

	// when npwp is "12312312312"
	pi.SetNPWP("12312312312")
	assert.Equal(t, "12312312312", pi.Npwp)
}

func TestGetSexInt(t *testing.T) {
	// when sex is male
	result := pi.GetSexInt("male")
	assert.Equal(t, int64(1), result)

	// when sex is female
	result = pi.GetSexInt("female")
	assert.Equal(t, int64(2), result)

	// when sex is other then male, female
	result = pi.GetSexInt("others")
	assert.Equal(t, int64(0), result)
}

func TestGetSex(t *testing.T) {
	// when sex int 1
	result := pi.GetSex(int64(1))
	assert.Equal(t, "male", result)

	// when sex int 2
	result = pi.GetSex(int64(2))
	assert.Equal(t, "female", result)

	// when sex is other then 1, 2
	result = pi.GetSex(int64(3))
	assert.Equal(t, "", result)
}
