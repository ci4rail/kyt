# kyt-cli
This repository contains kyt-cli sources, build environment and ci/cd pipeline.

Dependencies:
* git pre-commit hook
* docker

[Setup dependencies.](SetupDependencies.md)

## git pre-commit hook

Git pre-commit hook is used to ensure that a minimum set of quality is fullfilled by the checked in contents.

The configured pre-commit hooks are:

* Remove trailing whitespaces
* Ensure files end in a newline and only a newline
* Check yaml syntax
* Prevent from committing a large files
* Check python code style with flake8

## docker

As dobi is used for task automation, docker is required.

dobi is used to perform all actions and to provide a uniform interface.

Everything is packed into containers as far as possible and every action is linked to each other.

It is important not to call dobi directly, but to use the script 'dobi.sh'.

With `./dobi.sh list` you can get all 'annotated' tasks.

## Build build container
Build docker container for building kyt-cli.

Related dobi-tasks:

* image-kyt-cli-builder

To push docker container to docker registry execute `./dobi.sh image-kyt-cli-builder:push`. This requires `docker login harbor.ci4rail.com` to be executed before. See [Confluence Documentation](https://ci4rail.atlassian.net/l/c/61KodS7x) for further information.

## Build kyt-cli

Build kyt-cli within docker container.

Related dobi-tasks:

* build
* build-kyt-cli
