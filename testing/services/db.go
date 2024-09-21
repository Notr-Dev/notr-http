package services

import (
	"dev/globals"

	notrhttp "github.com/Notr-Dev/notr-http"
)

var DBService = notrhttp.NewService(
	notrhttp.WithServiceName("DB Service"),
	notrhttp.WithServiceInitFunction(func(service *notrhttp.Service, server *notrhttp.Server) error {
		_, err := globals.Data.DB.Exec("CREATE TABLE IF NOT EXISTS logs (id INTEGER PRIMARY KEY, log TEXT)")

		if err != nil {
			return err
		}

		return nil
	}),
)
