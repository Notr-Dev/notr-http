package db_service

import (
	"database/sql"
	"fmt"
)

type Migration struct {
	ID   string
	Up   func(db *sql.DB) error
	Down func(db *sql.DB) error
}

var initialMigration = Migration{
	ID: "initial",
	Up: func(db *sql.DB) error {
		fmt.Println("Creating migrations table...")
		_, err := db.Exec(`
		CREATE TABLE migrations (
			id TEXT NOT NULL,
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

	for _, mig := range migrations {
		for _, existingMig := range d.Migrations {
			if existingMig.ID == mig.ID {
				return fmt.Errorf("Migration with ID %s already exists", mig.ID)
			}
		}
	}

	fmt.Printf("Adding %d migrations\n", len(migrations))

	d.Migrations = append(d.Migrations, migrations...)

	for _, mig := range d.Migrations {
		wasApplied, err := checkIfMigrationWasApplied(d.Database, mig.ID)
		if err != nil {
			return err
		}

		if !wasApplied {
			err := mig.Up(d.Database)
			if err != nil {
				return err
			}
			err = insertMigrationRecord(d.Database, mig.ID)
			if err != nil {
				return err
			}

			fmt.Println("Migration", mig.ID, "completed")
		}
	}

	return nil
}

func insertMigrationRecord(db *sql.DB, id string) error {
	_, err := db.Exec("INSERT INTO migrations (id) VALUES (?)", id)
	if err != nil {
		return err
	}
	return nil
}

func checkIfMigrationWasApplied(db *sql.DB, id string) (bool, error) {
	row := db.QueryRow("SELECT id FROM migrations WHERE id = ?", id)
	var result string
	err := row.Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
