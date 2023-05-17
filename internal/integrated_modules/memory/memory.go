package memory

import (
	"encoding/json"
	"fmt"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"os"
)

var logger logging.Logger

func Tick() {
	logger = logging.GetLogger()

	memInfoFile, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		logger.Error(fmt.Sprintf("Could not read file '/proc/meminfo'. Reason: %s", err))
		return
	}

	entries, err := parseMemInfo(string(memInfoFile))
	if err != nil {
		logger.Error(fmt.Sprintf("Could not parse file '/proc/meminfo'. Reason: %s", err))
		return
	}

	meminfo := getMemInfoFromMap(entries)
	swapinfo := getSwapInfoFromMap(entries)

	memInfoJSON, err := json.Marshal(meminfo)
	if err != nil {
		logger.Error(fmt.Sprintf("Couldn't encode memory information! Reason: %s", err))
		return
	}

	swapInfoJSON, err := json.Marshal(swapinfo)
	if err != nil {
		logger.Error(fmt.Sprintf("Couldn't encode swap information! Reason: %s", err))
	}

	broker := ctx.GetContext().GetBroker()

	broker.Publish("Memory.MemInfo", string(memInfoJSON))
	broker.Publish("Memory.SwapInfo", string(swapInfoJSON))
}
