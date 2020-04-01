package models

import (
	"fmt"
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
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (ps *ProcessStatus) MapInsertProcessStatus(processType, tableName string, processID int64, reason interface{}) error {
	rString := fmt.Sprintf("%v", reason)

	ps.Reason = rString
	ps.ProcessID = processID
	ps.ProcessType = processType
	ps.TblName = tableName
	ps.CreatedAt = time.Now()

	return nil
}
