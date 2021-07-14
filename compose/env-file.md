---
description: Declare default environment variables in a file
keywords: fig, composition, compose, docker, orchestration, environment, env file
title: Declare default environment variables in file
---

Compose supports declaring default environment variables in an environment file
named `.env` placed in the project directory. Docker Compose versions earlier than `1.28`,
load the `.env` file from the current working directory, where the command is executed, or from the
project directory if this is explicitly set with the `--project-directory` option. This
inconsistency has been addressed starting with `+v1.28` by limiting the default `.env` file path
to the project directory. You can use the `--env-file` commandline option to override the default
`.env` and specify the path to a custom environment file.

The project directory is specified by the order of precedence:

- `--project-directory` flag
- Folder of the first `--file` flag
- Current directory

## Syntax rules

The following syntax rules apply to the `.env` file:

- Compose expects each line in an `env` file to be in `VAR=VAL` format.
- Lines beginning with `#` are processed as comments and ignored.
- Blank lines are ignored.
- There is no special handling of quotation marks. This means that
  **they are part of the VAL**.

### Example file
A example for the .env file:
```
# This is comment
PORT=3306
MYSQL_ROOT_PASSWORD=root
MYSQL_DATABASE=pass
```

## Compose file and CLI variables

The environment variables you define here are used for
[variable substitution](compose-file/compose-file-v3.md#variable-substitution)
in your Compose file, and can also be used to define the following
[CLI variables](reference/envvars.md):

- `COMPOSE_API_VERSION`
- `COMPOSE_CONVERT_WINDOWS_PATHS`
- `COMPOSE_FILE`
- `COMPOSE_HTTP_TIMEOUT`
- `COMPOSE_PROFILES`
- `COMPOSE_PROJECT_NAME`
- `COMPOSE_TLS_VERSION`
- `DOCKER_CERT_PATH`
- `DOCKER_HOST`
- `DOCKER_TLS_VERIFY`

> **Notes**
>
> * Values present in the environment at runtime always override those defined
>   inside the `.env` file. Similarly, values passed via command-line arguments
>   take precedence as well.
> * Environment variables defined in the `.env` file are not automatically
>   visible inside containers. To set container-applicable environment variables,
>   follow the guidelines in the topic
>   [Environment variables in Compose](environment-variables.md), which
>   describes how to pass shell environment variables through to containers,
>   define environment variables in Compose files, and more.

## More Compose documentation

- [User guide](index.md)
- [Installing Compose](install.md)
- [Getting Started](gettingstarted.md)
- [Command line reference](reference/index.md)
- [Compose file reference](compose-file/index.md)
- [Sample apps with Compose](samples-for-compose.md)
