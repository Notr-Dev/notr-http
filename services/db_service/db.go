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
	Migrations []Migration
}

func (d *DBService) GetDB() *sql.DB {
	if d.Database == nil {
		panic("Database is not initialized")
	}
	return d.Database
}

func NewDBService(config DBServiceConfig) *DBService {
	if config.DBPath == "" {
		panic("DBPath is required")
	}

	if config.Migrations == nil {
		config.Migrations = make([]Migration, 0)
	}

	wrapper := &DBService{}
	wrapper.Migrations = make([]Migration, 0)
	service := notrhttp.NewService(
		notrhttp.Service{
			PackageID: "db",
			Name:      config.Name,
			Path:      config.Subpath,
			InitFunction: func(service *notrhttp.Service, server *notrhttp.Server) error {
				db, err := sql.Open("sqlite3", config.DBPath)
				if err != nil {
					return err
				}

				if err := db.Ping(); err != nil {
					return err
				}

				wrapper.Database = db

				err = wrapper.AddMigrations(config.Migrations...)

				return err
			},
			Routes: []notrhttp.Route{
				{
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
				},
			},
		},
	)
	wrapper.Service = service

	return wrapper
}
