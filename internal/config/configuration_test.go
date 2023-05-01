package config

import (
	"github.com/knadh/koanf/providers/confmap"
	"github.com/stretchr/testify/assert"
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
