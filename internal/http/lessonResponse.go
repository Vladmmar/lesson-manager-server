package http

import "time"

type LessonResponse struct {
	Subject   string    `json:"subject"`
	Time      time.Time `json:"time"`
	MeetingId string    `json:"meeting_id"`
	Password  string    `json:"password"`
	Link      string    `json:"link"`
}
