package metrics

type Metrics struct {
	Hostname  string      `json:",omitempty"`
	Customers []Customers `json:",omitempty"`
}

type Customers struct {
	Name         string         `json:",omitempty"`
	Applications []Applications `json:",omitempty"`
}

type Applications struct {
	Name       string       `json:",omitempty"`
	Containers []Containers `json:",omitempty"`
}

type Containers struct {
	ID          string      `json:",omitempty"`
	Name        string      `json:",omitempty"`
	Cpu         Cpu         `json:",omitempty"`
	Memory      Memory      `json:",omitempty"`
	Disk        Disk        `json:",omitempty"`
	Networks    []Networks  `json:",omitempty"`
	Connections Connections `json:",omitempty"`
}

type Cpu struct {
	TotalUsage float64 `json:",omitempty"`
	//PerCPU     []int64 `json:",omitempty"`
}

type Memory struct {
	TotalMemory int64 `json:",omitempty"`
	TotalUsage  int64 `json:",omitempty"`
}

type Disk struct {
}

type Networks struct {
	Interface string `json:",omitempty"`
	Receive   Infos  `json:",omitempty"`
	Trasmit   Infos  `json:",omitempty"`
}

type Infos struct {
	Bytes   int64 `json:",omitempty"`
	Packets int64 `json:",omitempty"`
	Errors  int64 `json:",omitempty"`
	Drop    int64 `json:",omitempty"`
}

type Connections struct {
	Address string `json:",omitempty"`
	Status  string `json:",omitempty"`
	Start   string `json:",omitempty"`
	End     string `json:",omitempty"`
}
