package ctx

import "time"

func startClock() {
	go func() {
		for {
			modules := GetContext().GetModules()
			for _, module := range modules {
				module.tickFunction()
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}()
}
