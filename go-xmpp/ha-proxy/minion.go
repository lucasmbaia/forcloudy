package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"forcloudy/go-xmpp/ha-proxy/template"
	"forcloudy/go-xmpp/ha-proxy/watch"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	timeout = flag.Int("timeout", 5, "timeout of connect etcd")
	key     = flag.String("key", "/fc-haproxy/", "Key of watch etcd")
	hosts   = flag.String("host", "http://192.168.204.128:2379", "Host of etcd")
	path    = flag.String("path", "/etc/haproxy/", "Path to conf ha-proxy")
)

type InfosApplication struct {
	Name  string  `json:"name,omitempty"`
	Hosts []Hosts `json:"hosts,omitempty"`
	Dns   string  `json:"dns,omitempty"`
}

type Hosts struct {
	Protocol   string       `json:"protocol,omitempty"`
	PortSRC    string       `json:"portSRC,omitempty"`
	Containers []Containers `json:"containers,omitempty"`
	Address    []string     `json:"-"`
	Whitelist  string       `json:"-"`
}

type Containers struct {
	Name    string `json:"-"`
	Address string `json:"address,omitempty"`
}

func Whitelist(address []string) string {
	var addrs string

	for _, v := range address {
		addrs = fmt.Sprintf("%s%s ", addrs, strings.Split(v, ":")[0])
	}

	addrs = fmt.Sprintf("%s%s %s %s", addrs, "minion-1", "minion-2", "minion-3")
	return addrs
}

func main() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		err    error
		cli    *watch.Client
		values = make(chan watch.WatchInfos)
		sigs   = make(chan os.Signal, 1)
	)

	ctx, cancel = context.WithCancel(context.Background())

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancel()
	}()

	if cli, err = watch.New([]string{*hosts}, *timeout); err != nil {
		log.Fatalf("Error to connect etcd: %s", err.Error())
	}

	go func() {
		for {
			var ia InfosApplication
			infos := <-values

			if err = json.Unmarshal([]byte(infos.Values), &ia); err != nil {
				log.Printf("Error unmarshal: %s", err.Error())
				continue
			}

			fmt.Println(infos.Key, infos.Values)

			for key, _ := range ia.Hosts {
				for _, container := range ia.Hosts[key].Containers {
					ia.Hosts[key].Address = append(ia.Hosts[key].Address, container.Address)
				}

				fmt.Println(ia.Hosts[key].Address)
				ia.Hosts[key].Whitelist = Whitelist(ia.Hosts[key].Address)
			}

			ia.Name = strings.Replace(infos.Key, *key, "", 1)

			if err = template.ConfGenerate(*path, ia.Name, template.MINION, ia); err != nil {
				log.Printf("Error to generate conf: %s", err.Error())
				continue
			}
		}
	}()

	go func() {
		if err = cli.Watch(*key, values); err != nil {
			log.Fatalf("Watch error: %s", err.Error())
		}
	}()

	<-ctx.Done()
}
