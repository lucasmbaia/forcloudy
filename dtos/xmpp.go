package dtos

type Container struct {
	Customer    string              `json:",omitempty"`
	Application string              `json:",omitempty"`
	Name        string              `json:",omitempty"`
	Image       string              `json:",omitempty"`
	Ports       map[string][]string `json:",omitempty"`
}
