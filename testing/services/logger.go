package services

import (
	"dev/types"
	"fmt"

	notrhttp "github.com/Notr-Dev/notr-http"
)

var LoggerService = notrhttp.NewService[*types.Data]("Logger")

func init() {
	LoggerService.AddDependency(&DBService)
	LoggerService.SetInitFunction(func(service *notrhttp.Service[*types.Data], server *notrhttp.Server[*types.Data]) error {

		fmt.Println("Initializing logger")

		_, err := server.Data.DB.Exec("CREATE TABLE IF NOT EXISTS logs (id INTEGER PRIMARY KEY, log TEXT)")

		if err != nil {
			return err
		}

		return nil
	})
}
