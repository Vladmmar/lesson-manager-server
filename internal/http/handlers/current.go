package handlers

import (
	"encoding/json"
	"fmt"
	. "lesson-manager-server/internal/http"
	"lesson-manager-server/internal/storage"
	"log/slog"
	"net/http"
)

type CurrentResponse struct {
	Lessons []LessonResponse
}

func Init(db *storage.Storage, logging *slog.Logger) {
	http.HandleFunc("/current", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Db.Query(`
	SELECT subject, time, meeting_id, password, link
	FROM lessons
	WHERE group_id = $1 and time >= $2+':00' and dow = $3
	LIMIT 2
`, r.PathValue("group_id"), r.PathValue("time"), r.PathValue("dow"))
		if err != nil {
			logging.Error("internal.http.handlers.current.Init", err.Error())
			return
		}
		defer rows.Close()

		l := CurrentResponse{}
		n := 0
		for rows.Next() {
			err = rows.Scan(&l.Lessons[n].Subject,
				&l.Lessons[n].Time,
				&l.Lessons[n].MeetingId,
				&l.Lessons[n].Password,
				&l.Lessons[n].Link,
			)
			if err != nil {
				logging.Error("Failed to parse response from database")
			}
			n++
		}

		short := r.PathValue("short")
		if short == "true" {
			w.Header().Set("Content-Type", "text/json")
			lJson, err := json.Marshal(l)
			if err != nil {
				logging.Error("Failed to convert CurrentResponse object to JSON")
			}
			_, err = fmt.Fprint(w, lJson)
			if err != nil {
				logging.Error("Failed to write JSON to html response body")
			}
		} else {
			w.Header().Set("Content-Type", "text/html")
			WriteLesson(w, &l, 0, false)
			if l.Lessons[1].Link == "" {
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
	fmt.Fprintf(w, "The %s lesson is %s\n", next, l.Lessons[n].Subject)
	fmt.Fprintf(w, "It starts at %s", l.Lessons[n].Time)
	fmt.Fprintf(w, "Meeting id: %s, password: %s", l.Lessons[n].MeetingId, l.Lessons[n].Password)
	fmt.Fprintf(w, "Or you can join it with this link: %s", l.Lessons[n].Link)
}
