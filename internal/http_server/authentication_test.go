package http_server

import (
	"encoding/json"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/golang-jwt/jwt/v5"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	err := logging.SetDefaultLogger("CONSOLE")
	if err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestHandleAuthRequestInvalidMethod(t *testing.T) {

	methods := []string{http.MethodGet, http.MethodOptions, http.MethodDelete,
		http.MethodPatch, http.MethodHead, http.MethodConnect, http.MethodTrace}

	for _, m := range methods {
		req := httptest.NewRequest(m, "/auth", nil)
		w := httptest.NewRecorder()

		handleAuthRequest(w, req)

		res := w.Result()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)

		httpError := parseHTTPError(body)

		assert.Equal(t, "Method is not allowed!", httpError.Message)
		assert.Equal(t, "/auth", httpError.Path)
		assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)

		err = res.Body.Close()
		if err != nil {
			return
		}
	}
}

func TestHandleAuthRequestUndecipherableBody(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(""))
	req.RemoteAddr = "SampleAddress"
	w := httptest.NewRecorder()

	handleAuthRequest(w, req)

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

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "Can't decode message body!", httpError.Message)
	assert.Equal(t, "/auth", httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}

func TestHandleAuthRequestUnknownMethod(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/auth",
		strings.NewReader(`{ "method": "UnknownMethod", "credentials": { "something": "No credentials needed..." } }`))
	req.RemoteAddr = "SampleAddress"
	w := httptest.NewRecorder()

	handleAuthRequest(w, req)

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

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "Unsupported authentication method: UnknownMethod", httpError.Message)
	assert.Equal(t, "/auth", httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}

func TestHandleAuthRequestPAMInvalidCredentials(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/auth",
		strings.NewReader(`{ "method": "PAM", "credentials": { "username": "testuser", "password": "123456" } }`))
	req.RemoteAddr = "SampleAddress"
	w := httptest.NewRecorder()

	handleAuthRequest(w, req)

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

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, "Invalid username or password!", httpError.Message)
	assert.Equal(t, "/auth", httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}

func TestHandleRefreshRequestInvalidMethod(t *testing.T) {
	methods := []string{http.MethodGet, http.MethodOptions, http.MethodDelete,
		http.MethodPatch, http.MethodHead, http.MethodConnect, http.MethodTrace}

	for _, m := range methods {
		req := httptest.NewRequest(m, "/auth/refresh", nil)
		w := httptest.NewRecorder()

		handleRefreshRequest(w, req)

		res := w.Result()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)

		httpError := parseHTTPError(body)

		assert.Equal(t, "Method is not allowed!", httpError.Message)
		assert.Equal(t, "/auth/refresh", httpError.Path)
		assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)

		err = res.Body.Close()
		if err != nil {
			return
		}
	}
}

func TestHandleRefreshRequestNoHeader(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.RemoteAddr = "SampleAddress"
	w := httptest.NewRecorder()

	handleRefreshRequest(w, req)

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

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, "Bearer authentication is needed!", httpError.Message)
	assert.Equal(t, "Bearer", res.Header.Get("WWW-Authenticate"))
	assert.Equal(t, "/auth/refresh", httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}

func TestHandleRefreshRequestInvalidHeader(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.RemoteAddr = "SampleAddress"
	req.Header.Set("Authorization", "Basic dXNlcm5hbWU6cGFzc3dvcmQ=")
	w := httptest.NewRecorder()

	handleRefreshRequest(w, req)

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

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, "Bearer authentication is needed!", httpError.Message)
	assert.Equal(t, "Bearer", res.Header.Get("WWW-Authenticate"))
	assert.Equal(t, "/auth/refresh", httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}

func TestHandleRefreshRequestInvalidTokenExpired(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
		return
	}

	err = k.Load(confmap.Provider(map[string]interface{}{
		"http.auth.jwt.accessTokenSecret":  "123456",
		"http.auth.jwt.refreshTokenSecret": "abcdef",
	}, "."), nil)
	if err != nil {
		t.Error(err)
		return
	}

	token, err := signRefreshToken(jwt.MapClaims{
		"iss": "excubitor-backend",
		"sub": "testuser",
		"exp": time.Now().Add(-4 * time.Hour).Unix(),
	})
	if err != nil {
		t.Error(err)
		return
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.RemoteAddr = "SampleAddress"
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handleRefreshRequest(w, req)

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

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, "Token expired!", httpError.Message)
	assert.Equal(t, "/auth/refresh", httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}

func TestHandleRefreshRequestInvalidTokenInvalidSignature(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
		return
	}

	err = k.Load(confmap.Provider(map[string]interface{}{
		"http.auth.jwt.accessTokenSecret":  "123456",
		"http.auth.jwt.refreshTokenSecret": "abcdef",
	}, "."), nil)
	if err != nil {
		t.Error(err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "excubitor-backend",
		"sub": "testuser",
		"exp": time.Now().Add(-4 * time.Hour).Unix(),
	})
	signedToken, err := token.SignedString([]byte("someOtherKey"))
	if err != nil {
		t.Error(err)
		return
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.RemoteAddr = "SampleAddress"
	req.Header.Set("Authorization", "Bearer "+signedToken)
	w := httptest.NewRecorder()

	handleRefreshRequest(w, req)

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

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, "Invalid token!", httpError.Message)
	assert.Equal(t, "/auth/refresh", httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}

