package cmd

import (
	"errors"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/http_server"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/spf13/viper"
	"io/fs"
	"os"
)

func Execute() error {

	var err error

	logger, err := logging.GetLogger()

	if err != nil {
		return err
	}

	ctx.GetContext().RegisterModule(ctx.NewModule("main"))

	logger.Trace("Hallo Welt!")
	logger.Debug("Hallo Welt!")
	logger.Info("Hallo Welt!")
	logger.Warn("Hallo Welt!")
	logger.Error("Hallo Welt!")
	logger.Fatal("Hallo Welt!")

	err = http_server.Start()

	if err != nil {
		return err
	}

	return nil
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

func init() {
	err := initConfig()
	if err != nil {
		panic(err)
	}
}
