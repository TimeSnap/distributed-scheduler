package model

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type JobExecution struct {
	ID                 int         `json:"id"`
	JobID              uuid.UUID   `json:"job_id"`
	StartTime          time.Time   `json:"start_time"`
	EndTime            time.Time   `json:"end_time"`
	Success            bool        `json:"success"`
	NumberOfExecutions int         `json:"number_of_executions"`
	NumberOfRetries    int         `json:"number_of_retries"`
	ErrorMessage       null.String `json:"error_message,omitempty" swaggertype:"string"`
}

type JobExecutionStatus string

const (
	JobExecutionStatusSuccessful JobExecutionStatus = "SUCCESSFUL"
	JobExecutionStatusFailed     JobExecutionStatus = "FAILED"
)
