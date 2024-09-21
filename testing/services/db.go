package services

import (
	"database/sql"

	notrhttp "github.com/Notr-Dev/notr-http"
)

var db *sql.DB

var DBService = notrhttp.NewService("DB")

func init() {
	DBService.SetInitFunction(func(s *notrhttp.Service) error {
		var err error
		db, err = sql.Open("sqlite3", "./database.db")
		if err != nil {
			return err
		}

		return nil
	})
}
