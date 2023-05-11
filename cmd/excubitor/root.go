package excubitor

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/config"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/db"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/http_server"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/integrated_modules/cpuinfo"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pubsub"
)

func Execute() error {
	var err error

	if err := config.InitConfig(); err != nil {
		return err
	}
	if err := initLogging(); err != nil {
		return err
	}

	logger := logging.GetLogger()

	logger.Debug("Initializing database!")
	if err := db.InitDatabase(); err != nil {
		return err
	}

	logger.Debug("Loading context...")
	context := ctx.GetContext()
	context.RegisterModule(ctx.NewModule("main", func() {
	}))
	context.RegisterModule(ctx.NewModule("cpu", cpuinfo.Tick))
	logger.Debug("Registering broker...")
	context.RegisterBroker(pubsub.NewBroker())

	logger.Debug("Starting HTTP Server!")

	err = http_server.Start()

	if err != nil {
		return err
	}

	return nil
}

func initLogging() error {
	method := config.GetConfig().String("logging.method")

	err := logging.SetDefaultLogger(method)
	if err != nil {
		return err
	}

	return nil
}
