package http_server

import (
	"github.com/rs/cors"
)

func getCORSHandler() *cors.Cors {

	allowedOrigins := k.Strings("http.cors.allowed_origins")
	allowedMethods := k.Strings("http.cors.allowed_methods")
	allowedHeaders := k.Strings("http.cors.allowed_headers")

	return cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   allowedHeaders,
		AllowCredentials: false,
		Debug:            false,
	})
}
