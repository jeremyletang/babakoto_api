package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJsonResponse(w http.ResponseWriter, status int, body interface{}) {
	rawBody, _ := json.Marshal(body)
	WriteResponse(w, status, string(rawBody))
}

func WriteResponse(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(body)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `%v`, body)
}

func ReadRequestBody(r *http.Request, t interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(t)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}
