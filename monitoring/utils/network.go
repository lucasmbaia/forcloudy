package utils

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type Networks struct {
	Interface string `json:",omitempty"`
	Receive   infos  `json:",omitempty"`
	Trasmit   infos  `json:",omitempty"`
}

type infos struct {
	Bytes   int64 `json:",omitempty"`
	Packets int64 `json:",omitempty"`
	Errors  int64 `json:",omitempty"`
	Drop    int64 `json:",omitempty"`
}

func NetworkUtilization(pid string, update int) ([]Networks, error) {
	var (
		previusNet   []string
		currentNet   []string
		output       []byte
		networks     []Networks
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

		networks = append(networks, Networks{
			Interface: strings.Replace(previusInfos[0], ":", "", -1),
			Receive: infos{
				Bytes:   convertAndCalc(previusInfos[1], currentInfos[1]),
				Packets: convertAndCalc(previusInfos[2], currentInfos[2]),
				Errors:  convertAndCalc(previusInfos[3], currentInfos[3]),
				Drop:    convertAndCalc(previusInfos[4], currentInfos[4]),
			},
			Trasmit: infos{
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
