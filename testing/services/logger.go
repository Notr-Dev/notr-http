package services

import (
	"dev/globals"
	"fmt"

	notrhttp "github.com/Notr-Dev/notr-http"
)

var LoggerService = notrhttp.NewService(
	notrhttp.WithServiceName("Logger Service"),
	notrhttp.WithServiceInitFunction(func(service *notrhttp.Service, server *notrhttp.Server) error {

		fmt.Println("Initializing logger")

		_, err := globals.Data.DB.Exec("CREATE TABLE IF NOT EXISTS logs (id INTEGER PRIMARY KEY, log TEXT)")

		if err != nil {
			return err
		}

		return nil
	}),
	notrhttp.WithServiceDependencies(DBService),
)
