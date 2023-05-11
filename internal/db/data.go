package db

import (
	"time"
)

type HistoryMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Message   struct {
		Target string `json:"target"`
		Value  string `json:"value"`
	} `json:"message"`
}

type History []HistoryMessage
