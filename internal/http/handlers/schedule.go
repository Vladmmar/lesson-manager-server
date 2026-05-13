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

func scheduleHandler(db *storage.Storage, logging *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		full := query.Get("full")
		if full == "true" {
			fullSchedule(w, r, db, logging)
		} else {
			daySchedule(w, r, db, logging)
		}
	}
}

func fullSchedule(w http.ResponseWriter, r *http.Request, db *storage.Storage, logging *slog.Logger) {
	var fullSched [][]lessonHttp.ScheduleLesson
	for dow := 0; dow < 7; dow++ {
		daySched, err := fetchDaySchedule(db, logging, dow)
		if err != nil {
			logging.Error("Failed to fetch day schedule")
		}
		fullSched = append(fullSched, daySched)
	}

	w.Header().Set("Content-Type", "application/json")
	respJson, err := json.Marshal(fullSched)
	if err != nil {
		logging.Error("Failed to convert fullSched object to JSON")
		fmt.Fprint(w, "Failed to convert fullSched object to JSON")
		return
	}
	_, err = fmt.Fprint(w, string(respJson))
	if err != nil {
		logging.Error("Failed to write JSON to html response body")
	}
}

func daySchedule(w http.ResponseWriter, r *http.Request, db *storage.Storage, logging *slog.Logger) {
	query := r.URL.Query()
	dow, err := strconv.Atoi(query.Get("dow"))
	if err != nil && err.Error() != "strconv.Atoi: parsing \"\": invalid syntax" {
		fmt.Fprintf(w, "Incorrect day of week format\n")
		return
	} else if err != nil {
		dow = int(time.Now().Weekday()) - 1
	}

	schedule := lessonHttp.ScheduleResponse{}
	daySched, err := fetchDaySchedule(db, logging, dow)
	if err != nil {
		logging.Error("Failed to fetch schedule")
	}
	schedule.Schedule = daySched

	w.Header().Set("Content-Type", "application/json")
	if len(schedule.Schedule) == 0 {
		schedule.Present = false
	} else {
		schedule.Present = true
	}

	respJson, err := json.Marshal(schedule)
	if err != nil {
		logging.Error("Failed to convert ScheduleResponse object to JSON")
		fmt.Fprint(w, "Failed to convert ScheduleResponse object to JSON")
		return
	}
	_, err = fmt.Fprint(w, string(respJson))
	if err != nil {
		logging.Error("Failed to write JSON to html response body")
	}
}

func fetchDaySchedule(db *storage.Storage, logging *slog.Logger, dow int) ([]lessonHttp.ScheduleLesson, error) {
	rows, err := db.Db.Query(
		`
		SELECT subject, time_start, time_end, dow
		FROM lessons
		WHERE dow = $1
`,
		dow,
	)
	if err != nil {
		logging.Error("internal.http.handlers.current.Init", err.Error())
		return nil, err
	}
	defer rows.Close()

	var schedule []lessonHttp.ScheduleLesson
	n := 0
	for rows.Next() {
		schedule = append(schedule, lessonHttp.ScheduleLesson{})
		err = rows.Scan(
			&schedule[n].Subject,
			&schedule[n].TimeStart,
			&schedule[n].TimeEnd,
			&schedule[n].Dow,
		)
		n++
		if err != nil {
			logging.Error("Failed to parse response from database")
		}
	}

	return schedule, nil
}
