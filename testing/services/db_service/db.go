package db_service

import (
	"database/sql"

	notrhttp "github.com/Notr-Dev/notr-http"
)

type DBServiceConfig struct {
	Name       string
	DBPath     string
	Subpath    string
	Migrations []Migration
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

func NewDBService(config DBServiceConfig) *DBService {

	if config.Name == "" {
		config.Name = "Unnamed DB Service"
	}

	if config.Subpath == "" {
		panic("Subpath is required")
	}

	if config.DBPath == "" {
		panic("DBPath is required")
	}

	if config.Migrations == nil {
		config.Migrations = make([]Migration, 0)
	}

	wrapper := &DBService{}
	wrapper.Migrations = make([]versionedMigration, 0)
	service := notrhttp.NewService(
		notrhttp.WithServiceName(config.Name),
		notrhttp.WithServiceSubpath(config.Subpath),
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
		notrhttp.WithServiceRoutes(
			notrhttp.Route{
				Method: "GET",
				Path:   "/",
				Handler: func(rw notrhttp.Writer, r *notrhttp.Request) {
					_, err := wrapper.GetDB().Exec("INSERT INTO test (log) VALUES ('test')")
					if err != nil {
						rw.RespondWithInternalError(err.Error())
						return
					}
					rw.RespondWithSuccess(map[string]interface{}{
						"message": "Welcome to the db service.",
					})
				},
			}),
	)

	wrapper.Service = service

	return wrapper
}
