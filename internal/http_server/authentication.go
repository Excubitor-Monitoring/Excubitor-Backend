package http_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/http_server/helper"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pam"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"net/http"
	"strings"
	"time"
)

type Credentials interface {
	Authenticate() bool
}

type authRequest struct {
	Method      string                 `json:"method"`
	Credentials map[string]interface{} `json:"credentials"`
}

type authResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type refreshResponse struct {
	AccessToken string `json:"access_token"`
}

func handleAuthRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		helper.ReturnError(w, r, http.StatusMethodNotAllowed, "Method is not allowed!")
		return
	}

	if r.Body != nil {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Debug(fmt.Sprintf("Couldn't read message body of auth request from %s", r.RemoteAddr))
			helper.ReturnError(w, r, http.StatusBadRequest, "Can't read message body!")
			return
		}

		request := &authRequest{}
		err = json.Unmarshal(bytes, request)
		if err != nil {
			logger.Debug(fmt.Sprintf("Couldn't decode message body of auth request from %s", r.RemoteAddr))
			helper.ReturnError(w, r, http.StatusBadRequest, "Can't decode message body!")
			return
		}

		switch request.Method {
		case "PAM":
			if request.Credentials == nil {
				logger.Debug(fmt.Sprintf("Could not read credentials in auth request from %s!", r.RemoteAddr))
				helper.ReturnError(w, r, http.StatusBadRequest, "Credentials not specified!")
				return
			}

			if request.Credentials["username"] == nil {
				logger.Debug(fmt.Sprintf("Could not read username in pam auth request from %s!", r.RemoteAddr))
				helper.ReturnError(w, r, http.StatusBadRequest, "Username not specified!")
				return
			}

			if request.Credentials["password"] == nil {
				logger.Debug(fmt.Sprintf("Could not read password in pam auth request from %s!", r.RemoteAddr))
				helper.ReturnError(w, r, http.StatusBadRequest, "Password not specified!")
				return
			}

			username := request.Credentials["username"].(string)
			password := request.Credentials["password"].(string)

			pamCredentials := pam.PAMPasswordCredentials{Username: username, Password: password}

			if pamCredentials.Authenticate() {
				accessTokenClaims := jwt.MapClaims{
					"iss": "excubitor-backend",
					"sub": username,
					"exp": time.Now().Add(30 * time.Minute).Unix(),
				}

				accessToken, err := signAccessToken(accessTokenClaims)
				if err != nil {
					logger.Error(fmt.Sprintf("Couldn't sign access token for %s! Reason: %s", r.RemoteAddr, err))
					helper.ReturnError(w, r, http.StatusInternalServerError, "Internal Server Error!")
					return
				}

				refreshTokenClaims := jwt.MapClaims{
					"iss": "excubitor-backend",
					"sub": username,
					"exp": time.Now().Add(4 * time.Hour).Unix(),
				}

				refreshToken, err := signRefreshToken(refreshTokenClaims)
				if err != nil {
					logger.Error(fmt.Sprintf("Couldn't sign refresh token for %s! Reason: %s", r.RemoteAddr, err))
					helper.ReturnError(w, r, http.StatusInternalServerError, "Internal Server Error!")
					return
				}

				tokens := &authResponse{
					accessToken,
					refreshToken,
				}

				jsonResponse, err := json.Marshal(tokens)
				if err != nil {
					logger.Error(fmt.Sprintf("Couldn't assemble json response for auth request from %s.", r.RemoteAddr))
					helper.ReturnError(w, r, http.StatusInternalServerError, "Internal Server Error!")
					return
				}

				w.WriteHeader(http.StatusOK)
				_, err = w.Write(jsonResponse)
				if err != nil {
					return
				}
			} else {
				helper.ReturnError(w, r, http.StatusUnauthorized, "Invalid username or password!")
				return
			}
		default:
			helper.ReturnError(w, r, http.StatusBadRequest, "Unsupported authentication method: "+request.Method)
			return
		}
	}
}

func handleRefreshRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		helper.ReturnError(w, r, http.StatusMethodNotAllowed, "Method is not allowed!")
		return
	}

	authorization := r.Header.Get("Authorization")

	if !strings.HasPrefix(authorization, "Bearer ") {
		w.Header().Set("WWW-Authenticate", "Bearer")
		helper.ReturnError(w, r, http.StatusUnauthorized, "Bearer authentication is needed!")
		return
	}

	token := strings.Split(authorization, "Bearer ")[1]

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(k.String("http.auth.jwt.refresh_token_secret")), nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer("excubitor-backend"))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			logger.Debug(fmt.Sprintf("Attempt to refresh access token with expired token from %s!", r.RemoteAddr))
			helper.ReturnError(w, r, http.StatusUnauthorized, "Token expired!")
			return
		} else if errors.Is(err, jwt.ErrSignatureInvalid) {
			logger.Warn(fmt.Sprintf("Attempt to authenticate with invalid signature from %s!", r.RemoteAddr))
		} else {
			logger.Debug(fmt.Sprintf("Attempt to authenticate with invalid token from %s! Reason: %s", r.RemoteAddr, err))
		}

		helper.ReturnError(w, r, http.StatusUnauthorized, "Invalid token!")
		return
	}

	username, err := jwtToken.Claims.GetSubject()
	if err != nil {
		logger.Warn(fmt.Sprintf("Couldn't read subject claim of refresh token from %s! Reason: %s", r.RemoteAddr, err))
		helper.ReturnError(w, r, http.StatusBadRequest, "Token has no subject!")
		return
	}

	accessTokenClaims := jwt.MapClaims{
		"iss": "excubitor-backend",
		"sub": username,
		"exp": time.Now().Add(30 * time.Minute).Unix(),
	}

	accessToken, err := signAccessToken(accessTokenClaims)
	if err != nil {
		logger.Error(fmt.Sprintf("Couldn't sign access token for %s! Reason: %s", r.RemoteAddr, err))
		helper.ReturnError(w, r, http.StatusInternalServerError, "Internal Server Error!")
		return
	}

	jsonResponse, err := json.Marshal(refreshResponse{accessToken})
	if err != nil {
		logger.Error(fmt.Sprintf("Couldn't encode access token for %s! Reason: %s", r.RequestURI, err))
		helper.ReturnError(w, r, http.StatusInternalServerError, "Internal Server Error!")
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		logger.Error(fmt.Sprintf("Couldn't send access token to %s! Reason: %s", r.RemoteAddr, err))
		return
	}
}

func bearerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if !strings.HasPrefix(authorization, "Bearer ") {
			w.Header().Set("WWW-Authenticate", "Bearer")
			helper.ReturnError(w, r, http.StatusUnauthorized, "Bearer authentication is needed!")
			return
		}

		token := strings.Split(authorization, "Bearer ")[1]

		if checkToken(w, r, token) {
			next.ServeHTTP(w, r)
		}

	})
}

func queryAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")

		if token == "" {
			logger.Debug(fmt.Sprintf("Attempt to authenticate with invalid token format from %s!", r.RemoteAddr))
			helper.ReturnError(w, r, http.StatusBadRequest, "Token of invalid format!")
			return
		}

		if checkToken(w, r, token) {
			next.ServeHTTP(w, r)
		}

	})
}

func checkToken(w http.ResponseWriter, r *http.Request, token string) bool {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(k.String("http.auth.jwt.access_token_secret")), nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer("excubitor-backend"))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			logger.Debug(fmt.Sprintf("Attempt to authenticate with expired token from %s!", r.RemoteAddr))
			helper.ReturnError(w, r, http.StatusUnauthorized, "Token expired!")
			return false
		} else if errors.Is(err, jwt.ErrSignatureInvalid) {
			logger.Warn(fmt.Sprintf("Attempt to authenticate with invalid signature from %s!", r.RemoteAddr))
		} else {
			logger.Debug(fmt.Sprintf("Attempt to authenticate with invalid token from %s! Reason: %s", r.RemoteAddr, err))
		}

		helper.ReturnError(w, r, http.StatusUnauthorized, "Invalid token!")
		return false
	}

	user, err := jwtToken.Claims.GetSubject()
	if err != nil {
		logger.Warn(fmt.Sprintf("Couldn't read token subject from %s!", user))
	}

	logger.Trace(fmt.Sprintf("User %s authenticated successfully using JWT token!", user))

	return true
}
