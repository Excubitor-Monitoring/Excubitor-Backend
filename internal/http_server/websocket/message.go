package websocket

import (
	"encoding/json"
	"time"
)

type TargetAddress string

func GetEmptyTarget() TargetAddress {
	return ""
}

type Message struct {
	OpCode OpCode        `json:"op"`
	Target TargetAddress `json:"target"`
	Value  string        `json:"value,omitempty"`
}

func NewMessage(opcode OpCode, target TargetAddress, value string) Message {
	return Message{opcode, target, value}
}

func (msg Message) Bytes() ([]byte, error) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

type HistoryRequestParameters struct {
	From       time.Time   `json:"from,omitempty"`
	Until      time.Time   `json:"until,omitempty"`
	MaxDensity interface{} `json:"max_density,omitempty"`
}
