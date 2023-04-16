package http_server

import (
	"encoding/json"
	"net/http"
	"time"
)

type Error struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
}

func NewHTTPError(message string, path string) Error {
	return Error{
		time.Now(),
		message,
		path,
	}
}

func ReturnError(w http.ResponseWriter, r *http.Request, status int, reason string) {
	w.WriteHeader(status)

	httpError := NewHTTPError(reason, r.RequestURI)

	bytes, err := json.Marshal(httpError)
	if err != nil {
		logger.Error("Couldn't encode http error!")
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		logger.Error("Couldn't write http error!")
		return
	}
}
