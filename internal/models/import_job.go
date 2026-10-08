package models

import "time"

type ImportJob struct {
	ID           string     `gorm:"primaryKey" json:"id"`
	FileName     string     `json:"file_name"`
	Status       string     `json:"status"`
	TotalRows    int        `json:"total_rows"`
	SuccessRows  int        `json:"success_rows"`
	FailedRows   int        `json:"failed_rows"`
	ErrorDetails string     `json:"error_details"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}
