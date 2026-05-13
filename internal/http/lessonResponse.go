package http

import "time"

type Lesson struct {
	Subject   string    `json:"subject"`
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Dow       int       `json:"dow"`
	MeetingId string    `json:"meeting_id"`
	Password  string    `json:"password"`
	Link      string    `json:"link"`
}

type LessonResponse struct {
	Present bool   `json:"present"`
	Lesson  Lesson `json:"lesson"`
}
