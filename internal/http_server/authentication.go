package http_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pam"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
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
		ReturnError(w, r, http.StatusMethodNotAllowed, "Method is not allowed!")
		return
	}

	if r.Body != nil {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Debug(fmt.Sprintf("Couldn't read message body of auth request from %s", r.RemoteAddr))
			ReturnError(w, r, http.StatusBadRequest, "Can't read message body!")
			return
		}

		request := &authRequest{}
		err = json.Unmarshal(bytes, request)
		if err != nil {
			logger.Debug(fmt.Sprintf("Couldn't decode message body of auth request from %s", r.RemoteAddr))
			ReturnError(w, r, http.StatusBadRequest, "Can't decode message body!")
			return
		}

		switch request.Method {
		case "PAM":
			username := request.Credentials["username"].(string)
			password := request.Credentials["password"].(string)

			pamCredentials := pam.PAMPasswordCredentials{Username: username, Password: password}

			if pamCredentials.Authenticate() {
				logger.Info("Logged in successfully")

				accessTokenClaims := jwt.MapClaims{
					"iss": "excubitor-backend",
					"sub": username,
					"exp": time.Now().Add(30 * time.Minute).Unix(),
				}

				accessToken, err := signAccessToken(accessTokenClaims)
				if err != nil {
					logger.Error(fmt.Sprintf("Couldn't sign access token for %s! Reason: %s", r.RemoteAddr, err))
					ReturnError(w, r, http.StatusInternalServerError, "Internal Server Error!")
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
					ReturnError(w, r, http.StatusInternalServerError, "Internal Server Error!")
					return
				}

				tokens := &authResponse{
					accessToken,
					refreshToken,
				}

				jsonResponse, err := json.Marshal(tokens)
				if err != nil {
					logger.Error(fmt.Sprintf("Couldn't assemble json response for auth request from %s.", r.RemoteAddr))
					ReturnError(w, r, http.StatusInternalServerError, "Internal Server Error!")
					return
				}

				w.WriteHeader(http.StatusOK)
				_, err = w.Write(jsonResponse)
				if err != nil {
					return
				}
			} else {
				logger.Info("Login attempt was unsuccessful")
				ReturnError(w, r, http.StatusUnauthorized, "Invalid username or password!")
				return
			}
		default:
			ReturnError(w, r, http.StatusBadRequest, "Unsupported authentication method: "+request.Method)
			return
		}
	}
}

func handleRefreshRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorization := r.Header.Get("Authorization")

	if !strings.HasPrefix(authorization, "Bearer ") {
		w.Header().Set("WWW-Authenticate", "Bearer")
		ReturnError(w, r, http.StatusUnauthorized, "Bearer authentication is needed!")
		return
	}

	token := strings.Split(authorization, "Bearer ")[1]

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("http.auth.jwt.refreshTokenSecret")), nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer("excubitor-backend"))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			logger.Debug(fmt.Sprintf("Attempt to refresh access token with expired token from %s!", r.RemoteAddr))
			ReturnError(w, r, http.StatusUnauthorized, "Token expired!")
			return
		} else if errors.Is(err, jwt.ErrSignatureInvalid) {
			logger.Warn(fmt.Sprintf("Attempt to authenticate with invalid signature from %s!", r.RemoteAddr))
		} else {
			logger.Debug(fmt.Sprintf("Attempt to authenticate with invalid token from %s! Reason: %s", r.RemoteAddr, err))
		}

		ReturnError(w, r, http.StatusUnauthorized, "Invalid token!")
		return
	}

	username, err := jwtToken.Claims.GetSubject()
	if err != nil {
		logger.Warn(fmt.Sprintf("Couldn't read subject claim of refresh token from %s! Reason: %s", r.RemoteAddr, err))
		ReturnError(w, r, http.StatusBadRequest, "Token has no subject!")
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
		ReturnError(w, r, http.StatusInternalServerError, "Internal Server Error!")
		return
	}

	jsonResponse, err := json.Marshal(refreshResponse{accessToken})
	if err != nil {
		logger.Error(fmt.Sprintf("Couldn't encode access token for %s! Reason: %s", r.RequestURI, err))
		ReturnError(w, r, http.StatusInternalServerError, "Internal Server Error!")
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		logger.Error(fmt.Sprintf("Couldn't send access token to %s! Reason: %s", r.RemoteAddr, err))
		return
	}
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if !strings.HasPrefix(authorization, "Bearer ") {
			w.Header().Set("WWW-Authenticate", "Bearer")
			ReturnError(w, r, http.StatusUnauthorized, "Bearer authentication is needed!")
			return
		}

		token := strings.Split(authorization, "Bearer ")[1]

		jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("http.auth.jwt.accessTokenSecret")), nil
		}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer("excubitor-backend"))

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				logger.Debug(fmt.Sprintf("Attempt to authenticate with expired token from %s!", r.RemoteAddr))
				ReturnError(w, r, http.StatusUnauthorized, "Token expired!")
				return
			} else if errors.Is(err, jwt.ErrSignatureInvalid) {
				logger.Warn(fmt.Sprintf("Attempt to authenticate with invalid signature from %s!", r.RemoteAddr))
			} else {
				logger.Debug(fmt.Sprintf("Attempt to authenticate with invalid token from %s! Reason: %s", r.RemoteAddr, err))
			}

			ReturnError(w, r, http.StatusUnauthorized, "Invalid token!")
			return
		}

		user, err := jwtToken.Claims.GetSubject()
		if err != nil {
			logger.Warn(fmt.Sprintf("Couldn't read token subject from %s!", user))
		}

		logger.Trace(fmt.Sprintf("User %s authenticated successfully using JWT token!", user))

		next.ServeHTTP(w, r)
	})
}
