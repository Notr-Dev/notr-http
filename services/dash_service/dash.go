package dash_service

import (
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
				server.ServeHttpFileSystem(config.Subpath, dash_service_ui.BuildHTTPFS())

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
			},
		},
	)
	wrapper.Service = service

	return wrapper
}
