package models

type Alarm struct {
	ID       int64  `json:"id"`       
	Location string `json:"location"` 
}

type AlarmUser struct {
	ID      int64 `json:"id"`       
	AlarmID int64 `json:"alarm_id"` 
	UserID  int64 `json:"user_id"`  
}

type AlarmPoint struct {
	ID      int64  `json:"id"`       
	AlarmID int64  `json:"alarm_id"` 
	Point    string `json:"point"`     
}
