package cpu

import (
	"encoding/json"
	"fmt"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
)

var logger logging.Logger

// Tick is a function that is called whenever the context wants the module to report its values.
func Tick() {
	logger = logging.GetLogger()

	cpuInfo, err := readCPUInfoFile()
	if err != nil {
		logger.Error(fmt.Sprintf("Could not read /proc/cpuinfo! Reason: %s", err))
		return
	}

	cpus, err := readCPUInfo(string(cpuInfo))
	if err != nil {
		logger.Error(fmt.Sprintf("Could not gather cpu information! Reason: %s", err))
		return
	}

	jsonOutput, err := json.Marshal(cpus)
	if err != nil {
		logger.Error(fmt.Sprintf("Couldn't encode cpu information! Reason: %s", err))
		return
	}

	broker := ctx.GetContext().GetBroker()
	broker.Publish("CPU.CpuInfo", string(jsonOutput))
}
