package services

import (
	"database/sql"

	notrhttp "github.com/Notr-Dev/notr-http"
)

func NewDBService(dbPath string, database **sql.DB) *notrhttp.Service {
	return notrhttp.NewService(
		notrhttp.WithServiceName("DB"),
		notrhttp.WithServiceInitFunction(func(service *notrhttp.Service) error {
			db, err := sql.Open("sqlite3", dbPath)
			if err != nil {
				return err
			}

			if err := db.Ping(); err != nil {
				return err
			}

			*database = db

			return nil
		}),
	)
}