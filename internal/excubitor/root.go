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
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/plugins"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/pubsub"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/pkg/shared/modules"
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
		modules.NewModule(
			"main",
			modules.NewVersion(0, 0, 1),
			[]modules.Component{},
			func() {},
		),
	)

	context.RegisterModule(
		modules.NewModule(
			"cpu",
			modules.NewVersion(0, 0, 1),
			[]modules.Component{
				{
					TabName: "CPU Information",
					JSFile:  "static/internal/cpu/info.js",
					Tag:     "cpu-info",
				},
				{
					TabName: "CPU Clock History",
					JSFile:  "static/internal/cpu/clock-history.js",
					Tag:     "cpu-clock-history",
				},
				{
					TabName: "CPU Usage",
					JSFile:  "static/internal/cpu/usage.js",
					Tag:     "cpu-usage",
				},
				{
					TabName: "CPU Usage History",
					JSFile:  "static/internal/cpu/usage-history.js",
					Tag:     "cpu-usage-history",
				},
			},
			cpu.Tick,
		),
	)

	context.RegisterModule(
		modules.NewModule(
			"memory",
			modules.NewVersion(0, 0, 1),
			[]modules.Component{
				{
					TabName: "RAM Usage",
					JSFile:  "static/internal/memory/ram.js",
					Tag:     "ram-usage",
				},
				{
					TabName: "RAM Usage History",
					JSFile:  "static/internal/memory/ram-history.js",
					Tag:     "ram-usage-history",
				},
				{
					TabName: "Swap Usage",
					JSFile:  "static/internal/memory/swap.js",
					Tag:     "swap-usage",
				},
				{
					TabName: "Swap Usage History",
					JSFile:  "static/internal/memory/swap-history.js",
					Tag:     "swap-usage-history",
				},
			},
			memory.Tick,
		),
	)

	if err := plugins.LoadPlugins(); err != nil {
		return err
	}
	if err := plugins.InitPlugins(); err != nil {
		return err
	}

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
