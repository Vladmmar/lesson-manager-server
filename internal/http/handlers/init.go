package handlers

import (
	"lesson-manager-server/internal/storage"
	"log/slog"
	"net/http"
)

func Init(db *storage.Storage, logging *slog.Logger, mux *http.ServeMux) {
	mux.HandleFunc("/current-lesson", currentLessonHandler(db, logging))
	mux.HandleFunc("/next-lesson", nextLessonHandler(db, logging))
	mux.HandleFunc("/schedule", scheduleHandler(db, logging))
}

func EnableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
