package dockerclient

import (
	"github.com/docker/docker/client"
)

func GetDefaultClient() (*client.Client, error) {
	return client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
}
