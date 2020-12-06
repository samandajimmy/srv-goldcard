package models

import (
	"gade/srv-goldcard/logger"
	"strings"
	"time"
)

var (
	JobCategoryStr = map[int64]string{
		1: "Karyawan",
		2: "Profesional",
		3: "Pensiunan",
		4: "TNI/POLRI",
		5: "Wiraswasta",
		6: "Lain-lain",
	}

	JobBidangUsahaStr = map[int64]string{
		10: "Agricultural & Animal Rising",
		20: "Aneka Industry",
		30: "Customer Product",
		40: "Financial",
		50: "Goverment",
		60: "Industry and Chemical",
		70: "Infrastructure",
		80: "Mining",
		90: "Trading and Service",
		99: "Lain-lain",
	}
)

// DefJobTitle to store default value of job title
const DefJobTitle = "MANAGER"

// Occupation is a struct to store occupation data
type Occupation struct {
	ID                int64     `json:"id"`
	JobBidangUsaha    int64     `json:"jobBidangUsaha"`
	JobSubBidangUsaha int64     `json:"jobSubBidangUsaha"`
	JobCategory       int64     `json:"jobCategory"`
	JobStatus         int64     `json:"jobStatus"`
	TotalEmployee     int64     `json:"totalEmployee"`
	Company           string    `json:"company"`
	JobTitle          string    `json:"jobTitle"`
	WorkSince         string    `json:"workSince"`
	OfficeAddress1    string    `json:"officeAddress1" pg:"office_address_1"`
	OfficeAddress2    string    `json:"officeAddress2" pg:"office_address_2"`
	OfficeAddress3    string    `json:"officeAddress3" pg:"office_address_3"`
	OfficeZipcode     string    `json:"officeZipcode"`
	OfficeProvince    string    `json:"officeProvince"`
	OfficeCity        string    `json:"officeCity"`
	OfficeSubdistrict string    `json:"officeSubdistrict"`
	OfficeVillage     string    `json:"officeVillage"`
	OfficePhone       string    `json:"officePhone"`
	Income            int64     `json:"income"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// MappingOccupation a function to map all data occupation
func (occ *Occupation) MappingOccupation(pl PayloadOccupation, addrData AddressData) error {
	occ.JobBidangUsaha = pl.JobBidangUsaha
	occ.JobSubBidangUsaha = pl.JobSubBidangUsaha
	occ.JobCategory = pl.JobCategory
	occ.JobStatus = pl.JobStatus
	occ.TotalEmployee = pl.TotalEmployee
	occ.Company = pl.Company
	occ.JobTitle = pl.JobTitle
	occ.WorkSince = pl.WorkSince
	occ.OfficeZipcode = addrData.Zipcode
	occ.OfficeProvince = addrData.Province
	occ.OfficeCity = addrData.City
	occ.OfficeSubdistrict = addrData.Subdistrict
	occ.OfficeVillage = addrData.Village
	occ.OfficePhone = pl.OfficePhone
	occ.Income = pl.Income * 12

	addrData.AddressLine1 = pl.Company + " " + pl.OfficeAddress1 +
		" Kel " + strings.Title(strings.ToLower(pl.OfficeVillage)) +
		" Kec " + strings.Title(strings.ToLower(pl.OfficeSubdistrict))
	// set addressData
	addrData, err := RemappAddress(addrData, 30)

	if err != nil {
		logger.Make(nil, nil).Debug(err)

		return err
	}

	occ.OfficeAddress1 = addrData.AddressLine1
	occ.OfficeAddress2 = addrData.AddressLine2
	occ.OfficeAddress3 = addrData.AddressLine3
	return nil
}
