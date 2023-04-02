package cmd

import "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"

func Execute() error {

	logger, err := logging.GetMultiLoggerInstance()
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