func TestHandleRefreshRequest(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
		return
	}

	err = k.Load(confmap.Provider(map[string]interface{}{
		"http.auth.jwt.accessTokenSecret":  "123456",
		"http.auth.jwt.refreshTokenSecret": "abcdef",
	}, "."), nil)
	if err != nil {
		t.Error(err)
		return
	}

	token, err := signRefreshToken(jwt.MapClaims{
		"iss": "excubitor-backend",
		"sub": "testuser",
		"exp": time.Now().Add(4 * time.Hour).Unix(),
	})
	if err != nil {
		t.Error(err)
		return
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.RemoteAddr = "SampleAddress"
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handleRefreshRequest(w, req)

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

	responseObject := &refreshResponse{}

	err = json.Unmarshal(body, responseObject)
	if err != nil {
		t.Error(err)
		return
	}

	parsedToken, err := jwt.Parse(responseObject.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(k.String("http.auth.jwt.accessTokenSecret")), nil
	})
	if err != nil {
		t.Error(err)
		return
	}

	subject, err := parsedToken.Claims.GetSubject()
	if err != nil {
		t.Error(err)
		return
	}

	issuer, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, subject, "testuser")
	assert.Equal(t, issuer, "excubitor-backend")
}

func TestAuthNoHeader(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/someEndpoint", nil)
	req.RemoteAddr = "SampleAddress"
	w := httptest.NewRecorder()

	handler := auth(nil)
	handler.ServeHTTP(w, req)

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

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, "Bearer authentication is needed!", httpError.Message)
	assert.Equal(t, "Bearer", res.Header.Get("WWW-Authenticate"))
	assert.Equal(t, "/someEndpoint", httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}

func TestAuthInvalidHeader(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/someEndpoint", nil)
	req.RemoteAddr = "SampleAddress"
	req.Header.Set("Authorization", "Basic dXNlcm5hbWU6cGFzc3dvcmQ=")
	w := httptest.NewRecorder()

	handler := auth(nil)
	handler.ServeHTTP(w, req)

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

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, "Bearer authentication is needed!", httpError.Message)
	assert.Equal(t, "Bearer", res.Header.Get("WWW-Authenticate"))
	assert.Equal(t, "/someEndpoint", httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}

func TestAuthTokenExpired(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
		return
	}

	err = k.Load(confmap.Provider(map[string]interface{}{
		"http.auth.jwt.accessTokenSecret":  "123456",
		"http.auth.jwt.refreshTokenSecret": "abcdef",
	}, "."), nil)
	if err != nil {
		t.Error(err)
		return
	}

	token, err := signAccessToken(jwt.MapClaims{
		"iss": "excubitor-backend",
		"sub": "testuser",
		"exp": time.Now().Add(-30 * time.Minute).Unix(),
	})
	if err != nil {
		t.Error(err)
		return
	}

	req := httptest.NewRequest(http.MethodPost, "/someEndpoint", nil)
	req.RemoteAddr = "SampleAddress"
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler := auth(nil)
	handler.ServeHTTP(w, req)

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

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, "Token expired!", httpError.Message)
	assert.Equal(t, "/someEndpoint", httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}

func TestAuthInvalidSignature(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
		return
	}

	err = k.Load(confmap.Provider(map[string]interface{}{
		"http.auth.jwt.accessTokenSecret":  "123456",
		"http.auth.jwt.refreshTokenSecret": "abcdef",
	}, "."), nil)
	if err != nil {
		t.Error(err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "excubitor-backend",
		"sub": "testuser",
		"exp": time.Now().Add(-4 * time.Hour).Unix(),
	})
	signedToken, err := token.SignedString([]byte("someOtherKey"))
	if err != nil {
		t.Error(err)
		return
	}

	req := httptest.NewRequest(http.MethodPost, "/someEndpoint", nil)
	req.RemoteAddr = "SampleAddress"
	req.Header.Set("Authorization", "Bearer "+signedToken)
	w := httptest.NewRecorder()

	handler := auth(nil)
	handler.ServeHTTP(w, req)

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

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, "Invalid token!", httpError.Message)
	assert.Equal(t, "/someEndpoint", httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}

func TestAuth(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
		return
	}

	err = k.Load(confmap.Provider(map[string]interface{}{
		"http.auth.jwt.accessTokenSecret":  "123456",
		"http.auth.jwt.refreshTokenSecret": "abcdef",
	}, "."), nil)
	if err != nil {
		t.Error(err)
		return
	}

	token, err := signAccessToken(jwt.MapClaims{
		"iss": "excubitor-backend",
		"sub": "testuser",
		"exp": time.Now().Add(30 * time.Minute).Unix(),
	})
	if err != nil {
		t.Error(err)
		return
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.RemoteAddr = "SampleAddress"
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler := auth(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte{})
		if err != nil {
			t.Error(err)
			return
		}
	}))
	handler.ServeHTTP(w, req)

	res := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}
