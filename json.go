package notrhttp

import (
	"encoding/json"
	"net/http"
)

type Writer struct {
	http.ResponseWriter
}

func (r *Request) GetJSONBody(v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (rw *Writer) RespondWithJson(code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	rw.Write(response)
}

func (rw *Writer) RespondWithSuccess(payload interface{}) {
	rw.RespondWithJson(http.StatusOK, payload)
}

func (rw *Writer) RespondWithInternalError(message string) {
	rw.RespondWithJson(http.StatusInternalServerError, map[string]string{"message": message})
}

func (rw *Writer) RespondWithUnauthorized(message string) {
	rw.RespondWithJson(http.StatusUnauthorized, map[string]string{"message": message})
}

func (rw *Writer) RespondWithNotFound(message string) {
	rw.RespondWithJson(http.StatusNotFound, map[string]string{"message": message})
}

func (rw *Writer) RespondWithBadRequest(message string) {
	rw.RespondWithJson(http.StatusBadRequest, map[string]string{"message": message})
}
