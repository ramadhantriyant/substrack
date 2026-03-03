package models

import (
	"database/sql"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/database"
)

type AppConfig struct {
	DB        *sql.DB
	Queries   *database.Queries
	JWTSecret string
}
