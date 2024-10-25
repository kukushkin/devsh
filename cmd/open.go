// Copyright 2024 The devsh authors

package cmd

import (
	"github.com/spf13/cobra"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open a shell in the dev container",
	Long: `Open a shell in the development container.

`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := configLoad()
		openShell(cfg)
		statusDisplay(cfg)
	},
}

func openShell(cfg ConfigValues) {
	opts := []string{
		"-ti",
	}
	shellCmd := dockerConstructCmd("exec", opts, cfg.DevContainerName, cfg.ShellCmd)
	dockerRunInteractive(shellCmd)
}

func init() {
	rootCmd.AddCommand(openCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// openCmd.PersistentFlags().StringP("image", "i", "", "Help message for image")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// openCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
