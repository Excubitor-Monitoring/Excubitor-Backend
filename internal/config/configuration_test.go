package config

import (
	"github.com/knadh/koanf/providers/confmap"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCheckConfig(t *testing.T) {
	// HTTP port too low

	err := k.Load(confmap.Provider(map[string]interface{}{
		"http.port":                          0,
		"http.auth.jwt.access_token_secret":  "abcde",
		"http.auth.jwt.refresh_token_secret": "abcde",
	}, "."), nil)
	if err != nil {
		t.Error(err)
		return
	}

	err = checkConfig()
	assert.ErrorIs(t, err, ErrInvalidConfigParameter)
	assert.Equal(t, "invalid config parameter: port needs to be at least 1 and lower than 65536. Is: 0", err.Error())

	// HTTP port too high

	err = k.Load(confmap.Provider(map[string]interface{}{
		"http.port":                          65536,
		"http.auth.jwt.access_token_secret":  "abcde",
		"http.auth.jwt.refresh_token_secret": "abcde",
	}, "."), nil)
	if err != nil {
		t.Error(err)
		return
	}

	err = checkConfig()
	assert.ErrorIs(t, err, ErrInvalidConfigParameter)
	assert.Equal(t, "invalid config parameter: port needs to be at least 1 and lower than 65536. Is: 65536", err.Error())

	// Access Token Secret not set

	err = k.Load(confmap.Provider(map[string]interface{}{
		"http.auth.jwt.access_token_secret": "",
	}, "."), nil)
	if err != nil {
		t.Error(err)
		return
	}

	err = checkConfig()
	assert.ErrorIs(t, err, ErrInvalidConfigParameter)
	assert.Equal(t, "invalid config parameter: access token secret is not set", err.Error())

	// Refresh Token Secret not set

	err = k.Load(confmap.Provider(map[string]interface{}{
		"http.auth.jwt.access_token_secret":  "abcde",
		"http.auth.jwt.refresh_token_secret": "",
	}, "."), nil)
	if err != nil {
		t.Error(err)
		return
	}

	err = checkConfig()
	assert.ErrorIs(t, err, ErrInvalidConfigParameter)
	assert.Equal(t, "invalid config parameter: refresh token secret is not set", err.Error())
}

func TestInitConfigENV(t *testing.T) {
	t.Setenv("EXCUBITOR_LOGGING_LOG-LEVEL", "TRACE")
	t.Setenv("EXCUBITOR_HTTP_AUTH_JWT_ACCESS-TOKEN-SECRET", "abcde")
	t.Setenv("EXCUBITOR_HTTP_AUTH_JWT_REFRESH-TOKEN-SECRET", "abcde")

	movedConfig := false

	if _, err := os.Stat("config.yml"); !os.IsNotExist(err) {
		if err := os.Rename("config.yml", "config.yml.bak"); err != nil {
			t.Error(err)
			return
		}
		movedConfig = true
	}

	_, err := os.Create("config.yml")
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		if err := os.Remove("config.yml"); err != nil {
			t.Error(err)
		}

		if movedConfig {
			if err := os.Rename("config.yml.bak", "config.yml"); err != nil {
				t.Error(err)
				return
			}
		}
	}()

	if err := InitConfig(); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "TRACE", k.String("logging.log_level"))
}
