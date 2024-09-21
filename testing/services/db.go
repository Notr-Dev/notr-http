package services

import (
	"database/sql"
	"dev/types"

	notrhttp "github.com/Notr-Dev/notr-http"
)

var DBService = notrhttp.NewService[*types.Data]("DB")

func init() {
	DBService.SetInitFunction(func(service *notrhttp.Service[*types.Data], server *notrhttp.Server[*types.Data]) error {
		db, err := sql.Open("sqlite3", "./database.db")
		if err != nil {
			return err
		}

		server.Data.DB = db

		return nil
	})
}
