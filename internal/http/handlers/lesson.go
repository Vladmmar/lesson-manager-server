package handlers

import (
	"encoding/json"
	"fmt"
	lessonHttp "lesson-manager-server/internal/http"
	"lesson-manager-server/internal/storage"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func currentLessonHandler(db *storage.Storage, logging *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		hours, minutes, _ := time.Now().Clock()
		rows, err := db.Db.Query(
			`
	SELECT subject, time_start, time_end, dow, meeting_id, password, link
	FROM lessons
	WHERE group_id = $1 and time_start <= $2 and time_end >= $3 and dow = $4
	LIMIT 1
`,
			query.Get("group_id"),
			strconv.Itoa(hours)+":"+strconv.Itoa(minutes)+":00",
			strconv.Itoa(hours)+":"+strconv.Itoa(minutes)+":00",
			int(time.Now().Weekday())-1,
		)
		if err != nil {
			logging.Error("internal.http.handlers.current.Init", err.Error())
			return
		}
		defer rows.Close()

		lesson := lessonHttp.Lesson{}
		for rows.Next() {
			err = rows.Scan(
				&lesson.Subject,
				&lesson.TimeStart,
				&lesson.TimeEnd,
				&lesson.Dow,
				&lesson.MeetingId,
				&lesson.Password,
				&lesson.Link,
			)
			if err != nil {
				logging.Error("Failed to parse response from database")
			}
		}

		var resp lessonHttp.LessonResponse
		w.Header().Set("Content-Type", "application/json")
		if lesson.Link == "" {
			resp = lessonHttp.LessonResponse{
				Present: false,
				Lesson:  lesson,
			}
		} else {
			resp = lessonHttp.LessonResponse{
				Present: true,
				Lesson:  lesson,
			}
		}

		jsonResponse := query.Get("json")
		if jsonResponse == "true" {
			lJson, err := json.Marshal(resp)
			if err != nil {
				logging.Error("Failed to convert CurrentResponse object to JSON")
			}
			_, err = fmt.Fprint(w, string(lJson))
			if err != nil {
				logging.Error("Failed to write JSON to html response body")
			}
		} else {
			w.Header().Set("Content-Type", "text/html")
			WriteLesson(w, &resp.Lesson)
		}
	}
}

func nextLessonHandler(db *storage.Storage, logging *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		hours, minutes, _ := time.Now().Clock()

		lessonOn, err := strconv.Atoi(query.Get("lesson_on"))
		if err == nil {
			nextLessonOn(w, r, db, logging, lessonOn)
			return
		} else if err.Error() != "strconv.Atoi: parsing \"\": invalid syntax" {
			fmt.Fprintf(w, "Incorrect lesson on format\n")
			return
		}

		rows, err := db.Db.Query(`
	SELECT subject, time_start, time_end, dow, meeting_id, password, link
	FROM lessons
	WHERE group_id = $1 and time_start >= $2 and dow = $3
	LIMIT 1
`, query.Get("group_id"), strconv.Itoa(hours)+":"+strconv.Itoa(minutes)+":00", int(time.Now().Weekday())-1)
		if err != nil {
			logging.Error("internal.http.handlers.current.Init", err.Error())
			return
		}
		defer rows.Close()

		lesson := lessonHttp.Lesson{}
		for rows.Next() {
			err = rows.Scan(
				&lesson.Subject,
				&lesson.TimeStart,
				&lesson.TimeEnd,
				&lesson.Dow,
				&lesson.MeetingId,
				&lesson.Password,
				&lesson.Link,
			)
			if err != nil {
				logging.Error("Failed to parse response from database")
			}
		}

		var resp lessonHttp.LessonResponse
		if lesson.Link == "" {
			resp = lessonHttp.LessonResponse{
				Present: false,
				Lesson:  lesson,
			}
		} else {
			resp = lessonHttp.LessonResponse{
				Present: true,
				Lesson:  lesson,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		lJson, err := json.Marshal(resp)
		if err != nil {
			logging.Error("Failed to convert CurrentResponse object to JSON")
		}
		_, err = fmt.Fprint(w, string(lJson))
		if err != nil {
			logging.Error("Failed to write JSON to html response body")
		}
	}
}

func nextLessonOn(w http.ResponseWriter, r *http.Request, db *storage.Storage, logging *slog.Logger, lessonOn int) {
	query := r.URL.Query()
	dow, err := strconv.Atoi(query.Get("dow"))
	if err != nil && err.Error() != "strconv.Atoi: parsing \"\": invalid syntax" {
		fmt.Fprintf(w, "Incorrect day of week format\n")
		return
	} else if err != nil {
		dow = int(time.Now().Weekday()) - 1
	}

	rows, err := db.Db.Query(`
			SELECT subject, time_start, time_end, dow, meeting_id, password, link
			FROM lessons
			WHERE group_id = $1 and dow = $2
			ORDER BY time_start 
			LIMIT $3
			OFFSET $4
		`, query.Get("group_id"), dow, 1, lessonOn-1)
	if err != nil {
		logging.Error("internal.http.handlers.current.Init", err.Error())
		return
	}
	defer rows.Close()

	lesson := lessonHttp.Lesson{}
	for rows.Next() {
		err = rows.Scan(
			&lesson.Subject,
			&lesson.TimeStart,
			&lesson.TimeEnd,
			&lesson.Dow,
			&lesson.MeetingId,
			&lesson.Password,
			&lesson.Link,
		)
		if err != nil {
			logging.Error("Failed to parse response from database")
		}
	}

	var resp lessonHttp.LessonResponse
	if lesson.Link != "" {
		resp = lessonHttp.LessonResponse{
			Present: true,
			Lesson:  lesson,
		}
	} else {
		resp = lessonHttp.LessonResponse{
			Present: false,
			Lesson:  lesson,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	lJson, err := json.Marshal(resp)
	if err != nil {
		logging.Error("Failed to convert CurrentResponse object to JSON")
	}
	_, err = fmt.Fprint(w, string(lJson))
	if err != nil {
		logging.Error("Failed to write JSON to html response body")
	}
}

func WriteLesson(w http.ResponseWriter, l *lessonHttp.Lesson) {
	fmt.Fprintf(w, "The current lesson is %s<br>", l.Subject)
	fmt.Fprintf(w, "It starts at %s<br>", l.TimeStart)
	fmt.Fprintf(w, "Meeting id: %s, password: %s<br>", l.MeetingId, l.Password)
	fmt.Fprintf(w, "Or you can join it via this link: <a href=%s>%s<a><hr>", l.Link, l.Link)
}
