# kyt-cli
[![CI](https://concourse.ci4rail.com/api/v1/teams/main/pipelines/kyt-services/jobs/build-kyt-services/badge)](https://concourse.ci4rail.com/teams/main/pipelines/kyt-services) [![Go Report Card](https://goreportcard.com/badge/github.com/ci4rail/kyt-cli)](https://goreportcard.com/report/github.com/ci4rail/kyt-cli)

This repository contains kyt-cli and kyt-api-server sources, build environment and CI/CD pipeline.

# Build
Here you can find the build instructions for either locally with and without docker and via the CI/CD pipeline.

## Build locally

### KYT-API-SERVER
#### Docker image
To build (and deploy) the `kyt-api-server` docker image you can use the following commands:
```bash
$ ./dobi.sh image-kyt-api-server        # build only
$ ./dobi.sh image-kyt-api-server:push   # build and push do docker registry
```
To run the docker image with a specific `<tag>` use:
```bash
docker run --rm -p 8080:8080 -e IOTHUB_SERVICE_CONNECTION_STRING="HostName=ci4rail-eval-iot-hub.azure-devices.net;SharedAccessKeyName=iothubowner;SharedAccessKey=6..." harbor.ci4rail.com/ci4rail/kyt/kyt-api-server:<tag>
```
Have a look at available tags for the image: https://harbor.ci4rail.com/harbor/projects/7/repositories/kyt%2Fkyt-api-server

#### Plain binary

Containerized Build of the kyt-api-server tool. Builds x86 version for linux.

```bash
$ ./dobi.sh build-kyt-apiserver
```

Get the iot hub connection string from the Azure Portal.
Select the iot hub, then "shared access policies"

Run the kyt-api-server:

```bash
$ export IOTHUB_SERVICE_CONNECTION_STRING="HostName=ci4rail-eval-iot-hub.azure-devices.net;SharedAccessKeyName=iothubowner;SharedAccessKey=6...="
$ bin/kyt-api-server --addr :8080
```

Or, build/run it with your local go installation:

```bash
$ ./dobi.sh generate-server-sources
$ cd kyt-api-server
$ go run main.go  --addr :8080
```

#### Test-Server

Folder `kyt-api-server/test-server` contains a test-server that answers API request with dummy data.

```bash
$ cd kyt-api-server/test-server
$ go run main.go  --addr :8080
```

### KYT-CLI

Containerized Build of the kyt-cli tool. The binary is named just `kyt`. Builds x86 version for windows and linux.

```bash
$ ./dobi.sh build-kyt-cli
$ bin/kyt --server "http://localhost:8080/v1" dlm get devices
```

Or, build/run it with your local go installation:

```bash
$ ./dobi.sh generate-client-sources
$ cd kyt-cli
$ go run main.go --server "http://localhost:8080/v1" dlm get devices
```

## Build with CI/CD

[Concourse CI](https://concourse-ci.org/) will be used as CI/CD system.

## pipeline.yaml
The `pipeline.yaml` is the CI/CD pipeline that builds kyt-cli called `kyt` and kyt-api-server docker image. The kyt-cli goes to a [gitub pre-release](https://github.com/ci4rail/kyt-cli/releases) and the `kyt-api-server` will be published as docker image [here](https://harbor.ci4rail.com/harbor/projects/7/repositories/kyt%2Fkyt-api-server).

### Usage

Copy `ci/credentials.template.yaml` to `ci/credentials.yaml` and enter the credentials needed.
Download `fly` and apply the CI/CD pipeline to Concourse CI using:
```bash
# Download fly from concourse server
$ sudo curl -L https://concourse.ci4rail.com/api/v1/cli?arch=amd64&platform=linux -o /usr/local/bin/fly && sudo chmod +x /usr/local/bin/fly
# Apply pipeline
$ fly -t prod set-pipeline -p kyt-services -c pipeline.yaml -l ci/config.yaml  -l ci/credentials.yaml
```

# Repo Notes

## OpenAPI specification

The REST-API is defined with the OpenAPI document `kyt-api-spec/kytapi.yaml`. From this specification, parts of the server and client go code are automatically generated.

## Auto-Generated code

Openapi-generator (https://openapi-generator.tech/) is used to generate the go code for the server and client.

Note that the generated files are not placed in git.


## git pre-commit hook

Git pre-commit hook is used to ensure that a minimum set of quality is fullfilled by the checked in contents.

The configured pre-commit hooks are:

* Remove trailing whitespaces
* Ensure files end in a newline and only a newline
* Check yaml syntax
* Prevent from committing a large files
* Check python code style with flake8
* Autoformat go code
* Check go code style with golint
* Update go import lines (adding missing ones and removing unreferenced ones)
* Clean up unused dependencies in go.mod

## docker and dobi

As [dobi](https://github.com/dnephin/dobi) is used for task automation, docker is required.
dobi is used to perform all actions and to provide a uniform interface.

Everything is packed into containers as far as possible and every action is linked to each other.

It is important not to call dobi directly, but to use the script `dobi.sh` as it performs some version generation and several checks regarding the used dobi version itself.

Use `./dobi.sh list` to get all 'annotated' dobi resources like images or jobs.

# Development environment

If you want to contribute please see this section for the development environment and the dependencies.

## Dependencies
* git pre-commit hook
* golang (required for pre-commit hooks for golang)
* docker

See [Setup dependencies.](SetupDependencies.md) for further instructions.
