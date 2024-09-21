package main

import (
	"dev/services"
	"dev/services/db_service"

	_ "github.com/mattn/go-sqlite3"

	notrhttp "github.com/Notr-Dev/notr-http"
)

func main() {
	server := notrhttp.NewServer(
		notrhttp.WithServerName("Test Server"),
		notrhttp.WithServerPort(":8080"),
		notrhttp.WithServerVersion("1.0.0"),
	)
	server.Post("/test", func(rw notrhttp.Writer, r *notrhttp.Request) {
		type Response struct {
			Password string `json:"password"`
			Log      string `json:"log"`
		}
		var body Response
		err := r.GetJSONBody(body)

		if err != nil {
			rw.RespondWithBadRequest("Invalid body")
			return
		}

		if body.Password != "lazar" {
			rw.RespondWithUnauthorized("Invalid password")
			return
		}

		rw.RespondWithSuccess(map[string]interface{}{
			"message": "Logged",
		})
	})

	DBServiceWrapper := db_service.NewDBService(db_service.DBServiceConfig{
		DBPath:     "test.sqlite",
		Migrations: make([]db_service.Migration, 0),
	})
	LoggerService := services.NewLoggerService(DBServiceWrapper)

	server.RegisterService(DBServiceWrapper.Service)
	server.RegisterService(LoggerService)

	err := server.Run()
	if err != nil {
		panic(err)
	}
}
