package cmd

import (
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/http_server"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/integrated_modules/cpuinfo"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pubsub"
)

func Execute() error {
	var err error

	logger := logging.GetLogger()

	if err != nil {
		return err
	}

	context := ctx.GetContext()

	context.RegisterModule(ctx.NewModule("main", func() {
		logger.Trace("Tick!")
	}))
	context.RegisterModule(ctx.NewModule("cpu", cpuinfo.Tick))

	context.RegisterBroker(pubsub.NewBroker())

	logger.Debug("Starting HTTP Server!")

	err = http_server.Start()

	if err != nil {
		return err
	}

	return nil
}
