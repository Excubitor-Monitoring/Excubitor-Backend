package cpu

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type coreStat struct {
	name    string // Name of the cpu core
	User    int64  // Time spent with normal processes running in userspace
	Nice    int64  // Time spent with nice processes running in userspace
	System  int64  // Time spent with processes running in kernel space
	Idle    int64  // Idle time
	Iowait  int64  // Time spent waiting for I/O
	Irq     int64  // Time spent serving hardware interrupts
	Softirq int64  // Time spent serving software interrupts
	Steal   int64  // Time "stolen" by another operating system running in a virutal environment
	Guest   int64  // Time spent running a guest os under control of the kernel
}

type cpuUsage struct {
	Usage float64 `json:"usage"`
}

func calculateCPUUsage() (map[string]cpuUsage, error) {
	returnMap := make(map[string]cpuUsage)

	firstReading, err := os.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	time.Sleep(1000 * time.Millisecond)

	secondReading, err := os.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	firstStats, err := parseStat(string(firstReading))
	if err != nil {
		return nil, err
	}

	secondStats, err := parseStat(string(secondReading))
	if err != nil {
		return nil, err
	}

	for i := range firstStats {
		first := firstStats[i]
		second := secondStats[i]

		firstSum := first.User + first.Nice + first.System + first.Idle + first.Iowait + first.Irq + first.Softirq + first.Steal + first.Guest
		secondSum := second.User + second.Nice + second.System + second.Idle + second.Iowait + second.Irq + second.Softirq + second.Steal + second.Guest

		diff := secondSum - firstSum

		spentIdle := second.Idle - first.Idle
		spentWorking := diff - spentIdle

		usage := float64(100) * float64(spentWorking) / float64(diff)
		returnMap[first.name] = cpuUsage{Usage: usage}
	}

	return returnMap, nil
}

func parseStat(stat string) ([]coreStat, error) {
	statLines := strings.Split(stat, "\n")

	var cores []coreStat

	for _, line := range statLines {
		if strings.HasPrefix(line, "cpu") {
			core, err := parseCoreStat(line)
			if err != nil {
				return nil, err
			}

			cores = append(cores, *core)
		}
	}

	return cores, nil
}

func parseCoreStat(line string) (*coreStat, error) {
	split := strings.Fields(line)

	userValue, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing user value: %w", err)
	}

	niceValue, err := strconv.ParseInt(split[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing nice value: %w", err)
	}

	systemValue, err := strconv.ParseInt(split[3], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing system value: %w", err)
	}

	idleValue, err := strconv.ParseInt(split[4], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing idle value: %w", err)
	}

	iowaitValue, err := strconv.ParseInt(split[5], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing iowait value: %w", err)
	}

	irqValue, err := strconv.ParseInt(split[6], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing irq value: %w", err)
	}

	softirqValue, err := strconv.ParseInt(split[7], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing softirq value: %w", err)
	}

	stealValue, err := strconv.ParseInt(split[8], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing steal value: %w", err)
	}

	guestValue, err := strconv.ParseInt(split[9], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing guest value: %w", err)
	}

	return &coreStat{
		name:    split[0],
		User:    userValue,
		Nice:    niceValue,
		System:  systemValue,
		Idle:    idleValue,
		Iowait:  iowaitValue,
		Irq:     irqValue,
		Softirq: softirqValue,
		Steal:   stealValue,
		Guest:   guestValue,
	}, nil
}
