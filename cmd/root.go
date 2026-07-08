// Copyright 2024 The devsh authors

package cmd

import (
	"os"

	"github.com/kukushkin/devsh/version"
	"github.com/spf13/cobra"
)

const VERSION_TEMPLATE = "devsh {{.Version}}\n"

var (
	globalFlagVerbose bool = false
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "devsh",
	Version: version.Version,
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
		if !(dockerIsContainerPresent(cfg.ContainerName) && dockerIsContainerRunning(cfg.ContainerName)) {
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

	// Configuration flags mirror the configurable parameters and have the
	// highest priority, overriding values from the global and project config
	// files. They are persistent so they apply to every subcommand.
	rootCmd.PersistentFlags().StringP("image", "i", "", "Docker image for the dev container")
	rootCmd.PersistentFlags().StringP("name", "n", "", "Name of the project")
	rootCmd.PersistentFlags().StringP("shell-cmd", "s", "", "Shell to start inside the dev container (e.g. /bin/bash)")
	rootCmd.PersistentFlags().String("container-host", "", "Hostname for the dev container")
	rootCmd.PersistentFlags().String("container-dir", "", "Path inside the dev container where the project is mounted")
	rootCmd.PersistentFlags().String("container-name", "", "Human-readable name for the dev container")
	rootCmd.PersistentFlags().StringSliceP("ports", "p", nil, "Ports of the container exposed on the host")
	rootCmd.PersistentFlags().StringSliceP("volumes", "V", nil, "Additional volumes to be mounted inside the dev container")
	rootCmd.PersistentFlags().String("network", "", "Docker network for the dev container")
	rootCmd.PersistentFlags().String("dns", "", "Explicit DNS server to use for the dev container")

	rootCmd.SetVersionTemplate(VERSION_TEMPLATE)
}
