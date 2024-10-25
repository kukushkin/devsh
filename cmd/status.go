// Copyright 2024 The devsh authors

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of the dev container for the current project",
	Long: `Shows status of the development container for the current project.
`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := configLoad()
		statusDisplay(cfg)
	},
}

func statusDisplay(cfg ConfigValues) {
	containerName := cfg.DevContainerName
	if dockerIsContainerPresent(containerName) {
		if dockerIsContainerRunning(containerName) {
			fmt.Printf("* Dev container %s is running (%s)\n", containerName, dockerContainerIdShort(containerName))
		} else {
			fmt.Printf("* Dev container %s is stopped (%s)\n", containerName, dockerContainerIdShort(containerName))
		}
	} else {
		fmt.Printf("* Dev container %s does not exist (stopped and/or removed)\n", containerName)
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
