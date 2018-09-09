package utils

import (
	"fmt"
	"testing"
)

func TestNetworkUtilization(t *testing.T) {
	net := NewNetwork()
	fmt.Println(net.NetworkUtilization("1615", 1))
}

func TestSnifferInterface(t *testing.T) {
	net := NewSniffer()
	s := Sniffer{
		Device:  "ens33",
		Snaplen: 65535,
		Timeout: 30,
	}

	net.SnifferInterface(s)
}
