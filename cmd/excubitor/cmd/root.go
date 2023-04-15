package cmd

import (
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/http_server"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pubsub"
	"time"
)

func Execute() error {
	var err error

	logger := logging.GetLogger()

	if err != nil {
		return err
	}

	context := ctx.GetContext()
	context.RegisterModule(ctx.NewModule("main"))
	context.RegisterBroker(pubsub.NewBroker())

	logger.Debug("Starting HTTP Server!")

	go func() {
		for {
			time.Sleep(5 * time.Second)
			ctx.GetContext().GetBroker().Publish("Some.Monitor", "Test Message!")
		}
	}()

	err = http_server.Start()

	if err != nil {
		return err
	}

	return nil
}
