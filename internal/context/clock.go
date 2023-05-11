package ctx

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/config"
	"sync"
	"time"
)

var clockOnce sync.Once

func startClock() {
	clockOnce.Do(func() {
		clockString := config.GetConfig().String("data.module_clock")
		clock, err := time.ParseDuration(clockString)
		if err != nil {
			context.logger.Fatal("Could not parse module clock from configuration. Check your configuration values!")
			panic(err)
		}

		go func() {
			for {
				modules := GetContext().GetModules()
				for _, module := range modules {
					module.tickFunction()
				}
				time.Sleep(clock)
			}
		}()
	})
}
