package http_server

import "time"

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
