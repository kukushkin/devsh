// Copyright 2024 The devsh authors

package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	configFilename string = ".devsh"
)

// ConfigValues holds every configurable parameter of devsh. The same set of
// values can be provided by any of the three configuration sources: the global
// config file, the project .devsh file, and command-line flags.
type ConfigValues struct {
	Image         string   `yaml:"image,omitempty"`
	Name          string   `yaml:"name,omitempty"`
	ShellCmd      string   `yaml:"shell_cmd,omitempty"`
	ContainerHost string   `yaml:"container_host,omitempty"`
	ContainerDir  string   `yaml:"container_dir,omitempty"`
	ContainerName string   `yaml:"container_name,omitempty"`
	Ports         []string `yaml:"ports,omitempty"`
	Volumes       []string `yaml:"volumes,omitempty"`
	Network       string   `yaml:"network,omitempty"`
	DNS           string   `yaml:"dns,omitempty"`
}

// defaultConfigValues returns the built-in defaults. These have the lowest
// priority and are only used when a value is not provided by any of the three
// configuration sources.
func defaultConfigValues() ConfigValues {
	return ConfigValues{
		ShellCmd: "/bin/bash",
	}
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the current project configuration",
	Long: `Show the configuration of the project in the current folder.

The configuration is combined from the following sources, listed from the
lowest priority to the highest:

  1. Built-in defaults
  2. Global configuration file (default: ~/.config/devsh, overridable via
     the DEVSH_CONFIG environment variable)
  3. Project configuration file (.devsh in the current folder)
  4. Command-line flags

For every parameter, the value from the highest-priority source that provides
it takes precedence; values from lower-priority sources are inherited when a
higher-priority source does not set the parameter.

The .devsh file is a YAML file with the following format (all keys are optional):
  image: # docker image to be used for dev container
  name: # name of the project, if omitted the directory name is used
  shell_cmd: # shell to start inside the dev container, e.g. /bin/bash
  container_host: # name of the host for the dev container
  container_dir: # path inside the dev container where the project is going to be mounted
  container_name: # human-readable name for the dev container in docker
  ports: # ports of the container exposed on host
  volumes: # additional volumes to be mounted inside the dev container
  network: # docker network for the dev container
  dns: # explicit DNS server to use for the dev container
`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := configLoad(cmd)

		configYaml, err := yaml.Marshal(cfg)
		if err != nil {
			log.Fatalf("ERROR: Failed to serialize config: %s", err)
		}

		fmt.Printf("---\n%s\n", string(configYaml))
	},
}

// Loads and returns a combined config for the project in the current folder.
//
// Values are merged from four layers, each overriding the previous one:
//  1. built-in defaults
//  2. global configuration file
//  3. project configuration file (.devsh)
//  4. command-line flags
//
// Finally, values that are still empty are filled with dynamically constructed
// conventional defaults derived from the project directory.
func configLoad(cmd *cobra.Command) ConfigValues {
	cfg := defaultConfigValues()

	cfg = mergeConfig(cfg, configLoadGlobal())
	cfg = mergeConfig(cfg, configLoadLocal())
	cfg = mergeConfig(cfg, configLoadFlags(cmd))

	// fill values that are still empty with dynamically constructed defaults
	if cfg.Name == "" {
		cfg.Name = configDefaultProjectName()
	}

	if cfg.ContainerHost == "" {
		cfg.ContainerHost = configDefaultContainerHost(cfg)
	}
	if cfg.ContainerDir == "" {
		cfg.ContainerDir = configDefaultContainerDir(cfg)
	}
	if cfg.ContainerName == "" {
		cfg.ContainerName = configDefaultContainerName(cfg)
	}
	if cfg.Network == "" {
		cfg.Network = configDefaultNetwork(cfg)
	}
	if cfg.DNS == "" {
		cfg.DNS = configDefaultDNS(cfg)
	}

	return cfg
}

// mergeConfig returns base with every field overridden by the corresponding
// non-empty field from override. Empty fields in override are ignored so that
// lower-priority values are inherited. Slice fields are replaced (not
// concatenated) when override provides any values.
func mergeConfig(base, override ConfigValues) ConfigValues {
	if override.Image != "" {
		base.Image = override.Image
	}
	if override.Name != "" {
		base.Name = override.Name
	}
	if override.ShellCmd != "" {
		base.ShellCmd = override.ShellCmd
	}
	if override.ContainerHost != "" {
		base.ContainerHost = override.ContainerHost
	}
	if override.ContainerDir != "" {
		base.ContainerDir = override.ContainerDir
	}
	if override.ContainerName != "" {
		base.ContainerName = override.ContainerName
	}
	if len(override.Ports) > 0 {
		base.Ports = override.Ports
	}
	if len(override.Volumes) > 0 {
		base.Volumes = override.Volumes
	}
	if override.Network != "" {
		base.Network = override.Network
	}
	if override.DNS != "" {
		base.DNS = override.DNS
	}
	return base
}

