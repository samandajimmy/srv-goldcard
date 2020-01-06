package models

import "time"

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
	OfficeAddress1    string    `json:"officeAddress1"`
	OfficeAddress2    string    `json:"officeAddress2"`
	OfficeAddress3    string    `json:"officeAddress3"`
	OfficeZipcode     string    `json:"officeZipcode"`
	OfficeCity        string    `json:"officeCity"`
	OfficePhone       string    `json:"officePhone"`
	Income            int64     `json:"income"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
