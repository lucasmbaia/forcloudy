package dtos

type ApplicationEtcd struct {
	Protocol        map[string]string `json:"protocol,omitempty"`
	Image           string            `json:"image,omitempty"`
	PortsDST        []string          `json:"portsDst,omitempty"`
	Cpus            string            `json:"cpus,omitempty"`
	Dns             string            `json:"dns,omitempty"`
	Memory          string            `json:"memory,omitempty"`
	TotalContainers int               `json:"totalContainers,omitempty"`
}
