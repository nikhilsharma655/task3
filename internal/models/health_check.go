package models

import "time"

type HealthCheck struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CheckType string    `json:"check_type"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Details   string    `json:"details"`
	CheckedAt time.Time `json:"checked_at"`
}
