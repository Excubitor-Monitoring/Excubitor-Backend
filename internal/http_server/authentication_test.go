package http_server

import (
	"bytes"
	"encoding/json"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/golang-jwt/jwt/v5"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

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

func TestHandleAuthRequestPAMNilCredentials(t *testing.T) {
	type testParams struct {
		description          string
		method               string
		credentials          map[string]interface{}
		expectedErrorMessage string
		expectedStatusCode   int
	}

	for _, params := range []testParams{
		{
			description:          "Nil credentials",
			method:               "PAM",
			credentials:          nil,
			expectedErrorMessage: "Credentials not specified!",
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			description: "Nil username",
			method:      "PAM",
			credentials: map[string]interface{}{
				"username": nil,
				"password": "SomePassword",
			},
			expectedErrorMessage: "Username not specified!",
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			description: "Nil password",
			method:      "PAM",
			credentials: map[string]interface{}{
				"username": "SomeUser",
				"password": nil,
			},
			expectedErrorMessage: "Password not specified!",
			expectedStatusCode:   http.StatusBadRequest,
		},
	} {
		t.Run(params.description, func(t *testing.T) {
			var err error

			logger, err = logging.GetConsoleLoggerInstance()
			if err != nil {
				t.Error(err)
				return
			}

			payload, err := json.Marshal(authRequest{Method: params.method, Credentials: params.credentials})

			req := httptest.NewRequest(http.MethodPost, "/auth",
				bytes.NewReader(payload))
			req.RemoteAddr = "SampleAddress"
			w := httptest.NewRecorder()

			handleAuthRequest(w, req)

			res := w.Result()
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					t.Error(err)
					return
				}
			}(res.Body)

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Error(err)
				return
			}

			httpError := parseHTTPError(body)

			assert.Equal(t, params.expectedStatusCode, res.StatusCode)
			assert.Equal(t, params.expectedErrorMessage, httpError.Message)
			assert.Equal(t, httpError.Path, "/auth")
			assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
		})
	}

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
		"http.auth.jwt.access_token_secret":  "123456",
		"http.auth.jwt.refresh_token_secret": "abcdef",
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
		"http.auth.jwt.access_token_secret":  "123456",
		"http.auth.jwt.refresh_token_secret": "abcdef",
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
		"http.auth.jwt.access_token_secret":  "123456",
		"http.auth.jwt.refresh_token_secret": "abcdef",
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
		return []byte(k.String("http.auth.jwt.access_token_secret")), nil
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

func TestBearerAuthNoHeader(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/someEndpoint", nil)
	req.RemoteAddr = "SampleAddress"
	w := httptest.NewRecorder()

	handler := bearerAuth(nil)
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

func TestBearerAuthInvalidHeader(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/someEndpoint", nil)
	req.RemoteAddr = "SampleAddress"
	req.Header.Set("Authorization", "Basic dXNlcm5hbWU6cGFzc3dvcmQ=")
	w := httptest.NewRecorder()

	handler := bearerAuth(nil)
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

func TestBearerAuthTokenExpired(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
		return
	}

	err = k.Load(confmap.Provider(map[string]interface{}{
		"http.auth.jwt.access_token_secret":  "123456",
		"http.auth.jwt.refresh_token_secret": "abcdef",
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

	handler := bearerAuth(nil)
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

func TestBearerAuthInvalidSignature(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
		return
	}

	err = k.Load(confmap.Provider(map[string]interface{}{
		"http.auth.jwt.access_token_secret":  "123456",
		"http.auth.jwt.refresh_token_secret": "abcdef",
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

	handler := bearerAuth(nil)
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

func TestBearerAuth(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
		return
	}

	err = k.Load(confmap.Provider(map[string]interface{}{
		"http.auth.jwt.access_token_secret":  "123456",
		"http.auth.jwt.refresh_token_secret": "abcdef",
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

	req := httptest.NewRequest(http.MethodPost, "/someEndpoint", nil)
	req.RemoteAddr = "SampleAddress"
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler := bearerAuth(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
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

func TestQueryAuthNoToken(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/someEndpoint?token=", nil)
	req.RemoteAddr = "SampleAddress"
	w := httptest.NewRecorder()

	handler := queryAuth(nil)
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

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "Token of invalid format!", httpError.Message)
	assert.Equal(t, "/someEndpoint?token=", httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}

func TestQueryAuthTokenExpired(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
		return
	}

	err = k.Load(confmap.Provider(map[string]interface{}{
		"http.auth.jwt.access_token_secret":  "123456",
		"http.auth.jwt.refresh_token_secret": "abcdef",
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

	req := httptest.NewRequest(http.MethodPost, "/someEndpoint?token="+token, nil)
	req.RemoteAddr = "SampleAddress"
	w := httptest.NewRecorder()

	handler := queryAuth(nil)
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
	assert.Equal(t, "/someEndpoint?token="+token, httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}

func TestQueryAuthInvalidSignature(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
		return
	}

	err = k.Load(confmap.Provider(map[string]interface{}{
		"http.auth.jwt.access_token_secret":  "123456",
		"http.auth.jwt.refresh_token_secret": "abcdef",
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

	req := httptest.NewRequest(http.MethodPost, "/someEndpoint?token="+signedToken, nil)
	req.RemoteAddr = "SampleAddress"
	w := httptest.NewRecorder()

	handler := queryAuth(nil)
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
	assert.Equal(t, "/someEndpoint?token="+signedToken, httpError.Path)
	assert.True(t, time.Since(httpError.Timestamp) < time.Since(time.Now().Add(-time.Second)) && time.Until(httpError.Timestamp) < 0)
}

func TestQueryAuth(t *testing.T) {
	var err error

	logger, err = logging.GetConsoleLoggerInstance()
	if err != nil {
		t.Error(err)
		return
	}

	err = k.Load(confmap.Provider(map[string]interface{}{
		"http.auth.jwt.access_token_secret":  "123456",
		"http.auth.jwt.refresh_token_secret": "abcdef",
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

	req := httptest.NewRequest(http.MethodPost, "/someEndpoint?token="+token, nil)
	req.RemoteAddr = "SampleAddress"
	w := httptest.NewRecorder()

	handler := queryAuth(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
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
