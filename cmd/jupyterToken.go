package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/boazos/loklm/pkg/cfgmgr"
	dockerclient "github.com/boazos/loklm/pkg/dockerClient"
	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
)

var jupyterTokenCmd = &cobra.Command{
	Use:   "jupyterToken",
	Short: "Get list of all jupyter notebooks and thier tokens",
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

		execOpts := container.ExecOptions{
			Tty:          true,
			AttachStderr: true,
			AttachStdout: true,
			Cmd:          []string{"jupyter", "notebook", "list"},
		}

		execCreate, err := client.ContainerExecCreate(ctx, cfg.Notebooks.Name, execOpts)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		resp, err := client.ContainerExecAttach(ctx, execCreate.ID, container.ExecStartOptions{
			Tty: true,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		io.Copy(os.Stdout, resp.Reader)
	},
}

func init() {
	rootCmd.AddCommand(jupyterTokenCmd)
}
