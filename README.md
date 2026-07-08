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
brew install --cask kukushkin/tap/devsh
```

### Linux: Using "go install"

```
go install github.com/kukushkin/devsh@latest
```

This downloads and compiles `devsh` and places it into `$GOPATH/bin`.

## How to use


### Configure your project

Place a `.devsh` file into the root folder of your project.

It should contain at least one line:
```yaml
image: dev-go # docker image name to use for your development container
```

### Starting a shell in the development container

Run in your project root folder:

```
devsh
```
