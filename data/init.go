package data

import (
	"database/sql"
	"os"

	"github.com/pressly/goose/v3"

	_ "github.com/mattn/go-sqlite3"
	"github.com/odin-software/metro/control"
)

const (
	dbPath         = "./data/metro.db"
	migrationsPath = "./data/sql/migrations"
	seedsPath      = "./data/sql/seeds"
)

// InitDatabase checks if the database exists and runs migrations if needed.
func InitDatabase() error {
	dbExists := checkDatabaseExists()

	if !dbExists {
		control.Log("Database not found, creating and running migrations...")
	} else {
		control.Log("Database exists, checking for pending migrations...")
	}

	// Always run migrations (with versioning - will skip if already applied)
	if err := runMigrations(); err != nil {
		return err
	}
	control.Log("Migrations completed successfully")

	// Note: Seeds are NOT run automatically - use make seed_test_city or make seed_santo_domingo
	// to choose which city to load before running the program

	return nil
}

// checkDatabaseExists checks if the database file exists.
func checkDatabaseExists() bool {
	_, err := os.Stat(dbPath)
	return !os.IsNotExist(err)
}

// runMigrations runs all pending goose migrations.
func runMigrations() error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	goose.SetDialect("sqlite3")

	if err := goose.Up(db, migrationsPath); err != nil {
		return err
	}

	return nil
}

// runSeeds runs all seed files (useful for development).
func runSeeds() error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	control.Log("Running seeds...")
	goose.SetDialect("sqlite3")

	if err := goose.Up(db, seedsPath, goose.WithNoVersioning()); err != nil {
		return err
	}

	return nil
}
