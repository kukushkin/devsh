// Copyright 2024 The devsh authors

package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the dev container for the current project",
	Long: `Start the development container for the project in the current folder.

`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := startContainerConfig(cmd)

		if !dockerIsContainerPresent(cfg.ContainerName) {
			dockerCmd := startDockerCmd(cfg)
			dockerRunCmd(dockerCmd)
		}

		statusDisplay(cfg)
	},
}

// Returns the configuration constructed for the dev container, combined from
// all the sources and validated.
func startContainerConfig(cmd *cobra.Command) ConfigValues {
	cfg := configLoad(cmd)

	// construct container volumes configuration
	primaryVolume := configDevContainerPrimaryVolume(cfg)
	cfg.Volumes =
		append([]string{primaryVolume}, cfg.Volumes...)

	// validate mandatory config values
	if cfg.Image == "" {
		log.Fatal("ERROR: Docker image for dev container is not specified, consider specifying it in a .devsh file")
	}

	return cfg
}

// Constructs the docker command from the dev container configuration
func startDockerCmd(cfg ConfigValues) string {
	opts := []string{
		"--name " + cfg.ContainerName,
		"--hostname " + cfg.ContainerHost,
		"--workdir " + cfg.ContainerDir,
		"--detach",
	}
	if cfg.Network != "" {
		opts = append(opts, "--network "+cfg.Network)
	}
	if cfg.DNS != "" {
		opts = append(opts, "--dns "+cfg.DNS)
	}
	for _, ports := range cfg.Ports {
		opts = append(opts, "--publish "+ports)
	}
	for _, volume := range cfg.Volumes {
		opts = append(opts, "--volume "+volume)
	}

	return dockerConstructCmd("run", opts, cfg.Image)
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().StringP("image", "i", "", "Docker images for the dev container")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
