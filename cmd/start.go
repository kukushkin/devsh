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

		if !dockerIsContainerPresent(cfg.DevContainerName) {
			dockerCmd := startDockerCmd(cfg)
			dockerRunCmd(dockerCmd)
		}

		statusDisplay(cfg)
	},
}

func startImageFlag(cmd *cobra.Command) string {
	image, err := cmd.Flags().GetString("image")
	if err != nil {
		log.Fatalf("ERROR: Failed to fetch flag: ", err)
	}
	return image
}

// Returns the configuration constructed for the  dev container, combined from all the sources and validated
func startContainerConfig(cmd *cobra.Command) ConfigValues {
	cfg := configLoad()

	// process flags set for "start" command
	image := startImageFlag(cmd)
	if image != "" {
		cfg.Image = image
	}

	// construct container volumes configuration
	primaryVolume := configDevContainerPrimaryVolume(cfg)
	cfg.DevContainerVolumes =
		append([]string{primaryVolume}, cfg.DevContainerVolumes...)

	// validate mandatory config values
	if cfg.Image == "" {
		log.Fatal("ERROR: Docker image for dev container is not specified, consider specifying it in a .devsh file")
	}

	return cfg
}

// Constructs the docker command from the dev container configuration
func startDockerCmd(cfg ConfigValues) string {
	opts := []string{
		"--name " + cfg.Name,
		"--hostname " + cfg.DevContainerHost,
		"--workdir " + cfg.DevContainerDir,
		"--detach",
	}
	if cfg.DevContainerNetwork != "" {
		opts = append(opts, "--network "+cfg.DevContainerNetwork)
	}
	if cfg.DevContainerDNS != "" {
		opts = append(opts, "--dns "+cfg.DevContainerDNS)
	}
	for _, ports := range cfg.DevContainerPorts {
		opts = append(opts, "--publish "+ports)
	}
	for _, volume := range cfg.DevContainerVolumes {
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
	startCmd.Flags().StringP("image", "i", "", "Use this docker image for the dev container")
}
