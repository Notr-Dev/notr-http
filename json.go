package notrhttp

import (
	"encoding/json"
	"net/http"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
}

func DecodeJson(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (rw *ResponseWriterWrapper) RespondWithJson(code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	rw.Write(response)
}

func (rw *ResponseWriterWrapper) RespondWithSuccess(payload interface{}) {
	rw.RespondWithJson(http.StatusOK, payload)
}

func (rw *ResponseWriterWrapper) RespondWithInternalError(message string) {
	rw.RespondWithJson(http.StatusInternalServerError, map[string]string{"message": message})
}

func (rw *ResponseWriterWrapper) RespondWithUnauthorized(message string) {
	rw.RespondWithJson(http.StatusUnauthorized, map[string]string{"message": message})
}

func (rw *ResponseWriterWrapper) RespondWithNotFound(message string) {
	rw.RespondWithJson(http.StatusNotFound, map[string]string{"message": message})
}
