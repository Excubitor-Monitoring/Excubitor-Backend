package http_server

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func signAccessToken(claims jwt.MapClaims) (string, error) {
	return signToken(claims, []byte(viper.GetString("http.auth.jwt.accessTokenSecret")))
}

func signRefreshToken(claims jwt.MapClaims) (string, error) {
	return signToken(claims, []byte(viper.GetString("http.auth.jwt.refreshTokenSecret")))
}

func signToken(claims jwt.MapClaims, key interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedString, nil
}
