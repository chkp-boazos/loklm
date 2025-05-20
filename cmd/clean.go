package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/boazos/loklm/pkg/cfgmgr"
	dockerclient "github.com/boazos/loklm/pkg/dockerClient"
	"github.com/boazos/loklm/pkg/fs"
	"github.com/boazos/loklm/pkg/tasks"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	keepStateDir bool
	cleanImages  bool
)

func CleanStateDirectory(dirPath string) tasks.Task {
	return func() error {
		remover := fs.GetDefaultDirRemover()
		if !keepStateDir {
			if err := remover.Remove(dirPath); err != nil {
				return err
			}
		}
		return nil
	}
}

func CleanContainer(name string, client *client.Client, ctx context.Context) tasks.Task {
	return func() error {
		return client.ContainerRemove(ctx, name, container.RemoveOptions{Force: true})
	}
}

func CleanNetwork(name string, client *client.Client, ctx context.Context) tasks.Task {
	return func() error {
		return client.NetworkRemove(ctx, name)
	}
}

func CleanImage(name string, client *client.Client, ctx context.Context) tasks.Task {
	return func() error {
		_, err := client.ImageRemove(ctx, name, image.RemoveOptions{Force: true})
		return err
	}
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean environment",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := cfgmgr.LoadToml(configPath)
		if err != nil {
			fmt.Println("Failed loading configuration: ", err)
			os.Exit(1)
		}

		client, err := dockerclient.GetDefaultClient()
		if err != nil {
			fmt.Println("Failed setting docker client:", err)
			os.Exit(1)
		}

		ctx := context.Background()

		cleanTasks := []tasks.Task{
			tasks.WithResults(fmt.Sprintf("Container %s was removed", cfg.Notebooks.Name), fmt.Sprintf("Could not remove container %s", cfg.Notebooks.Name))(
				tasks.WithSpinner(fmt.Sprintf(":sparkles: Removing container %s", cfg.Notebooks.Name))(
					CleanContainer(cfg.Notebooks.Name, client, ctx),
				),
			),

			tasks.WithResults(fmt.Sprintf("Container %s was removed", cfg.Llm.Name), fmt.Sprintf("Could not remove container %s", cfg.Llm.Name))(
				tasks.WithSpinner(fmt.Sprintf(":sparkles: Removing container %s", cfg.Llm.Name))(
					CleanContainer(cfg.Llm.Name, client, ctx),
				),
			),

			tasks.WithResults(fmt.Sprintf("Network %s was removed", cfg.General.Network), fmt.Sprintf("Could not remove network %s", cfg.General.Network))(
				tasks.WithSpinner(fmt.Sprintf(":cookie: Removing Network %s", cfg.General.Network))(
					CleanNetwork(cfg.General.Network, client, ctx),
				),
			),
		}

		if cfg.VectorDB != nil {
			cleanTasks = append(
				[]tasks.Task{
					tasks.WithResults(fmt.Sprintf("Container %s was removed", cfg.VectorDB.Name), fmt.Sprintf("Could not remove container %s", cfg.VectorDB.Name))(
						tasks.WithSpinner(fmt.Sprintf(":sparkles: Removing container %s", cfg.VectorDB.Name))(
							CleanContainer(cfg.VectorDB.Name, client, ctx),
						),
					),
				},
				cleanTasks...,
			)
		}

		if cleanImages {
			cleanTasks = append(cleanTasks,
				tasks.WithResults(fmt.Sprintf("Image %s was removed", cfg.Notebooks.Image), fmt.Sprintf("Could not remove image %s", cfg.Llm.Image))(
					tasks.WithSpinner(fmt.Sprintf(":cookie: Removing image %s", cfg.Notebooks.Image))(
						CleanImage(cfg.Notebooks.Image, client, ctx),
					),
				),

				tasks.WithResults(fmt.Sprintf("Image %s was removed", cfg.Llm.Image), fmt.Sprintf("Could not remove image %s", cfg.Llm.Image))(
					tasks.WithSpinner(fmt.Sprintf(":cookie: Removing image %s", cfg.Llm.Image))(
						CleanImage(cfg.Llm.Image, client, ctx),
					),
				),
			)

			if cfg.VectorDB != nil {
				cleanTasks = append(
					cleanTasks,
					tasks.WithResults(fmt.Sprintf("Image %s was removed", cfg.VectorDB.Image), fmt.Sprintf("Could not remove image %s", cfg.VectorDB.Image))(
						tasks.WithSpinner(fmt.Sprintf(":cookie: Removing image %s", cfg.VectorDB.Image))(
							CleanImage(cfg.VectorDB.Image, client, ctx),
						),
					),
				)
			}
		}

		if !keepStateDir {
			cleanTasks = append(cleanTasks, tasks.WithResults("State directory successfully was removed", fmt.Sprintf("Could not remove state directory %s", cfg.General.StateDir))(
				tasks.WithSpinner(":rocket: Deleting state directory")(
					CleanStateDirectory(cfg.General.StateDir),
				),
			),
			)
		}
		errorCode := 0
		for _, cleanTask := range cleanTasks {
			if err := cleanTask(); err != nil {
				errorCode = 1
			}
		}
		os.Exit(errorCode)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.PersistentFlags().BoolVar(&keepStateDir, "keep-state", false, "Do not clean notebooks and models")
	cleanCmd.PersistentFlags().BoolVar(&cleanImages, "delete-images", false, "Clean Docker images")
}
