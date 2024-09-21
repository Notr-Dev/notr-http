package services

import (
	"fmt"

	notrhttp "github.com/Notr-Dev/notr-http"
)

func NewLoggerService(dbService *DBService) *notrhttp.Service {
	return notrhttp.NewService(
		notrhttp.WithServiceName("Logger"),
		notrhttp.WithServiceDependencies(dbService.Service),
		notrhttp.WithServiceInitFunction(func(service *notrhttp.Service) error {

			fmt.Println("Initializing logger")

			_, err := dbService.GetDB().Exec("CREATE TABLE IF NOT EXISTS logs (id INTEGER PRIMARY KEY, log TEXT)")

			if err != nil {
				return err
			}

			return nil
		}),
	)
}
