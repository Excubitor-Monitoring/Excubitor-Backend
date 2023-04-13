package main

import (
	"errors"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/cmd/excubitor/cmd"
	"github.com/spf13/viper"
	"io/fs"
	"os"
)

func main() {
	if err := initConfig(); err != nil {
		panic(err)
	}

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

func initConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.yml")
	viper.SetConfigType("yaml")

	viper.SetDefault("logging.log_level", "INFO")
	viper.SetDefault("logging.method", "CONSOLE")
	viper.SetDefault("http.port", 8080)
	viper.SetDefault("http.host", "localhost")

	if _, err := os.Stat("config.yml"); errors.Is(err, fs.ErrNotExist) {
		err := viper.WriteConfig()
		if err != nil {
			return err
		}
	}

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}
