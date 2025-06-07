package models

type Control struct {
	ID        int64  `json:"id"`
	AlarmID   int64  `json:"alarm_id"`
	UserID    *int64 `json:"user_id,omitempty"`
	Source    string `json:"source,omitempty"`
	Mode      string `json:"mode"`
	Action    string `json:"action"`
	Timestamp string `json:"timestamp"`
	Result    string `json:"result"`
}
