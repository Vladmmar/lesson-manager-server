package handlers

import (
	"lesson-manager-server/internal/storage"
	"log/slog"
	"net/http"
)

func Init(db *storage.Storage, logging *slog.Logger) {
	rows, err := db.Db.Query(`
	SELECT subject, time, meeting_id, password, link
	FROM lessons
	WHERE group_id = $1 and time >= '$2:00:00' and dow = $3
	LIMIT 2
`)
	if err != nil {
		logging.Error("internal.http.handlers.current.Init", err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {

	}

	http.HandleFunc("/current", func(w http.ResponseWriter, r *http.Request) {

		cookies := r.PathValue("short")
		if cookies == "true" {

		}
	})
}
