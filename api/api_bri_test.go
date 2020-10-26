package api_test

import (
	"fmt"
	"gade/srv-goldcard/api"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
)

var registerRequest = map[string]interface{}{
	"firstName":            "Jonathans",
	"lastName":             "BRS",
	"cardName":             "Jon Wick",
	"npwp":                 "12312312312",
	"nik":                  "3173080102640077",
	"birthPlace":           "DKI JAKARTA",
	"birthDate":            "1990-02-01",
	"addressLine1":         "dsf no. 100",
	"addressLine2":         "",
	"addressLine3":         "RT110 RW220, GROGOL PETAMBURAN",
	"sex":                  1,
	"homeStatus":           1,
	"addressCity":          "KOTA JAKARTA BARAT",
	"nationality":          "WNI",
	"stayedSince":          "2008-06-23",
	"education":            2,
	"zipcode":              "11470",
	"maritalStatus":        2,
	"motherName":           "Tina",
	"handPhoneNumber":      "0817123456",
	"homePhoneArea":        "021",
	"homePhoneNumber":      "58078901",
	"email":                "ferdian@cermati.com",
	"jobBidangUsaha":       10,
	"jobSubBidangUsaha":    10,
	"jobCategory":          3,
	"jobStatus":            2,
	"totalEmployee":        1,
	"company":              "Cermati",
	"jobTitle":             "APM",
	"workSince":            "2007-01-23",
	"officeAddress1":       "Jalan Tomang Raya no. 38",
	"officeAddress2":       "",
	"officeAddress3":       "KECAMATAN KOPO",
	"officeZipcode":        "42178",
	"officeCity":           "KABUPATEN SERANG",
	"officePhone":          "021121212",
	"income":               "168000000",
	"child":                7,
	"emergencyName":        "sdf",
	"emergencyRelation":    1,
	"emergencyAddress1":    "asdf no. 121",
	"emergencyAddress2":    "",
	"emergencyAddress3":    "RT 110 RW 110, KECAMATAN GROGOL PETAMBURAN",
	"emergencyCity":        "KOTA JAKARTA BARAT",
	"emergencyPhoneNumber": "0811233244",
	"productRequest":       "MEASY",
	"billingCycle":         3,
}

func TestNewBriAPI(t *testing.T) {
	response := api.BriResponse{}
	body := map[string]interface{}{
		"requestData": registerRequest,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	bri, _ := api.NewBriAPI(c)
	reqBody, _ := bri.Request("/register", echo.POST, body)
	resp, _ := bri.Do(reqBody, &response)

	fmt.Println(reqBody)
	fmt.Println(resp)
	fmt.Println(response)
}
