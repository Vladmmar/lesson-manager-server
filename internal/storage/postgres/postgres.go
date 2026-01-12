package postgres

import (
	"database/sql"
	"fmt"
	"lesson-manager-server/internal/config"
	"lesson-manager-server/internal/storage"

	_ "github.com/lib/pq"
)

func New(cfg *config.Database) (*storage.Storage, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password,
		cfg.Host, cfg.Port,
		cfg.Name,
	))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &storage.Storage{Db: db}, nil
}
