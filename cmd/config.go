// Copyright 2024 The devsh authors

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const configFilename string = ".devsh"

type ConfigValues struct {
	Image               string   `yaml:"image"`
	Name                string   `yaml:"name"`
	ShellCmd            string   `yaml:"shell_cmd"`
	DevContainerHost    string   `yaml:"dev_container_host"`
	DevContainerDir     string   `yaml:"dev_container_dir"`
	DevContainerName    string   `yaml:"dev_container_name"`
	DevContainerVolumes []string `yaml:"dev_container_volumes"`
	DevContainerNetwork string   `yaml:"dev_container_network"`
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the current project configuration",
	Long: `Show the configuration of the project in the current folder.

The configuration of the project is combined from the global devsh configuration
found in (TBD) and the local configuration file placed in the project root directory.

The local project configuration file should be named '.devsh' and placed in the root
directory of the project.

.devsh file is a YAML file with the following format (all keys are optional):
  image: # docker image to be used for dev container
  name: # name of the project, if omitted the directory name is used
  TBD_devenv: # name of the docker compose environment, to attach to
  shell_cmd: # shell to start inside the dev container, e.g. /bin/bash
  dev_container_host: # name of the host for the dev container
  dev_container_dir: # path inside the dev container where the project is going to be mounted
  dev_container_name: # human-readable name for the dev container in docker
  dev_container_volumes: # additional volumes to be mounted inside the dev container
  dev_container_network: # docker network for the dev container

  TBD
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
		configValues := configLoad()

		fmt.Printf("%v\n", configValues)
	},
}

func configLoad() ConfigValues {
	configFile, err := os.ReadFile(configFilename)
	if err != nil {
		log.Fatalf("ERROR: Failed to open config file: %s", err)
		panic(err)
	}

	var configValues ConfigValues
	yaml.Unmarshal(configFile, &configValues)

	return configValues
}

func init() {
	rootCmd.AddCommand(configCmd)
}
