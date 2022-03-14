package model

import (
	"time"
)

var (
	// Application process var
	FinalAppProcessType = "Final Application"

	// Table Names
	ApplicationTableName = "applications"
)

type ProcessStatus struct {
	ID          int64     `json:"id"`
	ProcessID   int64     `json:"processId"`
	ProcessType string    `json:"processType"`
	TblName     string    `json:"tblName"`
	Reason      string    `json:"reason"`
	ErrorCount  int64     `json:"errorCount"`
	Error       string    `json:"error"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (ps *ProcessStatus) MapUpdateProcessStatus(tableName string, processID int64) error {
	ps.ProcessID = processID
	ps.TblName = tableName

	return nil
}
