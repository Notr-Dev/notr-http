package dash_service_ui

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:generate npm i
//go:generate npm run build
//go:embed build/*
var files embed.FS

func BuildHTTPFS() http.FileSystem {
	build, err := fs.Sub(files, "build")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(build)
}
