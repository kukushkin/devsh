# devsh
Development shell to run inside a docker container

## Prerequisites

* docker
* prebuild docker images for your development containers

## Installation

### Using "go install"

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
