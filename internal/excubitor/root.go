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
	context.RegisterModule(
		ctx.NewModule(
			"main",
			ctx.NewVersion(0, 0, 1),
			[]ctx.Component{},
			func() {},
		),
	)

	context.RegisterModule(
		ctx.NewModule(
			"cpu",
			ctx.NewVersion(0, 0, 1),
			[]ctx.Component{
				{
					TabName: "CPU Information",
					JSFile:  "static/internal/CPU-Info/index.js",
					Tag:     "cpu-info",
				},
				{
					TabName: "CPU Clock History",
					JSFile:  "static/internal/CPU-Info/history.js",
					Tag:     "cpu-clock-history",
				},
			},
			cpu.Tick,
		),
	)

	context.RegisterModule(
		ctx.NewModule(
			"memory",
			ctx.NewVersion(0, 0, 1),
			[]ctx.Component{
				{
					TabName: "RAM Usage",
					JSFile:  "static/internal/RAM-Usage/index.js",
					Tag:     "ram-usage",
				},
				{
					TabName: "RAM Usage History",
					JSFile:  "static/internal/RAM-Usage/history.js",
					Tag:     "ram-usage-history",
				},
				{
					TabName: "Swap Usage",
					JSFile:  "static/internal/Swap-Usage/index.js",
					Tag:     "swap-usage",
				},
				{
					TabName: "Swap Usage History",
					JSFile:  "static/internal/Swap-Usage/history.js",
					Tag:     "swap-usage-history",
				},
			},
			memory.Tick,
		),
	)

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
