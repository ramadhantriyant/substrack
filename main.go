package main

import (
	"context"
	"embed"
	"log"
	"path/filepath"
	"time"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/database"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/models"
)

//go:embed sql/schema/*.sql
var embedMigrations embed.FS

const (
	dataDir          = "data"
	dbFileName       = "substrack.db"
	port             = ":8080"
	shutdownTimeout  = 10 * time.Second
	dbConnectionPool = 1
)

func main() {
	dbPath := filepath.Join(dataDir, dbFileName)
	db, err := connectDatabase(dbPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Running migrations...")
	if err := runMigrations(db); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	log.Println("Migrations complete")
	log.Printf("Database initialized in %s", dbPath)

	config := &models.AppConfig{
		DB:      db,
		Queries: database.New(db),
	}
	server := createServer(config, port)

	if err := runServer(context.Background(), server, shutdownTimeout); err != nil {
		log.Fatalf("Cannot bind to 0.0.0.0%s", err)
	}
}
