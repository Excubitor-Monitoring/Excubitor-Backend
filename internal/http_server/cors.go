package http_server

import (
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func getCORSHandler() *cors.Cors {

	allowedOrigins := viper.GetStringSlice("http.cors.allowed_origins")
	allowedMethods := viper.GetStringSlice("http.cors.allowed_methods")
	allowedHeaders := viper.GetStringSlice("http.cors.allowed_headers")

	return cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   allowedHeaders,
		AllowCredentials: false,
		Debug:            false,
	})
}
