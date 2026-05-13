package http

import "time"

type ScheduleLesson struct {
	Subject   string    `json:"subject"`
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Dow       int       `json:"dow"`
}

type ScheduleResponse struct {
	Present  bool             `json:"present"`
	Schedule []ScheduleLesson `json:"schedule"`
}
