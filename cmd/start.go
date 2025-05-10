package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/boaz0/cortexpoint/pkg/cfgmgr"
	dockerclient "github.com/boaz0/cortexpoint/pkg/dockerClient"
	"github.com/boaz0/cortexpoint/pkg/tasks"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"
)

func startContainerTask(containerCfg cfgmgr.ExposedContainer, generalCfg cfgmgr.General, client *client.Client, ctx context.Context) tasks.Task {
	return func() error {
		dirBinds := []string{}
		if containerCfg.Dir != "" {
			dirBinds = []string{fmt.Sprintf("%s/%s:%s", generalCfg.StateDir, containerCfg.Name, containerCfg.Dir)}
		}

		resp, err := client.ContainerCreate(ctx, &container.Config{
			Image:    containerCfg.Image,
			Hostname: containerCfg.Hostname,
		}, &container.HostConfig{
			Binds: dirBinds,
			PortBindings: nat.PortMap{
				containerCfg.ExposedPort: []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: fmt.Sprintf("%d", containerCfg.Port),
					},
				},
			},
		}, &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				generalCfg.Network: {},
			},
		}, nil, containerCfg.Name)

		if doesContainerExist(err) {
			return nil
		}

		if err != nil {
			return err
		}

		return client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	}
}

func wrapTask(task tasks.Task, name string) tasks.Task {
	return tasks.WithResults(
		fmt.Sprintf("Container %s is running successfully", name),
		fmt.Sprintf("Container %s failed to run successfully", name),
	)(
		tasks.WithSpinner(fmt.Sprintf(":rocket: Running container %s", name))(
			task,
		),
	)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the containers",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := cfgmgr.LoadToml(configPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		client, err := dockerclient.GetDefaultClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ctx := context.Background()

		startTasks := []tasks.Task{
			wrapTask(startContainerTask(cfg.Notebooks, cfg.General, client, ctx), cfg.Notebooks.Name),
			wrapTask(startContainerTask(cfg.Llm, cfg.General, client, ctx), cfg.Llm.Name),
		}

		errorCode := 0
		for _, startTask := range startTasks {
			if err := startTask(); err != nil {
				errorCode = 1
			}
		}

		os.Exit(errorCode)
	},
}

func doesContainerExist(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "You have to remove (or rename) that container to be able to reuse that name.")
}

func init() {
	rootCmd.AddCommand(startCmd)
}
