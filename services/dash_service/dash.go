package dash_service

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"runtime"

	notrhttp "github.com/Notr-Dev/notr-http"
	dash_service_ui "github.com/Notr-Dev/notr-http/services/dash_service/web"
)

type DashServiceConfig struct {
	Name    string
	Subpath string
}

type DashService struct {
	*notrhttp.Service
	server *notrhttp.Server
}

func NewDashService(config DashServiceConfig) *DashService {
	if config.Subpath == "" {
		panic("Subpath is required")
	}
	if config.Subpath[0] != '/' {
		panic("Subpath must start with a /")
	}

	wrapper := &DashService{}
	service := notrhttp.NewService(
		notrhttp.Service{
			Name: config.Name,
			Path: config.Subpath,
			InitFunction: func(service *notrhttp.Service, server *notrhttp.Server) error {
				wrapper.server = server

				return nil
			},
			Routes: []notrhttp.Route{
				{
					Method: "GET",
					Path:   "/details",
					Handler: func(rw notrhttp.Writer, r *notrhttp.Request) {
						var m runtime.MemStats
						runtime.ReadMemStats(&m)
						rw.RespondWithSuccess(map[string]interface{}{
							"name":    wrapper.server.Name,
							"version": wrapper.server.Version,
							"port":    wrapper.server.Port,
							"memory": map[string]interface{}{
								"alloc":      float64(m.Alloc) / 1024 / 1024,
								"totalAlloc": float64(m.TotalAlloc) / 1024 / 1024,
								"sys":        float64(m.Sys) / 1024 / 1024,
								"numGC":      m.NumGC,
							},
						})
					},
				},
				{
					Method: "GET",
					Path:   "/{filename...}",
					Handler: func(rw notrhttp.Writer, r *notrhttp.Request) {
						filename := r.Params["filename"]
						if filename == "" {
							filename = "index.html"
						}
						file, err := dash_service_ui.BuildHTTPFS().Open(filename)
						if err != nil {
							http.Error(rw, "File not found", http.StatusNotFound)
							return
						}
						defer file.Close()

						mimType := "application/octet-stream"
						if filepath.Ext(filename) == ".html" {
							fmt.Printf("Serving html file: %s\n", filename)
							mimType = "text/html"
						}
						if filepath.Ext(filename) == ".css" {
							fmt.Printf("Serving css file: %s\n", filename)
							mimType = "text/css"
						}
						if filepath.Ext(filename) == ".js" {
							fmt.Printf("Serving js file: %s\n", filename)
							mimType = "application/javascript"
						}

						rw.Header().Set("Content-Type", mimType)

						if _, err := io.Copy(rw, file); err != nil {
							http.Error(rw, "Error serving file", http.StatusInternalServerError)
						}
					},
				},
			},
		},
	)
	wrapper.Service = service

	return wrapper
}
