package websocket

import (
	"encoding/json"
	"time"
)

// TargetAddress models the target field in websocket messages
type TargetAddress string

// GetEmptyTarget returns an empty TargetAddress.
// This is used whenever a message needs to be sent that cannot be assigned to a specific target.
// I.e. whenever there is no target specified in an erroneous request
func GetEmptyTarget() TargetAddress {
	return ""
}

// Message is the model for a websocket message.
type Message struct {
	OpCode OpCode        `json:"op"`
	Target TargetAddress `json:"target"`
	Value  string        `json:"value,omitempty"`
}

// NewMessage returns a Message with chosen OpCode, TargetAddress and value.
func NewMessage(opcode OpCode, target TargetAddress, value string) Message {
	return Message{opcode, target, value}
}

// Bytes converts a Message to JSON and returns it as a byte slice.
func (msg Message) Bytes() ([]byte, error) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

// HistoryRequestParameters is a model for paramters that can be set with HIST requests.
type HistoryRequestParameters struct {
	From       time.Time   `json:"from,omitempty"`
	Until      time.Time   `json:"until,omitempty"`
	MaxDensity interface{} `json:"max_density,omitempty"`
}
