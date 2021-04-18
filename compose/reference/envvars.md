---
description: Compose CLI environment variables
keywords: fig, composition, compose, docker, orchestration, cli, reference
title: Compose CLI environment variables
---

Several environment variables are available for you to configure the Docker Compose command-line behaviour.

Variables starting with `DOCKER_` are the same as those used to configure the
Docker command-line client. If you're using `docker-machine`, then the `eval "$(docker-machine env my-docker-vm)"` command should set them to their correct values. (In this example, `my-docker-vm` is the name of a machine you created.)

> **Note**: Some of these variables can also be provided using an
> [environment file](../env-file.md).

## COMPOSE\_PROJECT\_NAME

Sets the project name. This value is prepended along with the service name to
the container on start up. For example, if your project name is `myapp` and it
includes two services `db` and `web`, then Compose starts containers named
`myapp_db_1` and `myapp_web_1` respectively.

Setting this is optional. If you do not set this, then Compose will 
default to using the `basename` of the project directory, and `COMPOSE_PROJECT_NAME`
defaults to an empty string. See also the `-p`
[command-line option](overview.md).

## COMPOSE\_FILE

Specify the path to a Compose file. If not provided, Compose looks for a file named
`docker-compose.yml` in the current directory and then each parent directory in
succession until a file by that name is found.

This variable supports multiple Compose files separated by a path separator (on
Linux and macOS the path separator is `:`, on Windows it is `;`). For example:
`COMPOSE_FILE=docker-compose.yml:docker-compose.prod.yml`. The path separator
can also be customized using `COMPOSE_PATH_SEPARATOR`.

See also the `-f` [command-line option](overview.md).

## COMPOSE\_PROFILES

Specify one or multiple active profiles to enable. Calling `docker-compose up`
with `COMPOSE_PROFILES=frontend` will start the services with the profile
`frontend` and services without specified profiles.

You can specify a list of profiles separated with a comma:
`COMPOSE_PROFILES=frontend,debug` will enable the profiles `frontend` and
`debug`.

See also [_Using profiles with Compose_](../profiles.md) and the `--profile`
[command-line option](overview.md).

## COMPOSE\_API\_VERSION

The Docker API only supports requests from clients which report a specific
version. If you receive a `client and server don't have same version` error using
`docker-compose`, you can workaround this error by setting this environment
variable. Set the version value to match the server version.

Setting this variable is intended as a workaround for situations where you need
to run temporarily with a mismatch between the client and server version. For
example, if you can upgrade the client but need to wait to upgrade the server.

Running with this variable set and a known mismatch does prevent some Docker
features from working properly. The exact features that fail would depend on the
Docker client and server versions. For this reason, running with this variable
set is only intended as a workaround and it is not officially supported.

If you run into problems running with this set, resolve the mismatch through
upgrade and remove this setting to see if your problems resolve before notifying
support.

## DOCKER\_HOST

Sets the URL of the `docker` daemon. As with the Docker client, defaults to `unix:///var/run/docker.sock`.

## DOCKER\_TLS\_VERIFY

When set to anything other than an empty string, enables TLS communication with
the `docker` daemon.

## DOCKER\_CERT\_PATH

Configures the path to the `ca.pem`, `cert.pem`, and `key.pem` files used for TLS verification. Defaults to `~/.docker`.

## COMPOSE\_HTTP\_TIMEOUT

Configures the time (in seconds) a request to the Docker daemon is allowed to hang before Compose considers
it failed. Defaults to 60 seconds.

## COMPOSE\_TLS\_VERSION

Configure which TLS version is used for TLS communication with the `docker`
daemon. Defaults to `TLSv1`.
Supported values are: `TLSv1`, `TLSv1_1`, `TLSv1_2`.

## COMPOSE\_CONVERT\_WINDOWS\_PATHS

Enable path conversion from Windows-style to Unix-style in volume definitions.
Users of Docker Machine on Windows should always set this. Defaults to `0`.
Supported values: `true` or `1` to enable, `false` or `0` to disable.

## COMPOSE\_PATH\_SEPARATOR

If set, the value of the `COMPOSE_FILE` environment variable is separated
using this character as path separator.

## COMPOSE\_FORCE\_WINDOWS\_HOST

If set, volume declarations using the [short syntax](../compose-file/#short-syntax-3)
are parsed assuming the host path is a Windows path, even if Compose is
running on a UNIX-based system.
Supported values: `true` or `1` to enable, `false` or `0` to disable.

## COMPOSE\_IGNORE\_ORPHANS

If set, Compose doesn't try to detect orphaned containers for the project.
Supported values: `true` or `1` to enable, `false` or `0` to disable.

## COMPOSE\_PARALLEL\_LIMIT

Sets a limit for the number of operations Compose can execute in parallel. The
default value is `64`, and may not be set lower than `2`.

## COMPOSE\_INTERACTIVE\_NO\_CLI

If set, Compose doesn't attempt to use the Docker CLI for interactive `run`
and `exec` operations. This option is not available on Windows where the CLI
is required for the aforementioned operations.
Supported: `true` or `1` to enable, `false` or `0` to disable.

## COMPOSE\_DOCKER\_CLI\_BUILD

Configure whether to use the Compose python client for building images or the
native docker cli. By default, Compose uses the `docker` CLI to perform builds,
which allows you to use [BuildKit](../../develop/develop-images/build_enhancements.md#to-enable-buildkit-builds)
to perform builds.

Set `COMPOSE_DOCKER_CLI_BUILD=0` to disable native builds, and to use the built-in
python client.

## Related information

- [User guide](../index.md)
- [Installing Compose](../install.md)
- [Compose file reference](../compose-file/index.md)
- [Environment file](../env-file.md)
