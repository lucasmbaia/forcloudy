package main

import (
  "context"
  "encoding/json"
  "flag"
  "fmt"
  "github.com/lucasmbaia/forcloudy/ha-proxy/template"
  "github.com/lucasmbaia/forcloudy/etcd"
  "log"
  "net"
  "os"
  "os/signal"
  "strings"
  "syscall"
)

var (
  timeout   = flag.Int("timeout", 5, "timeout of connect etcd")
  key       = flag.String("key", "/fc-haproxy/", "Key of watch etcd")
  hosts     = flag.String("host", "http://172.16.95.183:2379", "Host of etcd")
  path      = flag.String("path", "/usr/local/etc/", "Path to conf ha-proxy")
  model     = flag.String("model", "", "Model of template")
  addr      = flag.String("addr", "", "Address of Network Interface")
  lbproxy   = flag.String("lbproxy", "lb-proxy-1", "")
  exclusive = flag.Bool("exclusive", false, "")
  certs     = flag.String("certs", "", "Certs of ssl")

  protocols = map[string]string{
    "app-http":  "http",
    "app-https": "https",
  }
)

type InfosApplication struct {
  Name      string  `json:"name,omitempty"`
  Hosts     []Hosts `json:"hosts,omitempty"`
  Dns       string  `json:"dns,omitempty"`
  Interface string  `json:"-"`
  SSL       string  `json:"-"`
}

type Hosts struct {
  Name       string       `json:"name,omitempty"`
  Dns        string       `json:"dns,omitempty"`
  Protocol   string       `json:"protocol,omitempty"`
  PortSRC    string       `json:"portSRC,omitempty"`
  Containers []Containers `json:"containers,omitempty"`
  Address    []string     `json:"-"`
  Whitelist  string       `json:"-"`
  Minions    []string     `json:"-"`
}

type Containers struct {
  Name    string `json:"-"`
  Address string `json:"address,omitempty"`
  Minion  string `json:"minion,omitempty"`
}

func (ia InfosApplication) AddrToMinion(addr, minion string) string {
  return fmt.Sprintf("%s:%s", minion, strings.Split(addr, ":")[1])
}

func Whitelist(address []string, minions []string) string {
  var addrs string

  for _, v := range address {
    addrs = fmt.Sprintf("%s%s ", addrs, strings.Split(v, ":")[0])
  }

  addrs = fmt.Sprintf("%s%s %s", addrs, strings.Join(minions, " "), *lbproxy)
  return addrs
}

func AddressToMinion(containers []Containers, exclusive bool) []string {
  var addrs []string

  if exclusive {
    var hostname string

    hostname, _ = os.Hostname()

    for _, container := range containers {
      if hostname == container.Minion {
	addrs = append(addrs, fmt.Sprintf("%s:%s", container.Minion, strings.Split(container.Address, ":")[1]))
      }
    }
  } else {
    for _, container := range containers {
      addrs = append(addrs, fmt.Sprintf("%s:%s", container.Minion, strings.Split(container.Address, ":")[1]))
    }
  }

  return addrs
}

func AddressInterface(name string) string {
  var (
    ifaces []net.Interface
    err    error
    addr   string
  )

  if ifaces, err = net.Interfaces(); err != nil {
    log.Fatalf("Error to list interfaces: %s", err.Error())
  }

  for _, iface := range ifaces {
    if iface.Name == name {
      var addrs []net.Addr

      if addrs, err = iface.Addrs(); err != nil {
	log.Fatalf("Error to get infos of interface: %s", err.Error())
      }

      for _, a := range addrs {
	var ip net.IP

	if ip, _, err = net.ParseCIDR(a.String()); err != nil {
	  log.Fatalf("Error to parse CIDR: %s", err.Error())
	}

	addr = ip.String()
	break
      }
    }
  }

  return addr
}

func pending(name string, values []byte, mt string) {
  var (
    ia  InfosApplication
    err error
  )

  if err = json.Unmarshal([]byte(values), &ia); err != nil {
    log.Printf("Error unmarshal: %s", err.Error())
    return
  }

  for key, _ := range ia.Hosts {
    for _, container := range ia.Hosts[key].Containers {
      ia.Hosts[key].Address = append(ia.Hosts[key].Address, container.Address)
      ia.Hosts[key].Minions = append(ia.Hosts[key].Minions, container.Minion)
    }

    ia.Hosts[key].Whitelist = Whitelist(ia.Hosts[key].Address, removeStringDuplicates(ia.Hosts[key].Minions))
    ia.Hosts[key].Address = AddressToMinion(ia.Hosts[key].Containers, *exclusive)
    ia.Hosts[key].Minions = removeStringDuplicates(ia.Hosts[key].Minions)
  }

  if *certs != "" {
    ia.SSL = fmt.Sprintf("ssl %s", strings.Join(strings.Split(*certs, ","), "crt "))
  }

  ia.Name = strings.Replace(name, *key, "", 1)
  ia.Interface = *addr

  if err = template.ConfGenerate(*path, ia.Name, mt, ia); err != nil {
    log.Printf("Error to generate conf: %s", err.Error())
  }

  return
}

func main() {
  var (
    ctx    context.Context
    cancel context.CancelFunc
    err    error
    cli    etcd.Client
    values = make(chan etcd.Response)
    sigs   = make(chan os.Signal, 1)
  )

  flag.Parse()

  if *model == "" {
    log.Fatalf("Reporte the model of template")
  }

  if *addr == "" {
    log.Fatalf("Reporte the address of network usage in haproxy")
  }

  ctx, cancel = context.WithCancel(context.Background())
  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

  go func() {
    <-sigs
    cancel()
  }()

  if cli, err = etcd.NewClient(ctx, etcd.Config{
    Endpoints:	[]string{*hosts},
    Timeout:	int32(*timeout),
  }); err != nil {
    log.Fatalf("Error to connect etcd: %s", err.Error())
  }

  go func() {
    for {
      var mt string

      infos := <-values
      key := strings.Replace(infos.Key, *key, "", 1)

      fmt.Println(infos)

      switch infos.Action {
      case "set":
	if name, ok := protocols[key]; ok {
	  mt, err = template.ModelConf(fmt.Sprintf("%s-%s", *model, name))
	} else {
	  mt, err = template.ModelConf(fmt.Sprintf("%s-%s", *model, "tcpudp"))
	}

	if err != nil {
	  log.Printf("Error to get conf of template: %s", err.Error())
	  break
	}

	pending(infos.Key, []byte(infos.Values), mt)
      case "delete":
	if err = template.RemoveConf(*path, key); err != nil {
	  log.Printf("Error to delete file conf: %s", err.Error())
	}
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

func removeStringDuplicates(elem []string) []string {
  var (
    encountered = map[string]bool{}
    result      []string
  )

  for v := range elem {
    encountered[elem[v]] = true
  }

  for key, _ := range encountered {
    result = append(result, key)
  }

  return result
}
