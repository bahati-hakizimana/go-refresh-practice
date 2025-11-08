package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func PulseJson(r *http.Request, payload any) error {
	if r.Body == nil {
		fmt.Errorf("Missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJson(w http.ResponseWriter, status int, v any) error {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)

}

func WriteError(w http.ResponseWriter, status int, error error){
	WriteJson(w, status, map[string]string{"error":error.Error()})
}