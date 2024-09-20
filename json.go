package notrhttp

import (
	"encoding/json"
	"net/http"
)

func DecodeJson(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithSuccess(w http.ResponseWriter, payload interface{}) {
	RespondWithJson(w, http.StatusOK, payload)
}

func RespondWithInternalError(w http.ResponseWriter, message string) {
	RespondWithJson(w, http.StatusInternalServerError, map[string]string{"message": message})
}

func RespondWithUnauthorized(w http.ResponseWriter, message string) {
	RespondWithJson(w, http.StatusUnauthorized, map[string]string{"message": message})
}

func RespondWithNotFound(w http.ResponseWriter, message string) {
	RespondWithJson(w, http.StatusNotFound, map[string]string{"message": message})
}
