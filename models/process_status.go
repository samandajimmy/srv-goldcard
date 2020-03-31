package models

import (
	"time"
)

var (
	// Final application process var
	AppProcType = "Application"
	// Core Open
	FinalRegCoreOpenErr = "CoreOpenGC-Error"
	FinalRegCoreOpenSuc = "CoreOpenGC-Success"
	// BRI Registration
	FinalRegBRIRegisErr = "BRIRegGC-Error"
	FinalRegBRIRegisSuc = "BRIRegGC-Success"
	// BRI UploadDocument
	FinalRegBRIUploadDocErr = "BRIUploadDoc-Error"
	FinalRegBRIUploadDocSuc = "BRIUploadDoc-Success"
)

type ProcessStatus struct {
	ID          int64     `json:"id"`
	ProcessID   string    `json:"processId"`
	ProcessType string    `json:"processType"`
	Status      []string  `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (ps *ProcessStatus) MapInsertProcessStatus(process_id, process_type string, status string) error {
	ps.ProcessID = process_id
	ps.ProcessType = process_type
	ps.Status = append(ps.Status, status)
	ps.CreatedAt = time.Now()

	return nil
}

func (ps *ProcessStatus) MapUpdateProcessStatus(status string) error {
	if !Contains(ps.Status, status) {
		ps.Status = append(ps.Status, status)
	}
	ps.UpdatedAt = time.Now()

	return nil
}
