package services

import (
	"database/sql"
	"fmt"
	"strings"

	notrhttp "github.com/Notr-Dev/notr-http"
)

type DBServiceConfig struct {
	DBPath     string
	Migrations []Migration
}

type Migration struct {
	Up   func(db *sql.DB) error
	Down func(db *sql.DB) error
}

type versionedMigration struct {
	Version   int
	Migration Migration
}

type DBService struct {
	*notrhttp.Service
	Database   *sql.DB
	Migrations []versionedMigration
}

func (d *DBService) GetDB() *sql.DB {
	if d.Database == nil {
		panic("Database is not initialized")
	}
	return d.Database
}

func (d *DBService) AddMigrations(migrations ...Migration) error {

	fmt.Printf("Adding %d migrations\n", len(migrations))

	toAdd := make([]versionedMigration, 0)

	for index, mig := range migrations {
		toAdd = append(toAdd, versionedMigration{
			Version:   len(d.Migrations) + index,
			Migration: mig,
		})
	}

	d.Migrations = append(d.Migrations, toAdd...)

	latestVersion, err := getLatestVersion(d.Database)
	if err != nil {
		return err
	}

	for _, mig := range d.Migrations {
		if mig.Version > latestVersion {
			err := mig.Migration.Up(d.Database)
			if err != nil {
				return err
			}
			err = insertMigrationRecord(d.Database, mig.Version)
			if err != nil {

				return err
			}
		}
	}

	return nil
}

func NewDBService(config DBServiceConfig) *DBService {

	if config.DBPath == "" {
		panic("DBPath is required")
	}

	if config.Migrations == nil {
		panic("Migrations are required")
	}

	// Append a migration to the beginning of migrations
	initialMigration := Migration{
		Up: func(db *sql.DB) error {
			fmt.Println("Creating migrations table...")
			_, err := db.Exec(`
			CREATE TABLE migrations (
				version INTEGER NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
			`)
			return err
		},
		Down: func(db *sql.DB) error {
			_, err := db.Exec(`
			DROP TABLE migrations;
			`)
			return err
		},
	}

	wrapper := &DBService{}
	wrapper.Migrations = make([]versionedMigration, 0)
	service := notrhttp.NewService(
		notrhttp.WithServiceName("DB"),
		notrhttp.WithServiceInitFunction(func(service *notrhttp.Service) error {
			db, err := sql.Open("sqlite3", config.DBPath)
			if err != nil {
				return err
			}

			if err := db.Ping(); err != nil {
				return err
			}

			wrapper.Database = db
			err = wrapper.AddMigrations(initialMigration)
			if err != nil {
				return err
			}
			err = wrapper.AddMigrations(config.Migrations...)

			return err
		}),
	)

	wrapper.Service = service

	return wrapper
}

func getLatestVersion(db *sql.DB) (int, error) {
	var version int
	row := db.QueryRow("SELECT version FROM migrations ORDER BY version DESC LIMIT 1")
	err := row.Scan(&version)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		if strings.Contains(err.Error(), "no such table: migrations") {
			return -1, nil
		}

		return 0, err
	}
	return version, nil
}

func insertMigrationRecord(db *sql.DB, version int) error {
	_, err := db.Exec("INSERT INTO migrations (version) VALUES (?)", version)
	if err != nil {
		return err
	}
	return nil
}
