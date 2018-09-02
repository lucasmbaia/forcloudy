package utils

import (
  "forcloudy/monitoring/metrics"
  "fmt"
  "io/ioutil"
  "strconv"
  "strings"
  "time"
)

func NetworkUtilization(pid string, update int) ([]metrics.Networks, error) {
  var (
    previusNet   []string
    currentNet   []string
    output       []byte
    networks     []metrics.Networks
    previusInfos []string
    currentInfos []string
    err          error
  )

  if output, err = ioutil.ReadFile(fmt.Sprintf("/proc/%s/net/dev", pid)); err != nil {
    return networks, err
  }
  previusNet = strings.Split(string(output), "\n")

  time.Sleep(time.Duration(update) * time.Second)

  if output, err = ioutil.ReadFile(fmt.Sprintf("/proc/%s/net/dev", pid)); err != nil {
    return networks, err
  }
  currentNet = strings.Split(string(output), "\n")

  previusNet = previusNet[2 : len(previusNet)-1]
  currentNet = currentNet[2 : len(currentNet)-1]

  for index, _ := range previusNet {
    previusInfos = removeSpace(strings.Split(previusNet[index], " "))
    currentInfos = removeSpace(strings.Split(currentNet[index], " "))

    networks = append(networks, metrics.Networks{
      Interface: strings.Replace(previusInfos[0], ":", "", -1),
      Receive: Infos{
	Bytes:   convertAndCalc(previusInfos[1], currentInfos[1]),
	Packets: convertAndCalc(previusInfos[2], currentInfos[2]),
	Errors:  convertAndCalc(previusInfos[3], currentInfos[3]),
	Drop:    convertAndCalc(previusInfos[4], currentInfos[4]),
      },
      Trasmit: Infos{
	Bytes:   convertAndCalc(previusInfos[10], currentInfos[10]),
	Packets: convertAndCalc(previusInfos[11], currentInfos[11]),
	Errors:  convertAndCalc(previusInfos[12], currentInfos[12]),
	Drop:    convertAndCalc(previusInfos[13], currentInfos[13]),
      },
    })
  }

  return networks, nil
}

func convertAndCalc(previus, current string) int64 {
  var (
    p int64
    c int64
  )

  p, _ = strconv.ParseInt(previus, 10, 64)
  c, _ = strconv.ParseInt(current, 10, 64)

  return c - p
}

func removeSpace(values []string) []string {
  var infos []string

  for _, value := range values {
    if value != "" {
      infos = append(infos, value)
    }
  }

  return infos
}
