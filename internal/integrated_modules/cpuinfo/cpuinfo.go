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

func Tick() {
	logger = logging.GetLogger()

	cpus, err := readCPUInfo()
	if err != nil {
		logger.Error(fmt.Sprintf("Could not gather cpu information! Reason: %s", err))
	}

	jsonOutput, err := json.Marshal(cpus)
	if err != nil {
		logger.Error(fmt.Sprintf("Couldn't encode cpu information! Reason: %s", err))
	}

	broker := ctx.GetContext().GetBroker()
	broker.Publish("CPU.CpuInfo", string(jsonOutput))
}

func readCPUInfo() ([]cpu, error) {
	file, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}

	paragraphs := regexp.MustCompile(`\n\s*\n`).Split(string(file), -1)

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
