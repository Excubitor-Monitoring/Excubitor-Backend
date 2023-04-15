package http_server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMessage(t *testing.T) {
	msg := newMessage(UNSUB, "Some.Target.Address", "Some value")
	assert.Equal(t, UNSUB, msg.OpCode)
	assert.Equal(t, TargetAddress("Some.Target.Address"), msg.Target)
	assert.Equal(t, "Some value", msg.Value)
}
