package http_server

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func signAccessToken(claims jwt.MapClaims) (string, error) {
	return signToken(claims, []byte(config.GetConfig().String("http.auth.jwt.accessTokenSecret")))
}

func signRefreshToken(claims jwt.MapClaims) (string, error) {
	return signToken(claims, []byte(config.GetConfig().String("http.auth.jwt.refreshTokenSecret")))
}

func signToken(claims jwt.MapClaims, key interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedString, nil
}
