package services

import (
	"fmt"

	notrhttp "github.com/Notr-Dev/notr-http"
)

var LoggerService = notrhttp.NewService("Logger")

func init() {
	LoggerService.AddDependency(&DBService)
	LoggerService.SetInitFunction(func(s *notrhttp.Service) error {

		fmt.Println("Initializing logger")

		_, err := db.Exec("CREATE TABLE IF NOT EXISTS logs (id INTEGER PRIMARY KEY, log TEXT)")

		if err != nil {
			return err
		}

		return nil
	})
}
