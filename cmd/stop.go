// Copyright 2024 The devsh authors

package cmd

import (
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the dev container",
	Long: `Stop the development container.

`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := configLoad(cmd)

		if dockerIsContainerPresent(cfg.ContainerName) {
			// Stop the container gracefully with a short timeout. `docker stop`
			// sends SIGTERM to PID 1 and waits up to the timeout (here 1s) for
			// the process to exit on its own before escalating to SIGKILL. The
			// container's main process is an idle shell that typically does not
			// handle SIGTERM, so the short timeout keeps `devsh stop` fast while
			// still giving any background processes inside a brief grace period
			// to shut down.
			stopCmd := dockerConstructCmd("stop", []string{"-t 1"}, cfg.ContainerName)
			dockerRunCmd(stopCmd)
			rmCmd := dockerConstructCmd("rm", nil, cfg.ContainerName)
			dockerRunCmd(rmCmd)
		}

		statusDisplay(cfg)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
