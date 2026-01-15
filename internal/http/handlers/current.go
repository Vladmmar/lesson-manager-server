package handlers

import (
	"encoding/json"
	"fmt"
	. "lesson-manager-server/internal/http"
	"lesson-manager-server/internal/storage"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type CurrentResponse struct {
	Lessons []LessonResponse `json:"lessons"`
}

func Init(db *storage.Storage, logging *slog.Logger) {
	http.HandleFunc("/current", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		hours, minutes, _ := time.Now().Clock()
		rows, err := db.Db.Query(`
	SELECT subject, time, meeting_id, password, link
	FROM lessons
	WHERE group_id = $1 and time >= $2 and dow = $3
	LIMIT 2
`, query.Get("group_id"), strconv.Itoa(hours)+":"+strconv.Itoa(minutes)+":00", int(time.Now().Weekday())-1)
		if err != nil {
			logging.Error("internal.http.handlers.current.Init", err.Error())
			return
		}
		defer rows.Close()

		l := CurrentResponse{}
		l.Lessons = []LessonResponse{}
		for rows.Next() {
			var lesson LessonResponse
			err = rows.Scan(&lesson.Subject,
				&lesson.Time,
				&lesson.MeetingId,
				&lesson.Password,
				&lesson.Link,
			)
			if err != nil {
				logging.Error("Failed to parse response from database")
			}
			l.Lessons = append(l.Lessons, lesson)
		}

		if len(l.Lessons) == 0 {
			fmt.Fprintln(w, "No lessons were found with group_id "+query.Get("group_id"))
			return
		}

		short := query.Get("short")
		if short == "true" {
			w.Header().Set("Content-Type", "application/json")
			lJson, err := json.Marshal(l)
			if err != nil {
				logging.Error("Failed to convert CurrentResponse object to JSON")
			}
			_, err = fmt.Fprint(w, string(lJson))
			if err != nil {
				logging.Error("Failed to write JSON to html response body")
			}
		} else {
			w.Header().Set("Content-Type", "text/html")
			WriteLesson(w, &l, 0, false)
			if len(l.Lessons) > 1 {
				WriteLesson(w, &l, 1, true)
			}
		}
	})
}

func WriteLesson(w http.ResponseWriter, l *CurrentResponse, n int, isNext bool) {
	next := "current"
	if isNext {
		next = "next"
	}
	fmt.Fprintf(w, "The %s lesson is %s<br>", next, l.Lessons[n].Subject)
	fmt.Fprintf(w, "It starts at %s<br>", l.Lessons[n].Time)
	fmt.Fprintf(w, "Meeting id: %s, password: %s<br>", l.Lessons[n].MeetingId, l.Lessons[n].Password)
	fmt.Fprintf(w, "Or you can join it via this link: <a href=%s>%s<a><hr>", l.Lessons[n].Link, l.Lessons[n].Link)
}