// configLoadFlags collects values provided via command-line flags. Only flags
// that were explicitly set on the command line are taken into account, so that
// unset flags do not clobber values coming from the config files.
func configLoadFlags(cmd *cobra.Command) ConfigValues {
	var cfg ConfigValues
	if cmd == nil {
		return cfg
	}

	flags := cmd.Flags()

	if flags.Changed("image") {
		cfg.Image, _ = flags.GetString("image")
	}
	if flags.Changed("name") {
		cfg.Name, _ = flags.GetString("name")
	}
	if flags.Changed("shell-cmd") {
		cfg.ShellCmd, _ = flags.GetString("shell-cmd")
	}
	if flags.Changed("container-host") {
		cfg.ContainerHost, _ = flags.GetString("container-host")
	}
	if flags.Changed("container-dir") {
		cfg.ContainerDir, _ = flags.GetString("container-dir")
	}
	if flags.Changed("container-name") {
		cfg.ContainerName, _ = flags.GetString("container-name")
	}
	if flags.Changed("ports") {
		cfg.Ports, _ = flags.GetStringSlice("ports")
	}
	if flags.Changed("volumes") {
		cfg.Volumes, _ = flags.GetStringSlice("volumes")
	}
	if flags.Changed("network") {
		cfg.Network, _ = flags.GetString("network")
	}
	if flags.Changed("dns") {
		cfg.DNS, _ = flags.GetString("dns")
	}

	return cfg
}

// configGlobalPath returns the path to the global configuration file. The
// location can be overridden with the DEVSH_CONFIG environment variable; it
// defaults to ~/.config/devsh. A leading '~' is expanded to the user's home
// directory.
func configGlobalPath() string {
	if p := os.Getenv("DEVSH_CONFIG"); p != "" {
		return expandTilde(p)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("ERROR: Failed to determine home directory: %s", err)
	}
	return filepath.Join(home, ".config", "devsh")
}

// expandTilde replaces a leading '~' with the user's home directory.
func expandTilde(p string) string {
	if p == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return p
		}
		return home
	}
	if strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return p
		}
		return filepath.Join(home, p[2:])
	}
	return p
}

func configLoadLocal() ConfigValues {
	// If the config file does not exist, return an empty configuration
	_, err := os.Stat(configFilename)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("WARN: Config file is not found: %s\n", configFilename)
		return ConfigValues{}
	}
	if err != nil {
		log.Fatalf("ERROR: Failed to stat config file %s: %s", configFilename, err)
	}

	configFile, err := os.ReadFile(configFilename)
	if err != nil {
		log.Fatalf("ERROR: Failed to read config file %s: %s", configFilename, err)
	}

	var configValues ConfigValues
	if err := yaml.Unmarshal(configFile, &configValues); err != nil {
		log.Fatalf("ERROR: Failed to parse config file %s: %s", configFilename, err)
	}

	return configValues
}

func configLoadGlobal() ConfigValues {
	path := configGlobalPath()

	// If the config file does not exist, return an empty configuration; a
	// global config file is optional.
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return ConfigValues{}
	}
	if err != nil {
		log.Fatalf("ERROR: Failed to stat global config file %s: %s", path, err)
	}

	configFile, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("ERROR: Failed to read global config file %s: %s", path, err)
	}

	var configValues ConfigValues
	if err := yaml.Unmarshal(configFile, &configValues); err != nil {
		log.Fatalf("ERROR: Failed to parse global config file %s: %s", path, err)
	}

	return configValues
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
func configDefaultContainerHost(configValues ConfigValues) string {
	return configValues.Name
}

// Returns the default path where the project is mounted inside the dev container
func configDefaultContainerDir(configValues ConfigValues) string {
	return "/" + configValues.Name
}

// Returns the default name for the dev container
func configDefaultContainerName(configValues ConfigValues) string {
	return configValues.Name
}

// Returns the default network for the dev container
func configDefaultNetwork(_configValues ConfigValues) string {
	return ""
}

// Returns the default DNS for the dev container
func configDefaultDNS(_configValues ConfigValues) string {
	return ""
}

// Returns the primary volume (mounting the project folder) for the dev container
// Example:
//
//	configDevContainerPrimaryVolume() // "/home/alex/Projects/devsh:/devsh"
func configDevContainerPrimaryVolume(configValues ConfigValues) string {
	return configProjectDir() + ":" + configValues.ContainerDir
}

func init() {
	rootCmd.AddCommand(configCmd)
}
