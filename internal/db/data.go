package db

import (
	"time"
)

// HistoryMessage serves as a model to describe entries in the history table.
type HistoryMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Message   struct {
		Target string `json:"target"`
		Value  string `json:"value"`
	} `json:"message"`
}

type History []HistoryMessage
