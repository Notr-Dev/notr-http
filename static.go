package notrhttp

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func (s *Server) StaticServe(path string, dir string) {
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
