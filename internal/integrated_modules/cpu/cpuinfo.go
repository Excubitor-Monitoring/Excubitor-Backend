package cpu

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type cpu struct {
	Id         uint    `json:"id"`
	CoreId     uint    `json:"core_id"`
	SocketId   uint    `json:"socket_id"`
	ModelName  string  `json:"model_name"`
	ClockSpeed float64 `json:"clock_speed"`
	CacheSize  uint    `json:"cache_size"`
	Flags      string  `json:"flags"`
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

// readCPU parses a single thread/core from a cpuinfo file into a cpu struct.
func readCPU(paragraph string) (*cpu, error) {
	id, err := getUInt("processor", paragraph)
	if err != nil {
		return nil, fmt.Errorf("could not parse processor id: %w", err)
	}

	coreId, err := getUInt("core id", paragraph)
	if err != nil {
		return nil, fmt.Errorf("could not parse core id: %w", err)
	}

	socketId, err := getUInt("physical id", paragraph)
	if err != nil {
		return nil, fmt.Errorf("could not parse physical id: %w", err)
	}

	modelName, err := getString("model name", paragraph)
	if err != nil {
		return nil, fmt.Errorf("could not parse model name: %w", err)
	}

	clockSpeed, err := getFloat64("cpu MHz", paragraph)
	if err != nil {
		return nil, fmt.Errorf("could not parse clock speed: %w", err)
	}

	cache, err := getUInt("cache size", paragraph)
	if err != nil {
		return nil, fmt.Errorf("could not parse cache size: %w", err)
	}

	flags, err := getString("flags", paragraph)
	if err != nil {
		return nil, fmt.Errorf("could not parse flags: %w", err)
	}

	return &cpu{
		Id:         id,
		CoreId:     coreId,
		SocketId:   socketId,
		ModelName:  modelName,
		ClockSpeed: clockSpeed,
		CacheSize:  cache,
		Flags:      flags,
	}, nil
}

// HELPER METHODS

func getUInt(name string, paragraph string) (uint, error) {
	regex := regexp.MustCompile(name + `\s+:\s+(\d+)`)
	matches := regex.FindStringSubmatch(paragraph)
	value, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return 9999, err
	}

	return uint(value), nil
}

func getFloat64(name string, paragraph string) (float64, error) {
	regex := regexp.MustCompile(name + `\s+:\s+([[:digit:]]+.[[:digit:]]+)`)
	matches := regex.FindStringSubmatch(paragraph)
	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 9999, err
	}

	return value, nil
}

func getString(name string, paragraph string) (string, error) {
	regex := regexp.MustCompile(name + `\s+:\s+([[:print:]]+)`)
	matches := regex.FindStringSubmatch(paragraph)

	return matches[1], nil
}
