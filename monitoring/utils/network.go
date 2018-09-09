package utils

import (
	"fmt"
	"forcloudy/monitoring/metrics"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Utilization struct {
	sync.RWMutex
	pids map[string][]byte
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
		previusNet   []string
		currentNet   []string
		output       []byte
		networks     []metrics.Networks
		previusInfos []string
		currentInfos []string
		err          error
		ok           bool
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

type Network struct {
	sync.RWMutex
}

type Sniffer struct {
	Device      string
	Snaplen     int32
	Promiscuous bool
	Timeout     int
}

func NewSniffer() *Network {
	return &Network{}
}

func (n *Network) SnifferInterface(s Sniffer) error {
	var (
		handle   *pcap.Handle
		err      error
		packet   *gopacket.PacketSource
		ethLayer layers.Ethernet
		ipLayer  layers.IPv4
		tcpLayer layers.TCP
		parser   *gopacket.DecodingLayerParser
		size     uint16
	)

	if handle, err = pcap.OpenLive(s.Device, s.Snaplen, s.Promiscuous, time.Duration(s.Timeout)*time.Second); err != nil {
		return err
	}
	defer handle.Close()

	packet = gopacket.NewPacketSource(handle, handle.LinkType())
	size = 0

	for p := range packet.Packets() {
		parser = gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &ethLayer, &ipLayer, &tcpLayer)
		decoded := []gopacket.LayerType{}

		if err = parser.DecodeLayers(p.Data(), &decoded); err != nil {
			fmt.Println(err)
			continue
		}

		for _, lt := range decoded {
			switch lt {
			case layers.LayerTypeIPv4:
				fmt.Println("IPV4", ipLayer.SrcIP, ipLayer.DstIP)
				fmt.Printf("%+v", ipLayer)
				fmt.Println("SIZE: ", len(ipLayer.Payload))
				if ipLayer.SrcIP.String() == "192.168.204.128" {
					size += ipLayer.Length * uint16(ipLayer.IHL)
					//size += ipLayer.Length
				}
			case layers.LayerTypeTCP:
				fmt.Println("TCP", tcpLayer.SrcPort, tcpLayer.DstPort)
			}
		}

		fmt.Println(size)
		fmt.Println(p)
		/*fmt.Println(ethLayer)
		fmt.Println(ipLayer)
		fmt.Println(tcpLayer)
		fmt.Println(p)*/
	}

	return nil
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
