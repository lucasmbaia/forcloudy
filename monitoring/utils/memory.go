package utils

import (
  "forcloudy/monitoring/metrics"
  "io/ioutil"
  "strconv"
  "strings"
  "fmt"
)

func MemoryUtilization(id string) (metrics.Memory, error) {
  var (
    output  []byte
    memory  metrics.Memory
    err	    error
  )

  if output, err = ioutil.ReadFile(fmt.Sprintf("/sys/fs/cgroup/memory/docker/%s/memory.usage_in_bytes", id)); err != nil {
    return memory, err
  }

  if memory.TotalUsage, err = strconv.ParseInt(strings.TrimSuffix(string(output), "\n"), 10, 64); err != nil {
    return memory, err
  }

  if output, err = ioutil.ReadFile(fmt.Sprintf("/sys/fs/cgroup/memory/docker/%s/memory.limit_in_bytes", id)); err != nil {
    return memory, err
  }

  if memory.TotalMemory, err = strconv.ParseInt(strings.TrimSuffix(string(output), "\n"), 10, 64); err != nil {
    return memory, err
  }

  return memory, err
}
