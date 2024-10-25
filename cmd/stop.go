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
		cfg := configLoad()

		if dockerIsContainerPresent(cfg.DevContainerName) {
			stopCmd := dockerConstructCmd("stop", nil, cfg.DevContainerName)
			dockerRunCmd(stopCmd)
			rmCmd := dockerConstructCmd("rm", nil, cfg.DevContainerName)
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
