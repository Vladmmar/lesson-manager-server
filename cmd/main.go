package main

import (
	"fmt"
	"lesson-manager-server/internal/app"
	"lesson-manager-server/internal/config"
	"lesson-manager-server/internal/http/handlers"
	"net/http"
)

func main() {
	//load config
	cfg := config.MustLoad()
	fmt.Println("Successfully loaded config")

	//init logger
	logging := app.SetupLogger(cfg)
	logging.Info("Successfully initialized logger")

	//init storage
	db, err := app.SetupStorage(&cfg.Db)
	if err != nil {
		logging.Error("Could not connect to database: ", err.Error())
		return
	}
	if err = db.Db.Ping(); err != nil {
		logging.Error(err.Error())
		return
	}
	logging.Info("Successfully established a connection to the database")

	//setup http server

	logging.Info(fmt.Sprintf("Starting http server on port %s", cfg.Net.Port))
	handlers.Init(db, logging)
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Net.Host, cfg.Net.Port), nil)
	if err != nil {
		logging.Error(err.Error())
	} else {
		logging.Info("Successfully started http server on port %s", cfg.Net.Port)
	}
}
