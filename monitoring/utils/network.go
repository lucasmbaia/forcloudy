package utils

import (
  "forcloudy/monitoring/metrics"
  "fmt"
  "io/ioutil"
  "strconv"
  "strings"
  "time"
  "sync"
)

type Utilization struct {
  sync.RWMutex
  pids	map[string][]byte
}

func NewNetwork() *Utilization {
  return &Utilization{pids: make(map[string][]byte)}
}

func (u *Utilization) SetPid(pid string) error {
  var err error

  u.Lock()
  if _, ok := u.pids[pid]; !ok {
    if u.pids[pid], err = ioutil.ReadFile(fmt.Sprintf("/proc/%s/net/dev", pid)); err != nil {
      u.Unlock()
      return err
    }
  }
  u.Unlock()

  return nil
}

func (u *Utilization) DelPid(pid string) {
  u.Lock()
  if _, ok := u.pids[pid]; ok {
    delete(u.pids, pid)
  }
  u.Unlock()
}

func (u *Utilization) NetworkUtilization(pid string, update int) ([]metrics.Networks, error) {
  var (
    previusNet	  []string
    currentNet	  []string
    output	  []byte
    networks	  []metrics.Networks
    previusInfos  []string
    currentInfos  []string
    err		  error
    ok		  bool
  )

  u.Lock()
  if output, ok = u.pids[pid]; ok {
    previusNet = strings.Split(string(u.pids[pid]), "\n")
  } else {
    if output, err = ioutil.ReadFile(fmt.Sprintf("/proc/%s/net/dev", pid)); err != nil {
      u.Unlock()
      return networks, err
    }
    previusNet = strings.Split(string(output), "\n")
    time.Sleep(time.Duration(update) * time.Second)
  }
  u.Unlock()

  if output, err = ioutil.ReadFile(fmt.Sprintf("/proc/%s/net/dev", pid)); err != nil {
    return networks, err
  }
  currentNet = strings.Split(string(output), "\n")

  u.Lock()
  u.pids[pid] = output
  u.Unlock()

  previusNet = previusNet[2 : len(previusNet)-1]
  currentNet = currentNet[2 : len(currentNet)-1]

  for index, _ := range previusNet {
    previusInfos = removeSpace(strings.Split(previusNet[index], " "))
    currentInfos = removeSpace(strings.Split(currentNet[index], " "))

    networks = append(networks, metrics.Networks{
      Interface: strings.Replace(previusInfos[0], ":", "", -1),
      Receive: metrics.Infos{
	Bytes:   convertAndCalc(previusInfos[1], currentInfos[1]),
	Packets: convertAndCalc(previusInfos[2], currentInfos[2]),
	Errors:  convertAndCalc(previusInfos[3], currentInfos[3]),
	Drop:    convertAndCalc(previusInfos[4], currentInfos[4]),
      },
      Trasmit: metrics.Infos{
	Bytes:   convertAndCalc(previusInfos[9], currentInfos[9]),
	Packets: convertAndCalc(previusInfos[10], currentInfos[10]),
	Errors:  convertAndCalc(previusInfos[11], currentInfos[11]),
	Drop:    convertAndCalc(previusInfos[12], currentInfos[12]),
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
