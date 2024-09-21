package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	notrhttp "github.com/Notr-Dev/notr-http"
	"github.com/Notr-Dev/notr-http/services/db_service"
	"github.com/Notr-Dev/notr-http/services/logger_service"
)

func main() {
	server := notrhttp.NewServer(
		notrhttp.Server{
			Name:    "Test Server",
			Port:    ":8080",
			Version: "1.0.0",
		},
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

	mig := db_service.Migration{
		Up: func(db *sql.DB) error {
			_, err := db.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, log TEXT)")
			return err
		},
		Down: func(db *sql.DB) error {
			_, err := db.Exec("DROP TABLE test")
			return err
		},
	}

	DBServiceWrapper := db_service.NewDBService(db_service.DBServiceConfig{
		DBPath:     "test.sqlite",
		Migrations: []db_service.Migration{mig},
		Name:       "DB Service",
		Subpath:    "/db",
	})
	LoggerService := logger_service.NewLoggerService(DBServiceWrapper)

	server.RegisterJob(notrhttp.Job{
		Name:     "Test Job",
		Interval: 5 * time.Second,
		Job: func() error {
			fmt.Println("Test Job")
			return nil
		},
	})

	server.RegisterService(DBServiceWrapper.Service)
	server.RegisterService(LoggerService)

	server.ServeStatic("/static", "photos")

	err := server.Run()
	if err != nil {
		panic(err)
	}
}
