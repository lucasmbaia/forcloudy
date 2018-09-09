package utils

import (
	"fmt"
	"testing"
)

func TestSystemCPUUsage(t *testing.T) {
	fmt.Println(systemUsage())
}

func TestCpuUsageContainerUnix(t *testing.T) {
	for {
		fmt.Println(CpuUsageContainerUnix("06bc381a5d0344d0c90d1b774d077a847dcb66c604b6bb5de2e2e5cce2983219", 1))
	}
}

func TestGetCpuUsagePerCPU(t *testing.T) {
	fmt.Println(getCpuUsagePerCPU("/sys/fs/cgroup/cpuacct/docker/06bc381a5d0344d0c90d1b774d077a847dcb66c604b6bb5de2e2e5cce2983219/cpuacct.usage_percpu"))
}
