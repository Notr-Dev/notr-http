package services

import (
	"database/sql"

	notrhttp "github.com/Notr-Dev/notr-http"
)

type DBService struct {
	*notrhttp.Service
	Database *sql.DB
}

func (d *DBService) GetDB() *sql.DB {
	if d.Database == nil {
		panic("Database is not initialized")
	}
	return d.Database
}

func NewDBService(dbPath string) *DBService {
	wrapper := &DBService{}
	service := notrhttp.NewService(
		notrhttp.WithServiceName("DB"),
		notrhttp.WithServiceInitFunction(func(service *notrhttp.Service) error {
			db, err := sql.Open("sqlite3", dbPath)
			if err != nil {
				return err
			}

			if err := db.Ping(); err != nil {
				return err
			}

			wrapper.Database = db

			return nil
		}),
	)

	wrapper.Service = service

	return wrapper
}
