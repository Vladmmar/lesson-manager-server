package server

import (
	"fmt"
	"lesson-manager-server/internal/app"
	"lesson-manager-server/internal/config"
	"net/http"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//init logger
	logging := app.SetupLogger(cfg)
	logging.Info("Initialized logger")

	//init storage
	db, err := app.SetupStorage(cfg.Db)
	if err != nil {
		logging.Error("Could not connect to database: ", err.Error())
		return
	}
	if err = db.Db.Ping(); err != nil {
		logging.Error(err.Error())
	}

	//setup http server

	logging.Info(fmt.Sprintf("Starting http server on port %s", cfg.Net.Port))
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Net.Host, cfg.Net.Port), nil)
	if err != nil {
		logging.Error(err.Error())
	} else {
		logging.Info("Successfully started http server on port %s", cfg.Net.Port)
	}
}
