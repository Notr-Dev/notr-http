package logger_service

import (
	"database/sql"
	"fmt"

	notrhttp "github.com/Notr-Dev/notr-http"
	"github.com/Notr-Dev/notr-http/services/db_service"
)

func NewLoggerService(dbService *db_service.DBService) *notrhttp.Service {
	return notrhttp.NewService(
		notrhttp.Service{
			Name:         "Logger",
			Dependencies: []*notrhttp.Service{dbService.Service},
			InitFunction: func(service *notrhttp.Service, server *notrhttp.Server) error {
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
			},
		},
	)
}
