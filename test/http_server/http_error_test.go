package http_server

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/http_server"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewError(t *testing.T) {
	httpError := http_server.NewHTTPError("Some message", "Something.Something")

	assert.Equal(t, "Some message", httpError.Message)
	assert.Equal(t, "Something.Something", httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}
