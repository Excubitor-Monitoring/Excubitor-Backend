package cpuinfo

import (
	"encoding/json"
	"fmt"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"os"
	"regexp"
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

// readCPUInfoFile reads the contents of /proc/cpuinfo and returns them in a byte slice.
func readCPUInfoFile() ([]byte, error) {
	file, err := os.ReadFile("/proc/cpuinfo")
	return file, err
}

// readCPUInfo can parse multiple threads from a cpuinfo file.
func readCPUInfo(cpuInfo string) ([]cpu, error) {
	paragraphs := regexp.MustCompile(`\n\s*\n`).Split(cpuInfo, -1)

	cpus := make([]cpu, len(paragraphs)-1)

	for i, p := range paragraphs {
		if len(p) != 0 {
			cpu, err := readCPU(p)
			if err != nil {
				return nil, err
			}
			cpus[i] = *cpu
		}
	}

	return cpus, nil
}
