package helper

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewError(t *testing.T) {
	httpError := NewHTTPError("Some message", "Something.Something")

	assert.Equal(t, "Some message", httpError.Message)
	assert.Equal(t, "Something.Something", httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}

func TestReturnError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/samplePath", nil)
	w := httptest.NewRecorder()

	ReturnError(w, req, 400, "Some reason")

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response Error
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "Some reason", response.Message)
	assert.Equal(t, "/samplePath", response.Path)
}
