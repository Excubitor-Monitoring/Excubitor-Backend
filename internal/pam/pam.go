package pam

import (
	"errors"
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/msteinert/pam"
)

type PAMPasswordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (cred PAMPasswordCredentials) Authenticate() bool {
	return passwordAuthentication(cred.Username, cred.Password)
}

var ErrUnrecognizedMessageStyle = errors.New("unrecognized message style")

func passwordAuthentication(username string, password string) bool {
	logger := logging.GetLogger()

	t, err := pam.StartFunc("excubitor", username, func(style pam.Style, msg string) (string, error) {
		switch style {
		case pam.PromptEchoOff:
			fallthrough
		case pam.PromptEchoOn:
			return password, nil
		case pam.ErrorMsg:
			return "", fmt.Errorf("authentication error: %s", msg)
		case pam.TextInfo:
			return "", nil
		default:
			return "", ErrUnrecognizedMessageStyle
		}
	})

	if err != nil {
		logger.Warn(fmt.Sprintf("Authentication of user %s failed on start: %s", username, err))
		return false
	}

	err = t.Authenticate(0)
	if err != nil {
		logger.Warn(fmt.Sprintf("Authentication of user %s failed on authentication: %s", username, err))
		return false
	}

	logger.Trace(fmt.Sprintf("Password-based authentication of user %s was successful!", username))
	return true
}
