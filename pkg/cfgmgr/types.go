package cfgmgr

import "github.com/docker/go-connections/nat"

type Configuration struct {
	General   General
	Notebooks ExposedContainer
	Llm       ExposedContainer
	VectorDB  *ExposedContainer
}

type General struct {
	Network  string
	StateDir string
}

type ExposedContainer struct {
	Name        string
	Port        int
	ExposedPort nat.Port
	Hostname    string
	Image       string
	Dir         string
}
