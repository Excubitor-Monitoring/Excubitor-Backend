package http_server

import (
	"encoding/json"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pam"
	"io"
	"net/http"
)

type Credentials interface {
	Authenticate() bool
}

type authRequest struct {
	Method      string      `json:"method"`
	Credentials Credentials `json:"credentials"`
}

type authResponse struct {
	Token string `json:"token"`
}

func handleAuthRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body != nil {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			ReturnError(w, r, http.StatusBadRequest, "Can't read message body!")
			return
		}

		request := &authRequest{}
		err = json.Unmarshal(bytes, request)
		if err != nil {
			ReturnError(w, r, http.StatusBadRequest, "Can't decode message body!")
		}

		switch request.Method {
		case "PAM":
			if request.Credentials.(pam.PAMPasswordCredentials).Authenticate() {
				logger.Info("Logged in successfully")
				_, err := w.Write(nil)
				if err != nil {
					return
				}
			} else {
				logger.Info("Login attempt was unsuccessful")
			}
		}
	}
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Authenticate based on JWT token provided as bearer
		next.ServeHTTP(w, r)
	})
}
