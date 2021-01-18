# kyt-cli
This repository contains kyt-cli and kyt-api-server sources, build environment and ci/cd pipeline (ci/cd currently missing).

Dependencies:
* git pre-commit hook
* golang (required for pre-commit hooks for golang)
* docker

[Setup dependencies.](SetupDependencies.md)

## Build

### KYT-API-SERVER

Containerized Build of the kyt-api-server tool. Builds x86 version for linux.

```bash
$ ./dobi.sh build-kyt-apiserver
```

Get the iot hub connection string from the Azure Portal.
Select the iot hub, then "shared access policies"

Run the kyt-api-server:

```bash
export IOTHUB_SERVICE_CONNECTION_STRING="HostName=ci4rail-eval-iot-hub.azure-devices.net;SharedAccessKeyName=iothubowner;SharedAccessKey=6...="
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

## Repo Notes

### OpenAPI specification

The REST-API is defined with the OpenAPI document `kyt-api-spec/kytapi.yaml`. From this specification, parts of the server and client go code are automatically generated.

### Auto-Generated code

Openapi-generator (https://openapi-generator.tech/) is used to generate the go code for the server and client.

Note that the generated files are not placed in git.


### git pre-commit hook

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

## docker

As dobi is used for task automation, docker is required.

dobi is used to perform all actions and to provide a uniform interface.

Everything is packed into containers as far as possible and every action is linked to each other.

It is important not to call dobi directly, but to use the script 'dobi.sh'.

With `./dobi.sh list` you can get all 'annotated' tasks.

## Build build container
Build docker container for building kyt-cli.

Related dobi-tasks:

* image-kyt-go-builder

To push docker container to docker registry execute `./dobi.sh image-kyt-go-builder:push`. This requires `docker login harbor.ci4rail.com` to be executed before. See [Confluence Documentation](https://ci4rail.atlassian.net/l/c/61KodS7x) for further information.
