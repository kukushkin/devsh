// Copyright 2024 The devsh authors

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	Version           string = "debug-version"
	globalFlagVerbose bool   = false
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "devsh",
	Version: Version,
	Short:   "Run a shell in a development container",
	Long: `Devsh is a CLI tool that allows you to start an isolated docker container
with your project mounted inside (development container), and open a shell into it.

For example:
	devsh start # starts a development container
	devsh open # opens a shell into it
Or:
	devsh # default action starts a development container and opens a shell in one go
`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := startContainerConfig(cmd)

		// start the dev container if it is not started yet
		if !(dockerIsContainerPresent(cfg.DevContainerName) && dockerIsContainerRunning(cfg.DevContainerName)) {
			dockerCmd := startDockerCmd(cfg)
			dockerRunCmd(dockerCmd)
		}

		openShell(cfg)

		statusDisplay(cfg)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.devsh.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&globalFlagVerbose, "verbose", "v", false, "Produce verbose output")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolVarP(&globalFlagVerbose, "verbose", "v", false, "Produce verbose output")
	rootCmd.Flags().StringP("image", "i", "", "Use this docker image for the dev container")
}
