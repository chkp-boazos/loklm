package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/boaz0/cortexpoint/pkg/cfgmgr"
	dockerclient "github.com/boaz0/cortexpoint/pkg/dockerClient"
	"github.com/boaz0/cortexpoint/pkg/fs"
	"github.com/boaz0/cortexpoint/pkg/tasks"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/pcj/mobyprogress"
	"github.com/spf13/cobra"
)

func pullImageTask(name string, client *client.Client) tasks.Task {
	return func() error {
		events, err := client.ImagePull(context.Background(), name, image.PullOptions{})
		if err != nil {
			return err
		}
		defer events.Close()
		streamPullEvents(events)
		return nil
	}
}

func createNetworkTask(name string, client *client.Client) tasks.Task {
	return func() error {
		return createNetwork(client, name)
	}
}

func createDirectoryTask(name string, dirCreator fs.DirCreator) tasks.Task {
	return func() error {
		return dirCreator.CreateDirectory(name, 0755)
	}
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set resources to run the local environment",
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

		setupTasks := []tasks.Task{
			tasks.WithResults(
				fmt.Sprintf("Network %s was successfully created", cfg.General.Network),
				fmt.Sprintf("Network %s was not created successfully", cfg.General.Network),
			)(
				tasks.WithSpinner(":rocket: Creating network")(
					createNetworkTask(cfg.General.Network, client),
				),
			),

			tasks.WithResults(
				fmt.Sprintf("State directory %s was created successfully", cfg.General.StateDir),
				fmt.Sprintf("State directory %s was not created successfully", cfg.General.StateDir),
			)(
				tasks.WithSpinner(":hard_disk: Creating state directory")(
					createDirectoryTask(cfg.General.StateDir, fs.GetDefaultDirCreator()),
				),
			),

			tasks.WithResults(
				fmt.Sprintf("Image %s was pulled successfully", cfg.Notebooks.Image),
				fmt.Sprintf("Image %s was not pulled successfully", cfg.Notebooks.Image),
			)(
				pullImageTask(cfg.Notebooks.Image, client),
			),

			tasks.WithResults(
				fmt.Sprintf("Image %s was pulled successfully", cfg.Llm.Image),
				fmt.Sprintf("Image %s was not pulled successfully", cfg.Llm.Image),
			)(
				pullImageTask(cfg.Llm.Image, client),
			),
		}

		if cfg.VectorDB != nil {
			setupTasks = append(
				setupTasks,
				tasks.WithResults(
					fmt.Sprintf("Image %s was pulled successfully", cfg.VectorDB.Image),
					fmt.Sprintf("Image %s was not pulled successfully", cfg.VectorDB.Image),
				)(
					pullImageTask(cfg.VectorDB.Image, client),
				),
			)
		}

		errorCode := 0
		for _, setupTask := range setupTasks {
			err := setupTask()
			if err != nil {
				errorCode = 1
			}
		}
		os.Exit(errorCode)
	},
}

func streamPullEvents(events io.Reader) {
	out := mobyprogress.NewOut(os.Stdout)
	progress := mobyprogress.NewProgressOutput(out)
	decoder := json.NewDecoder(events)
	for {
		var event struct {
			Status   string `json:"status"`
			ID       string `json:"id,omitempty"`
			Progress string `json:"progress,omitempty"`
			Error    string `json:"error,omitempty"`
		}
		if err := decoder.Decode(&event); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error decoding event:", err)
			break
		}

		if event.Error != "" {
			fmt.Fprintf(os.Stderr, "Error: %s\n", event.Error)
			continue
		}

		if event.ID != "" {
			progress.WriteProgress(mobyprogress.Progress{
				ID:      event.ID,
				Action:  event.Status,
				Current: int64(len(event.Progress)),
				Total:   int64(len(event.Progress)),
				Units:   "bytes",
			})
		} else {
			progress.WriteProgress(mobyprogress.Progress{
				ID:      event.Status,
				Action:  event.Status,
				Current: int64(len(event.Progress)),
				Total:   int64(len(event.Progress)),
				Units:   "bytes",
			})
		}
	}
}

func doesNetworkNameExist(err error, name string) bool {
	// I wish I could import libnetwork.NetorkNameError and do errors.Is, but the compiler does not allow me to import this :(
	return strings.Contains(err.Error(), fmt.Sprintf("network with name %s already exists", name))
}

func createNetwork(client dockerclient.NetworkCreator, networkName string) error {
	_, err := client.NetworkCreate(context.Background(), networkName, network.CreateOptions{
		Driver: network.NetworkBridge,
		Scope:  "local",
	})
	if err != nil && !doesNetworkNameExist(err, networkName) {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
