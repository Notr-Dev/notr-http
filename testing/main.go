package main

import (
	"dev/services"

	_ "github.com/mattn/go-sqlite3"

	notrhttp "github.com/Notr-Dev/notr-http"
)

func main() {
	server := notrhttp.NewServer("8080", "1.0")
	server.SetName("My Server")
	server.Post("/log", func(rw notrhttp.Writer, r *notrhttp.Request) {
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

	server.RegisterService(services.DBService)
	server.RegisterService(services.LoggerService)

	err := server.Run()
	if err != nil {
		panic(err)
	}
}
