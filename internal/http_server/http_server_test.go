package http_server

import (
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInfo(t *testing.T) {
	ctx.GetContext().RegisterModule(ctx.NewModule("TestModule"))

	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	w := httptest.NewRecorder()

	info(w, req)

	res := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.JSONEq(t, `{"authentication": { "method": "PAM" }, "modules": [ { "name": "TestModule" } ] }`, string(body))
}
