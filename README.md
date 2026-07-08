[![cli-app-release](https://github.com/kukushkin/devsh/actions/workflows/cli-app-release.yaml/badge.svg)](https://github.com/kukushkin/devsh/actions/workflows/cli-app-release.yaml)

# devsh

A shell that runs inside a development container:

* Starts a development container with your development tools (e.g. NodeJS or Python)
* Mounts your current folder into the container
* Runs a shell inside the container

## Prerequisites

* Docker
* prebuild Docker images for your development containers

## Installation

### MacOS: Using Homebrew

```
brew install kukushkin/tap/devsh
```

### Linux: Using "go install"

```
go install github.com/kukushkin/devsh@latest
```

This downloads and compiles `devsh` and places it into `$GOPATH/bin`.

## How to use

### Configure your project

Place a `.devsh` file into the root folder of your project to configure the
development container. It is a YAML file; when present it must specify at least
an image:

```yaml
image: dev-go # docker image to use for the development container
```

All other keys are optional:

```yaml
image: dev-go                        # docker image for the dev container (required)
name: my-project                     # project name (default: current directory name)
shell_cmd: /bin/bash                 # shell to start inside the container (default: /bin/bash)
container_host: my-project           # hostname for the container (default: project name)
container_dir: /my-project           # path where the project is mounted inside the container
container_name: my-project           # name of the container (default: <dir_name>-<hash>)
ports:                               # container ports exposed on the host
  - 8080:8080
volumes:                             # additional volumes to mount inside the container
  - /home/alex/data:/data
network: my-network                  # docker network for the container
dns: 8.8.8.8                         # explicit DNS server for the container
```

### Global configuration

A global configuration file can be placed at `~/.config/devsh`. It uses the same
format as the project `.devsh` file and provides defaults for all projects. The
location of the global config file can be overridden with the `DEVSH_CONFIG`
environment variable:

```
DEVSH_CONFIG=~/my-devsh-config devsh
```

Configuration values are combined from the following sources, listed from the
lowest priority to the highest:

1. Built-in defaults
2. Global configuration file (`~/.config/devsh` or `$DEVSH_CONFIG`)
3. Project configuration file (`.devsh` in the current folder)
4. Command-line flags

For each parameter, the value from the highest-priority source that provides it
takes precedence; values from lower-priority sources are inherited when a
higher-priority source does not set the parameter.

### Minimal configuration (without a `.devsh` file)

A `.devsh` file is not required. You can run `devsh` in any folder as long as
the docker image is provided through the global configuration or the `--image`
flag. For example, set a default image once for all projects:

```
echo "image: dev-go" > ~/.config/devsh
```

then run `devsh` in any folder:

```
devsh
```

or pass the image on the command line:

```
devsh --image dev-go
```

All other parameters fall back to their built-in defaults (e.g. the container
name is derived from the current directory name).

### Starting a shell in the development container

Run in your project root folder:

```
devsh
```

This starts the development container (if it is not running yet) and opens a
shell into it.

### Commands

| Command | Description |
|---|---|
| `devsh` | Start the container (if needed) and open a shell (default action) |
| `devsh start` | Start the development container |
| `devsh open` | Open a shell in the running container |
| `devsh status` | Show the status of the container |
| `devsh stop` | Stop and remove the container |
| `devsh config` | Show the effective configuration for the current project |

### Command-line flags

Every configuration parameter can also be set via a command-line flag. Flags
have the highest priority and override values from the config files:

| Flag | Description |
|---|---|
| `-i, --image` | Docker image for the dev container |
| `-n, --name` | Name of the project |
| `-s, --shell-cmd` | Shell to start inside the dev container |
| `--container-host` | Hostname for the dev container |
| `--container-dir` | Path inside the container where the project is mounted |
| `--container-name` | Human-readable name for the dev container |
| `-p, --ports` | Ports of the container exposed on the host |
| `-V, --volumes` | Additional volumes to be mounted inside the dev container |
| `--network` | Docker network for the dev container |
| `--dns` | Explicit DNS server to use for the dev container |

Use `-v`/`--verbose` to print the docker commands devsh runs.

Example:

```
devsh --image dev-go --network my-network -p 8080:8080
```
