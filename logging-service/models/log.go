package models

import "time"

// LogEntry representa um evento de armacao, desarmacao ou disparo registrado.
type LogEntry struct {
	ID        int64     `json:"id"`
	Service   string    `json:"service"`
	AlarmID   int64     `json:"alarm_id"`
	UserID    *int64    `json:"user_id,omitempty"`
	Action    string    `json:"action"`
	Mode      *string   `json:"mode,omitempty"`
	Point     *string   `json:"point,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
