package models

import "time"

// Notification representa uma notificação enviada ao usuário.
type Notification struct {
	UserID    int64     `json:"user_id"`
	AlarmID   int64     `json:"alarm_id"`
	Event     string    `json:"event"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
