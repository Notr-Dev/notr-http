package services

import (
	"database/sql"
	"dev/services/db_service"
	"fmt"

	notrhttp "github.com/Notr-Dev/notr-http"
)

func NewLoggerService(dbService *db_service.DBService) *notrhttp.Service {
	return notrhttp.NewService(
		notrhttp.WithServiceName("Logger"),
		notrhttp.WithServiceDependencies(dbService.Service),
		notrhttp.WithServiceInitFunction(func(service *notrhttp.Service) error {

			fmt.Println("Initializing logger")

			return dbService.AddMigrations(
				db_service.Migration{
					Up: func(db *sql.DB) error {
						_, err := db.Exec("CREATE TABLE logs (id INTEGER PRIMARY KEY, log TEXT)")
						return err
					},
					Down: func(db *sql.DB) error {
						_, err := db.Exec("DROP TABLE logs")
						return err
					},
				},
			)
		}),
	)
}
