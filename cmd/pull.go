package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/boaz0/cortexpoint/pkg/cfgmgr"
	dockerclient "github.com/boaz0/cortexpoint/pkg/dockerClient"
	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull a new LLM model",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			os.Exit(1)
		}
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
			Cmd:          []string{"ollama", "pull", args[0]},
		}

		execCreate, err := client.ContainerExecCreate(ctx, cfg.Llm.Name, execOpts)
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
	rootCmd.AddCommand(pullCmd)
}
