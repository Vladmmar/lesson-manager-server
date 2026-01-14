package server

import (
	"lesson-manager-server/internal/app"
	"lesson-manager-server/internal/config"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//init logger
	slog := app.SetupLogger(cfg)
	slog.Info("Initialized logger")

	//init storage
	db, err := app.SetupStorage(cfg.Db)
	if err != nil {
		slog.Error("Could not connect to database: ", err.Error())
		return
	}
	if err = db.Db.Ping(); err != nil {
		slog.Error(err.Error())
	}
}
