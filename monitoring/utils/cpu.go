package utils

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	FILE_SYSTEM_CPU_USAGE    = "/sys/fs/cgroup/cpuacct/cpuacct.usage"
	FILE_CONTAINER_CPU_USAGE = "/sys/fs/cgroup/cpuacct/docker/{container}/cpuacct.usage"
)

func getCpuUsage(name string) (int64, error) {
	var (
		content []byte
		err     error
	)

	if content, err = ioutil.ReadFile(name); err != nil {
		return 0, err
	}

	return strconv.ParseInt(strings.TrimSuffix(string(content), "\n"), 10, 64)
}

func CpuUsageContainerUnix(container string) (float64, error) {
	var (
		previusSystem     int64
		previusContainer  int64
		currentSystem     int64
		currenteContainer int64
		err               error
		percentage        float64
	)

	if previusSystem, err = getCpuUsage(FILE_SYSTEM_CPU_USAGE); err != nil {
		return percentage, err
	}

	if previusContainer, err = getCpuUsage(strings.Replace(FILE_CONTAINER_CPU_USAGE, "{container}", container, -1)); err != nil {
		return percentage, err
	}

	time.Sleep(1 * time.Second)

	if currentSystem, err = getCpuUsage(FILE_SYSTEM_CPU_USAGE); err != nil {
		return percentage, err
	}

	if currenteContainer, err = getCpuUsage(strings.Replace(FILE_CONTAINER_CPU_USAGE, "{container}", container, -1)); err != nil {
		return percentage, err
	}

	//percentage = float64(float64((currenteContainer-previusContainer))/float64((currentSystem-previusSystem))) * float64(2) * 100.0
	fmt.Println(float64(float64((currenteContainer - previusContainer)) / float64((currentSystem - previusSystem))))
	percentage = float64(float64((currenteContainer-previusContainer))/float64((currentSystem-previusSystem))) * 100.0

	return percentage, nil
}
