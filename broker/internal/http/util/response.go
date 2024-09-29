package util

import (
	"encoding/json"
	"net/http"
)

var defaultHeaders = map[string]string{
	"Content-Type":           "application/json; charset=utf-8",
	"X-Content-Type-Options": "nosniff",
}

func Respond(w http.ResponseWriter, data any, statusCode int) {
	setDefaultHeaders(w)
	w.WriteHeader(statusCode)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func setDefaultHeaders(w http.ResponseWriter) {
	for k, v := range defaultHeaders {
		w.Header().Set(k, v)
	}
}
