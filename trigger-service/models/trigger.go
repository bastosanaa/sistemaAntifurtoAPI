package models

import "time"

// Trigger representa um disparo de alarme registrado no sistema.
type Trigger struct {
	ID        int64     `json:"id"`
	AlarmID   int64     `json:"alarm_id"`
	Point     string    `json:"point"`
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
}
