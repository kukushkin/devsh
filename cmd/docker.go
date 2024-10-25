package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	dockerCli         = "docker"
	dockerIdShortSize = 12 // hex characters
)

// Constructs a command to a docker and returns it as a string.
//
// The dockerCommand parameter is a command to a docker CLI, e.g. "run", "stop" etc.
// opts is a list of optional flags and args are additional arguments to a command.
//
// For example:
//
//	dockerConstructCmd("run", []string{"--rm",}, "ubuntu") // docker run --rm ubuntu
func dockerConstructCmd(dockerCommand string, opts []string, args ...string) string {
	var cmd []string

	cmd = append(cmd, dockerCli)
	cmd = append(cmd, dockerCommand)

	for _, opt := range opts {
		cmd = append(cmd, opt)
	}

	for _, arg := range args {
		cmd = append(cmd, arg)
	}

	return strings.Join(cmd, " ")
}

// Runs a shell command and returns the resulting output as a string
func dockerRunCmd(shellCommand string) string {
	if globalFlagVerbose {
		fmt.Println("+ " + shellCommand) // if echo/verbose
	}
	out, err := exec.Command("/bin/sh", "-c", shellCommand).Output()
	if err != nil {
		log.Fatalf("ERROR: Failed to run docker command: %s", err)
	}

	return strings.TrimSpace(string(out))
}

// Runs a command in an interactive shell
func dockerRunInteractive(shellCommand string) {
	if globalFlagVerbose {
		fmt.Println("+ " + shellCommand) // if echo/verbose
	}
	cmd := exec.Command("/bin/sh", "-c", shellCommand)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run() // add error checking
	if err != nil {
		log.Fatalf("WARN: Shell exited with an error: %s", err)
	}
}

// Returns true if the container with given name is present (either running or stopped)
func dockerIsContainerPresent(name string) bool {
	shellCmd := dockerConstructCmd("inspect", nil, name)
	if globalFlagVerbose {
		fmt.Println("+ " + shellCmd) // if echo/verbose
	}
	_, err := exec.Command("/bin/sh", "-c", shellCmd).Output()
	if err != nil {
		return false
	} else {
		return true
	}
}

// Returns true if the container with the given name is running
func dockerIsContainerRunning(name string) bool {
	opts := []string{
		"-f '{{.State.Status}}'",
	}
	shellCmd := dockerConstructCmd("inspect", opts, name)
	out := dockerRunCmd(shellCmd)
	if out == "running" {
		return true
	} else {
		return false
	}
}

// Returns ID of the docker container with the given name
func dockerContainerId(name string) string {
	opts := []string{
		"-f '{{.Id}}'",
	}
	dockerCmd := dockerConstructCmd("inspect", opts, name)
	out := dockerRunCmd(dockerCmd)

	return out
}

// Returns shortened ID of the docker container with the given name
func dockerContainerIdShort(name string) string {
	return dockerContainerId(name)[0:dockerIdShortSize]
}
