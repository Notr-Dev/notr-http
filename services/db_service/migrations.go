package db_service

import (
	"database/sql"
	"fmt"
	"strings"
)

type Migration struct {
	Up   func(db *sql.DB) error
	Down func(db *sql.DB) error
}

type versionedMigration struct {
	Version   int
	Migration Migration
}

var initialMigration = Migration{
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

func (d *DBService) AddMigrations(migrations ...Migration) error {

	if len(migrations) == 0 {
		return nil
	}

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

	fmt.Println("Latest version:", latestVersion)

	for _, mig := range d.Migrations {
		if mig.Version > latestVersion {
			fmt.Println("Running migration", mig.Version)
			err := mig.Migration.Up(d.Database)
			if err != nil {
				return err
			}
			fmt.Println("Migration", mig.Version, "completed")
			err = insertMigrationRecord(d.Database, mig.Version)
			if err != nil {

				return err
			}
			fmt.Println("Migration record inserted")
		}
	}

	return nil
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
