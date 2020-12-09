---
description: Describes how to use the json-file logging driver.
keywords: json-file, docker, logging, driver
redirect_from:
- /engine/reference/logging/json-file/
- /engine/admin/logging/json-file/
title: JSON File logging driver
---


By default, Docker captures the standard output (and standard error) of all your containers,
and writes them in files using the JSON format. The JSON format annotates each line with its
origin (`stdout` or `stderr`) and its timestamp. Each log file contains information about
only one container.

```json
{"log":"Log line is here\n","stream":"stdout","time":"2019-01-01T11:11:11.111111111Z"}
```

## Usage

To use the `json-file` driver as the default logging driver, set the `log-driver`
and `log-opts` keys to appropriate values in the `daemon.json` file, which is
located in `/etc/docker/` on Linux hosts or
`C:\ProgramData\docker\config\` on Windows Server. For more information about
configuring Docker using `daemon.json`, see
[daemon.json](../../../engine/reference/commandline/dockerd.md#daemon-configuration-file).

The following example sets the log driver to `json-file` and sets the `max-size` and `max-file` options.

```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3" 
  }
}
```

> **Note**
>
> `log-opts` configuration options in the `daemon.json` configuration file must
> be provided as strings. Boolean and numeric values (such as the value for
> `max-file` in the example above) must therefore be enclosed in quotes (`"`).

Restart Docker for the changes to take effect for newly created containers.
Existing containers do not use the new logging configuration.

You can set the logging driver for a specific container by using the
`--log-driver` flag to `docker container create` or `docker run`:

```bash
$ docker run \
      --log-driver json-file --log-opt max-size=10m \
      alpine echo hello world
```

### Options

The `json-file` logging driver supports the following logging options:

| Option         | Default           | Description                                                                                                                                                                                    | Example  value                                    |
|:---------------|:------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:--------------------------------------------------|
| `max-size`     | `-1` (unlimited)  | The maximum size of the log before it is rolled. A positive integer plus a modifier representing the unit of measure (`k`, `m`, or `g`).                                                       | `--log-opt max-size=10m`                          |
| `max-file`     | `1` (no rotation) | The maximum number of log files that can be present. If rolling the logs creates excess files, the oldest file is removed. **Only effective when `max-size` is also set.** A positive integer. | `--log-opt max-file=3`                            |
| `labels`       |                   | Applies when starting the Docker daemon. A comma-separated list of logging-related labels this daemon accepts. Used for advanced [log tag options](log_tags.md).                               | `--log-opt labels=production_status,geo`          |
| `env`          |                   | Applies when starting the Docker daemon. A comma-separated list of logging-related environment variables this daemon accepts. Used for advanced [log tag options](log_tags.md).                | `--log-opt env=os,customer`                       |
| `env-regex`    |                   | Similar to and compatible with `env`. A regular expression to match logging-related environment variables. Used for advanced [log tag options](log_tags.md).                                   | `--log-opt env-regex=^(os|customer).`             |
| `compress`     | `false`           | Toggles compression for rotated logs. Disabled by default.                                                                                                                                     | `--log-opt compress=true`                         |


### Examples

This example starts an `alpine` container which can have a maximum of 3 log
files no larger than 10 megabytes each.

```bash
$ docker run -it --log-opt max-size=10m --log-opt max-file=3 alpine ash
```
