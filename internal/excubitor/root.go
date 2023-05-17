package excubitor

import (
	"errors"
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/config"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/db"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/http_server"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/integrated_modules/cpu"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/integrated_modules/memory"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pubsub"
)

func Execute() error {
	var err error

	if err := config.InitConfig(); err != nil {
		return err
	}
	if err := logging.InitLogging(); err != nil {
		return err
	}

	if config.GetConfig().Bool("main.print_startup_banner") {
		printBanner()
	}

	logger := logging.GetLogger()

	logger.Debug("Starting startup check...")
	if err := check(); err != nil {
		if !errors.Is(err, SoftfailError) {
			return err
		}
	}

	logger.Debug("Initializing database!")
	if err := db.InitDatabase(); err != nil {
		return err
	}

	logger.Debug("Loading context...")
	context := ctx.GetContext()
	context.RegisterModule(ctx.NewModule("main", func() {
	}))
	context.RegisterModule(ctx.NewModule("cpu", cpu.Tick))
	context.RegisterModule(ctx.NewModule("memory", memory.Tick))
	logger.Debug("Registering broker...")
	context.RegisterBroker(pubsub.NewBroker())

	logger.Debug("Starting HTTP Server!")

	err = http_server.Start()

	if err != nil {
		return err
	}

	return nil
}

func printBanner() {
	banner := `
		 ▄▄▄▄▄▄▄▄▄▄▄  ▄       ▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄         ▄  ▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄ 
		▐░░░░░░░░░░░▌▐░▌     ▐░▌▐░░░░░░░░░░░▌▐░▌       ▐░▌▐░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌
		▐░█▀▀▀▀▀▀▀▀▀  ▐░▌   ▐░▌ ▐░█▀▀▀▀▀▀▀▀▀ ▐░▌       ▐░▌▐░█▀▀▀▀▀▀▀█░▌▀▀▀▀█░█▀▀▀▀  ▀▀▀▀█░█▀▀▀▀ ▐░█▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀█░▌
		▐░▌            ▐░▌ ▐░▌  ▐░▌          ▐░▌       ▐░▌▐░▌       ▐░▌    ▐░▌          ▐░▌     ▐░▌       ▐░▌▐░▌       ▐░▌
		▐░█▄▄▄▄▄▄▄▄▄    ▐░▐░▌   ▐░▌          ▐░▌       ▐░▌▐░█▄▄▄▄▄▄▄█░▌    ▐░▌          ▐░▌     ▐░▌       ▐░▌▐░█▄▄▄▄▄▄▄█░▌
		▐░░░░░░░░░░░▌    ▐░▌    ▐░▌          ▐░▌       ▐░▌▐░░░░░░░░░░▌     ▐░▌          ▐░▌     ▐░▌       ▐░▌▐░░░░░░░░░░░▌
		▐░█▀▀▀▀▀▀▀▀▀    ▐░▌░▌   ▐░▌          ▐░▌       ▐░▌▐░█▀▀▀▀▀▀▀█░▌    ▐░▌          ▐░▌     ▐░▌       ▐░▌▐░█▀▀▀▀█░█▀▀ 
		▐░▌            ▐░▌ ▐░▌  ▐░▌          ▐░▌       ▐░▌▐░▌       ▐░▌    ▐░▌          ▐░▌     ▐░▌       ▐░▌▐░▌     ▐░▌  
		▐░█▄▄▄▄▄▄▄▄▄  ▐░▌   ▐░▌ ▐░█▄▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄█░▌▄▄▄▄█░█▄▄▄▄      ▐░▌     ▐░█▄▄▄▄▄▄▄█░▌▐░▌      ▐░▌ 
		▐░░░░░░░░░░░▌▐░▌     ▐░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░▌▐░░░░░░░░░░░▌     ▐░▌     ▐░░░░░░░░░░░▌▐░▌       ▐░▌
		 ▀▀▀▀▀▀▀▀▀▀▀  ▀       ▀  ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀▀       ▀       ▀▀▀▀▀▀▀▀▀▀▀  ▀         ▀ 
		v0.0.1-alpha


		`

	fmt.Println(banner)
}
