package http

import "time"

type LessonResponse struct {
	Subject   string
	Time      time.Time
	MeetingId string
	Password  string
	Link      string
}
