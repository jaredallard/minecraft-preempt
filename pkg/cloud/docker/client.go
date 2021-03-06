package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"github.com/jaredallard/minecraft-preempt/pkg/cloud"
	"github.com/pkg/errors"
)

var (
	ErrNotStopped = errors.New("not stopped")
)

type Client struct {
	d dockerclient.APIClient
}

func NewClient() (*Client, error) {
	c, err := dockerclient.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create docker client")
	}

	return &Client{c}, nil
}

func (c *Client) Start(ctx context.Context, containerID string) error {
	cont, err := c.d.ContainerInspect(ctx, containerID)
	if err != nil {
		return err
	}

	if cont.State.Status != "exited" {
		return ErrNotStopped
	}

	return c.d.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{})
}

func (c *Client) Status(ctx context.Context, containerID string) (cloud.ProviderStatus, error) {
	cont, err := c.d.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", err
	}

	switch cont.State.Status {
	case "exited", "dead":
		return cloud.StatusStopped, nil
	case "removing":
		return cloud.StatusStopping, nil
	case "running":
		return cloud.StatusRunning, nil
	case "created":
		return cloud.StatusStarting, nil
	}

	return cloud.StatusUnknown, nil
}
