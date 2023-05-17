package memory

import (
	"regexp"
	"strconv"
	"strings"
)

type meminfo struct {
	MemTotal     int64 `json:"mem_total"`
	MemFree      int64 `json:"mem_free"`
	MemAvailable int64 `json:"mem_available"`
}

type swapInfo struct {
	SwapTotal int64 `json:"swap_total"`
	SwapFree  int64 `json:"swap_free"`
}

func parseMemInfo(input string) (map[string]int64, error) {
	entries := make(map[string]int64)

	lines := strings.Split(input, "\n")

	regex := regexp.MustCompile(":\\s+")

	for _, line := range lines[:len(lines)-1] {

		splitted := regex.Split(line, -1)

		var value string
		if strings.HasSuffix(splitted[1], "kB") {
			value = splitted[1][:len(splitted[1])-3]
		} else {
			value = splitted[1]
		}

		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}

		entries[splitted[0]] = parsed
	}

	return entries, nil
}

func getMemInfoFromMap(entries map[string]int64) meminfo {
	return meminfo{
		MemTotal:     entries["MemTotal"],
		MemFree:      entries["MemFree"],
		MemAvailable: entries["MemAvailable"],
	}
}

func getSwapInfoFromMap(entries map[string]int64) swapInfo {
	return swapInfo{
		SwapTotal: entries["SwapTotal"],
		SwapFree:  entries["SwapFree"],
	}
}
