package notrhttp

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func (s *Server) ServeStatic(path string, dir string) {
	s.Routes = append(s.Routes,
		Route{
			Method: http.MethodGet,
			Path:   path + "/{filename}",
			Handler: func(rw Writer, r *Request) {
				filename := r.Params["filename"]
				filePath := filepath.Join(dir, filename)
				fmt.Println("Serving file: ", filePath)

				file, err := os.Open(filePath)
				if err != nil {
					http.Error(rw, "File not found", http.StatusNotFound)
					return
				}
				defer file.Close()

				rw.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(filePath)))

				if _, err := io.Copy(rw, file); err != nil {
					http.Error(rw, "Error serving file", http.StatusInternalServerError)
				}
			},
		},
	)
}

func (s *Server) ServeHttpFileSystem(path string, fs http.FileSystem) {
	s.Routes = append(s.Routes, Route{
		Method: http.MethodGet,
		Path:   path + "/{filename...}",
		Handler: func(rw Writer, r *Request) {
			filename := r.Params["filename"]
			if filename == "" {
				filename = "index.html"
			}
			file, err := fs.Open(filename)
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
	})
}
