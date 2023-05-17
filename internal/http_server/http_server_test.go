package http_server

import (
	"fmt"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestInfo(t *testing.T) {
	ctx.GetContext().RegisterModule(ctx.NewModule("TestModule", func() {}))

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

func TestInfoMethodNotAllowed(t *testing.T) {
	type testParams struct {
		description string
		method      string
	}

	for _, params := range []testParams{
		{
			description: "Method POST",
			method:      http.MethodPost,
		},
		{
			description: "Method PUT",
			method:      http.MethodPut,
		},
		{
			description: "Method PATCH",
			method:      http.MethodPatch,
		},
		{
			description: "Method DELETE",
			method:      http.MethodDelete,
		},
		{
			description: "Method TRACE",
			method:      http.MethodTrace,
		},
	} {
		t.Run(params.description, func(t *testing.T) {
			req := httptest.NewRequest(params.method, "/info", nil)
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

			httpError := parseHTTPError(body)

			assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
			assert.Equal(t, "/info", httpError.Path)
			assert.Equal(t, fmt.Sprintf("Method %s not allowed!", params.method), httpError.Message)
			assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
		})
	}
}
