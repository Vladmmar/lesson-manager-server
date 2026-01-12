package storage

import (
	"database/sql"
)

type Storage struct {
	Db *sql.DB
}
