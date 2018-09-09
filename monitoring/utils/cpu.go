package utils

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/opencontainers/runc/libcontainer/system"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	//FILE_SYSTEM_CPU_USAGE    = "/sys/fs/cgroup/cpuacct/cpuacct.usage"
	FILE_SYSTEM_CPU_USAGE           = "/proc/stat"
	FILE_CONTAINER_CPU_USAGE        = "/sys/fs/cgroup/cpuacct/docker/{container}/cpuacct.usage"
	FILE_CONTAINER_CPU_USAGE_PERCPU = "/sys/fs/cgroup/cpuacct/docker/{container}/cpuacct.usage_percpu"
	nanoSecondsPerSecond            = 1e9
)

func getCpuUsage(name string) (uint64, error) {
	var (
		content []byte
		err     error
	)

	if content, err = ioutil.ReadFile(name); err != nil {
		return 0, err
	}

	return strconv.ParseUint(strings.TrimSuffix(string(content), "\n"), 10, 64)
}

func getCpuUsagePerCPU(name string) ([]uint64, error) {
	var (
		content []byte
		err     error
		cpus    []uint64
		usage   uint64
	)

	if content, err = ioutil.ReadFile(name); err != nil {
		return cpus, err
	}

	for _, cpu := range strings.Split(strings.TrimSuffix(string(content), "\n"), " ") {
		if cpu != "" {
			if usage, err = strconv.ParseUint(cpu, 10, 64); err != nil {
				return cpus, err
			}

			cpus = append(cpus, usage)
		}
	}

	return cpus, nil
}

func systemUsage() (uint64, error) {
	var (
		file                *os.File
		err                 error
		buf                 *bufio.Reader
		usage               uint64
		line                string
		parts               []string
		clockTicksPerSecond uint64
	)

	buf = bufio.NewReaderSize(nil, 128)
	clockTicksPerSecond = uint64(system.GetClockTicks())

	if file, err = os.Open(FILE_SYSTEM_CPU_USAGE); err != nil {
		return usage, err
	}

	defer func() {
		buf.Reset(nil)
		file.Close()
	}()

	buf.Reset(file)

	for err == nil {
		if line, err = buf.ReadString('\n'); err != nil {
			break
		}

		parts = strings.Fields(line)
		switch parts[0] {
		case "cpu":
			if len(parts) < 8 {
				return 0, fmt.Errorf("invalid number of cpu fields")
			}

			var totalClockTicks uint64
			for _, i := range parts[1:8] {
				v, err := strconv.ParseUint(i, 10, 64)
				if err != nil {
					return 0, fmt.Errorf("Unable to convert value %s to int: %s", i, err)
				}
				totalClockTicks += v
			}
			return (totalClockTicks * nanoSecondsPerSecond) /
				clockTicksPerSecond, nil
		}
	}

	return usage, errors.New("invalid stat format. Error trying to parse the '/proc/stat' fil")
}

func CpuUsageContainerUnix(container string, interval int) (float64, error) {
	var (
		previusSystem     uint64
		previusContainer  uint64
		currentSystem     uint64
		currenteContainer uint64
		usagePerCPU       []uint64
		err               error
		percentage        float64
	)

	if previusSystem, err = systemUsage(); err != nil {
		return percentage, err
	}

	if previusContainer, err = getCpuUsage(strings.Replace(FILE_CONTAINER_CPU_USAGE, "{container}", container, -1)); err != nil {
		return percentage, err
	}

	time.Sleep(time.Duration(interval) * time.Second)

	if currentSystem, err = systemUsage(); err != nil {
		return percentage, err
	}

	if currenteContainer, err = getCpuUsage(strings.Replace(FILE_CONTAINER_CPU_USAGE, "{container}", container, -1)); err != nil {
		return percentage, err
	}

	if usagePerCPU, err = getCpuUsagePerCPU(strings.Replace(FILE_CONTAINER_CPU_USAGE_PERCPU, "{container}", container, -1)); err != nil {
		return percentage, err
	}

	percentage = float64(float64((currenteContainer-previusContainer))/float64((currentSystem-previusSystem))) * float64(len(usagePerCPU)) * 100.0

	return percentage, nil
}
