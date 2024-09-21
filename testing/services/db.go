package services

import (
	"database/sql"
	"dev/state"
	"fmt"

	notrhttp "github.com/Notr-Dev/notr-http"
)

var DBService = notrhttp.NewService("db", "/db")

func init() {
	DBService.SetInitFunction(func(s *notrhttp.Service) error {
		fmt.Println("Initializing db")
		database, err := sql.Open("sqlite3", "test.sqlite")
		if err != nil {
			return err
		}

		err = database.Ping()
		if err != nil {
			return err
		}

		state.DB = database

		return nil
	})
}
