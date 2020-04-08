package models

import (
	"time"
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
	OfficeCity        string    `json:"officeCity"`
	OfficePhone       string    `json:"officePhone"`
	Income            int64     `json:"income"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// MappingOccupation a function to map all data occupation
func (occ *Occupation) MappingOccupation(pl PayloadOccupation) error {
	occ.JobBidangUsaha = pl.JobBidangUsaha
	occ.JobSubBidangUsaha = pl.JobSubBidangUsaha
	occ.JobCategory = pl.JobCategory
	occ.JobStatus = pl.JobStatus
	occ.TotalEmployee = pl.TotalEmployee
	occ.Company = pl.Company
	occ.JobTitle = pl.JobTitle
	occ.WorkSince = pl.WorkSince
	occ.OfficeAddress1 = pl.OfficeAddress1
	occ.OfficeAddress2 = pl.OfficeAddress2
	occ.OfficeAddress3 = pl.OfficeAddress3
	occ.OfficeZipcode = pl.OfficeZipcode
	occ.OfficeCity = pl.OfficeCity
	occ.OfficePhone = pl.OfficePhone
	occ.Income = pl.Income * 12

	return nil
}
