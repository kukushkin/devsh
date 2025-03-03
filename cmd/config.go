// Copyright 2024 The devsh authors

package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	globalConfigFilename string = "~/.devsh/config"
	configFilename       string = ".devsh"
)

type GlobalConfigValues struct {
	Image               string   `yaml:"image,omitempty"`
	ShellCmd            string   `yaml:"shell_cmd,omitempty"`
	DevContainerVolumes []string `yaml:"dev_container_volumes,omitempty"`
	DevContainerNetwork string   `yaml:"dev_container_network,omitempty"`
	DevContainerDNS     string   `yaml:"dev_container_dns,omitempty"`
}

var defaultGlobalConfigValues = GlobalConfigValues{
	Image:               "",
	ShellCmd:            "/bin/bash",
	DevContainerVolumes: nil,
	DevContainerNetwork: "",
	DevContainerDNS:     "",
}

// Project config values replace global config values if set
type ConfigValues struct {
	Image               string   `yaml:"image,omitempty"`
	Name                string   `yaml:"name,omitempty"`
	ShellCmd            string   `yaml:"shell_cmd,omitempty"`
	DevContainerHost    string   `yaml:"dev_container_host,omitempty"`
	DevContainerDir     string   `yaml:"dev_container_dir,omitempty"`
	DevContainerName    string   `yaml:"dev_container_name,omitempty"`
	DevContainerVolumes []string `yaml:"dev_container_volumes,omitempty"`
	DevContainerNetwork string   `yaml:"dev_container_network,omitempty"`
	DevContainerDNS     string   `yaml:"dev_container_dns,omitempty"`
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
  dev_container_dns: # explicit DNS server to use for the dev container

  TBD
`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := configLoad()

		configYaml, err := yaml.Marshal(cfg)
		if err != nil {
			log.Fatalf("ERROR: Failed to serialize config: ", err)
		}

		fmt.Printf("---\n%s\n", string(configYaml))
	},
}

// Loads and returns a combined config for the project in the current folder.
func configLoad() ConfigValues {
	globalConfigValues := configLoadGlobal()
	configValues := configLoadLocal()

	// fill empty values with defaults from the global configuration
	if configValues.Image == "" {
		configValues.Image = globalConfigValues.Image
	}
	if configValues.ShellCmd == "" {
		configValues.ShellCmd = globalConfigValues.ShellCmd
	}
	if len(configValues.DevContainerVolumes) == 0 {
		configValues.DevContainerVolumes = globalConfigValues.DevContainerVolumes
	}
	if configValues.DevContainerNetwork == "" {
		configValues.DevContainerNetwork = globalConfigValues.DevContainerNetwork
	}
	if configValues.DevContainerDNS == "" {
		configValues.DevContainerDNS = globalConfigValues.DevContainerDNS
	}

	// fill values that are still empty with dynamically constructed conventional defaults
	if configValues.Name == "" {
		configValues.Name = configDefaultProjectName()
	}

	if configValues.DevContainerHost == "" {
		configValues.DevContainerHost = configDefaultDevContainerHost(configValues)
	}
	if configValues.DevContainerDir == "" {
		configValues.DevContainerDir = configDefaultDevContainerDir(configValues)
	}
	if configValues.DevContainerName == "" {
		configValues.DevContainerName = configDefaultDevContainerName(configValues)
	}
	if configValues.DevContainerNetwork == "" {
		configValues.DevContainerNetwork = configDefaultDevContainerNetwork(configValues)
	}
	if configValues.DevContainerDNS == "" {
		configValues.DevContainerDNS = configDefaultDevContainerDNS(configValues)
	}

	return configValues
}

func configLoadLocal() ConfigValues {
	// If the config file does not exist, return an empty configuration
	_, err := os.Stat(configFilename)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("WARN: Config file is not found: %s\n", configFilename)
		return ConfigValues{}
	}

	configFile, err := os.ReadFile(configFilename)
	if err != nil {

		log.Fatalf("ERROR: Failed to open config file: %s", err)
		panic(err)
	}

	var configValues ConfigValues
	yaml.Unmarshal(configFile, &configValues)

	return configValues
}

func configLoadGlobal() GlobalConfigValues {
	// If the config file does not exist, return an empty/default configuration
	_, err := os.Stat(globalConfigFilename)
	if errors.Is(err, os.ErrNotExist) {
		return defaultGlobalConfigValues
	}

	configFile, err := os.ReadFile(globalConfigFilename)
	if err != nil {

		log.Fatalf("ERROR: Failed to open config file: %s", err)
	}

	var globalConfigValues GlobalConfigValues
	yaml.Unmarshal(configFile, &globalConfigValues)

	// fill empty values with defaults from the global configuration
	if globalConfigValues.Image == "" {
		globalConfigValues.Image = defaultGlobalConfigValues.Image
	}
	if globalConfigValues.ShellCmd == "" {
		globalConfigValues.ShellCmd = defaultGlobalConfigValues.ShellCmd
	}
	if len(globalConfigValues.DevContainerVolumes) == 0 {
		globalConfigValues.DevContainerVolumes = defaultGlobalConfigValues.DevContainerVolumes
	}
	if globalConfigValues.DevContainerNetwork == "" {
		globalConfigValues.DevContainerNetwork = defaultGlobalConfigValues.DevContainerNetwork
	}
	if globalConfigValues.DevContainerDNS == "" {
		globalConfigValues.DevContainerDNS = defaultGlobalConfigValues.DevContainerDNS
	}

	return globalConfigValues
}

// Returns the project folder (i.e. current folder) on the host
func configProjectDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return cwd
}

// Returns the project name derived from the name of the project (i.e. current) folder
func configDefaultProjectName() string {
	return filepath.Base(configProjectDir())
}

// Returns the default hostname for the dev container
func configDefaultDevContainerHost(configValues ConfigValues) string {
	return configValues.Name
}

// Returns the default path where the project is mounted inside the dev container
func configDefaultDevContainerDir(configValues ConfigValues) string {
	return "/" + configValues.Name
}

// Returns the default name for the dev container
func configDefaultDevContainerName(configValues ConfigValues) string {
	return configValues.Name
}

// Returns the default network for the dev container
func configDefaultDevContainerNetwork(_configValues ConfigValues) string {
	return ""
}

// Returns the default DNS for the dev container
func configDefaultDevContainerDNS(_configValues ConfigValues) string {
	return ""
}

// Returns the primary volume (mounting the project folder) for the dev container
// Example:
//
//	configDevContainerPrimaryVolume() // "/home/alex/Projects/devsh:/devsh"
func configDevContainerPrimaryVolume(configValues ConfigValues) string {
	return configProjectDir() + ":" + configValues.DevContainerDir
}

func init() {
	rootCmd.AddCommand(configCmd)
}
