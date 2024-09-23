package auth_service

import (
	"database/sql"

	notrhttp "github.com/Notr-Dev/notr-http"
	"github.com/Notr-Dev/notr-http/services/db_service"
)

type AuthServiceConfig struct {
	Name      string
	Subpath   string
	JWTConfig JWTConfig
}

type JWTConfig struct {
	Issuer string
	Secret []byte
	Types  []string
	Roles  []string
}

type AuthService struct {
	*notrhttp.Service
	DBService *db_service.DBService
	JWTConfig JWTConfig
}

func NewAuthService(config AuthServiceConfig, dbService *db_service.DBService) *notrhttp.Service {
	wrapper := &AuthService{}
	wrapper.DBService = dbService
	return notrhttp.NewService(
		notrhttp.Service{
			PackageID:    "auth",
			Name:         config.Name,
			Path:         config.Subpath,
			Dependencies: []*notrhttp.Service{dbService.Service},
			InitFunction: func(service *notrhttp.Service, server *notrhttp.Server) error {
				err := dbService.AddMigrations(
					db_service.Migration{
						ID: "auth-users",
						Up: func(db *sql.DB) error {
							_, err := db.Exec(`
								CREATE TABLE users (
									id BLOB PRIMARY KEY,
									created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
									email TEXT NOT NULL,
									password TEXT NOT NULL,
									role TEXT NOT NULL
								)
							`)
							return err
						},
						Down: func(db *sql.DB) error {
							_, err := db.Exec("DROP TABLE users")
							return err
						},
					},
				)
				if err != nil {
					return err
				}
				return nil
			},
			Routes: []notrhttp.Route{},
		},
	)
}
