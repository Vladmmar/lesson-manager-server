package server

import (
	. "lesson-manager-server/internal/app"
	"log"
)

func main() {
	//load config
	cfg, err = SetupConfig()
	if err != nil {
		log.Fatalln(err)
	}

	//init logger
	slog := SetupLogger(cfg)
	slog.Info("Initialized logger")

	//init storage

}
