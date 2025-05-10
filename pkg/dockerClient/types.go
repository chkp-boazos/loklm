package dockerclient

import (
	"context"

	"github.com/docker/docker/api/types/network"
)

type NetworkCreator interface {
	NetworkCreate(context context.Context, name string, options network.CreateOptions) (network.CreateResponse, error)
}
