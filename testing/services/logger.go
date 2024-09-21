package services

import (
	"database/sql"
	"fmt"

	notrhttp "github.com/Notr-Dev/notr-http"
)

func NewLoggerService(dbService *notrhttp.Service, database func() *sql.DB) *notrhttp.Service {
	return notrhttp.NewService(
		notrhttp.WithServiceName("Logger"),
		notrhttp.WithServiceDependencies(dbService),
		notrhttp.WithServiceInitFunction(func(service *notrhttp.Service) error {

			fmt.Println("Initializing logger")

			_, err := database().Exec("CREATE TABLE IF NOT EXISTS logs (id INTEGER PRIMARY KEY, log TEXT)")

			if err != nil {
				return err
			}

			return nil
		}),
	)
}
