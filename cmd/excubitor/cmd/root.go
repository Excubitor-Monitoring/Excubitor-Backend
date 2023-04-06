package cmd

import (
	"errors"
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/spf13/viper"
	"io/fs"
	"os"
	"strings"
)

func Execute() error {

	var err error

	if err = initConfig(); err != nil {
		return err
	}

	var logger logging.Logger
	loggingMethod := viper.GetString("logging.method")

	switch strings.ToUpper(loggingMethod) {
	case "CONSOLE":
		logger, err = logging.GetConsoleLoggerInstance()
		break
	case "FILE":
		logger, err = logging.GetFileLoggerInstance()
		break
	case "HYBRID":
		logger, err = logging.GetMultiLoggerInstance()
		break
	default:
		return fmt.Errorf("unknown logging method %s", loggingMethod)
	}

	if err != nil {
		return err
	}

	logger.Trace("Hallo Welt!")
	logger.Debug("Hallo Welt!")
	logger.Info("Hallo Welt!")
	logger.Warn("Hallo Welt!")
	logger.Error("Hallo Welt!")
	logger.Fatal("Hallo Welt!")

	return nil
}

func initConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.yml")
	viper.SetConfigType("yaml")

	viper.SetDefault("logging.log_level", "INFO")
	viper.SetDefault("logging.method", "CONSOLE")

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
