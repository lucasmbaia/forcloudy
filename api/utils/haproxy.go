package utils

const (
	KEY_ETCD = "/fc-haproxy"
)

var (
	PROTOCOL_HTTP = []string{"http", "https"}
	PORTS_HTTP    = []string{"80", "443"}
	KEY_HTTP      = map[string]string{"80": "app-http", "443": "app-https"}
)

type Haproxy struct {
	ApplicationsName string
	ContainerName    string
	PortsContainer   map[string][]string
	Protocol         map[string]string
	AddressContainer string
	Dns              string
}

type httpHttps struct {
	ApplicationsName  string
	ContainerName     string
	PortSource        string
	PortsDestionation []string
	AddressContainer  string
	Dns               string
}

func GenerateConf(h Haproxy) error {
	var (
		exists bool
	)

	for src, dst := range h.PortsContainer {
		if _, exists = ExistsStringElement(src, PORTS_HTTP); exists {

		}
	}
}

func httpAndHttps(h httpHttps) {
	var (
		key string
	)

	key = fmt.Sprintf("%s%s", KEY_ETCD, KEY_HTTP[h.PortSource])
}

func ExistsStringElement(f string, s []string) (int, bool) {
	for idx, str := range s {
		if str == f {
			return idx, true
		}
	}

	return 0, false
}
