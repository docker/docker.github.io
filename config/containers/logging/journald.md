---
description: Describes how to use the Journald logging driver.
keywords: Journald, docker, logging, driver
redirect_from:
- /engine/reference/logging/journald/
- /engine/admin/logging/journald/
title: Journald logging driver
---

The `journald` logging driver sends container logs to the
[`systemd` journal](https://www.freedesktop.org/software/systemd/man/systemd-journald.service.html).
Log entries can be retrieved using the `journalctl` command, through use of the
`journal` API, or using the `docker logs` command.

In addition to the text of the log message itself, the `journald` log driver
stores the following metadata in the journal with each message:

| Field                                | Description                                                                                                                                            |
|:-------------------------------------|:-------------------------------------------------------------------------------------------------------------------------------------------------------|
| `CONTAINER_ID`                       | The container ID truncated to 12 characters.                                                                                                           |
| `CONTAINER_ID_FULL`                  | The full 64-character container ID.                                                                                                                    |
| `CONTAINER_NAME`                     | The container name at the time it was started. If you use `docker rename` to rename a container, the new name is not reflected in the journal entries. |
| `CONTAINER_TAG`, `SYSLOG_IDENTIFIER` | The container tag ([log tag option documentation](log_tags.md)).                                                                                       |
| `CONTAINER_PARTIAL_MESSAGE`          | A field that flags log integrity. Improve logging of long log lines.                                                                                   |
| `IMAGE_NAME`                         | The image name and tag for the container, displayed as `image:tag`. Note that tag will not be displayed if its value is `latest`.                      |

## Usage

To use the `journald` driver as the default logging driver, set the `log-driver`
and `log-opt` keys to appropriate values in the `daemon.json` file, which is
located in `/etc/docker/` on Linux hosts or
`C:\ProgramData\docker\config\daemon.json` on Windows Server. For more about
configuring Docker using `daemon.json`, see
[daemon.json](../../../engine/reference/commandline/dockerd.md#daemon-configuration-file).

The following example sets the log driver to `journald`:

```json
{
  "log-driver": "journald"
}
```

Restart Docker for the changes to take effect.

To configure the logging driver for a specific container, use the `--log-driver`
flag on the `docker run` command.

```bash
$ docker run --log-driver=journald ...
```

## Options

Use the `--log-opt NAME=VALUE` flag to specify additional `journald` logging
driver options.

| Option         | Required | Description                                                                                                                                                                   |
|:---------------|:---------|:------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `tag`          | optional | Specify template to set `CONTAINER_TAG` and `SYSLOG_IDENTIFIER` value in journald logs. Refer to [log tag option documentation](log_tags.md) to customize the log tag format. |
| `labels`       | optional | Comma-separated list of keys of labels, which should be included in message, if these labels are specified for the container.                                                 |
| `labels-regex` | optional | Similar to and compatible with labels. A regular expression to match logging-related labels. Used for advanced [log tag options](log_tags.md).                                |
| `env`          | optional | Comma-separated list of keys of environment variables, which should be included in message, if these variables are specified for the container.                               |
| `env-regex`    | optional | Similar to and compatible with env. A regular expression to match logging-related environment variables. Used for advanced [log tag options](log_tags.md).                    |

If a collision occurs between label and env keys, the value of the env takes
precedence. Each option adds additional fields to the attributes of a logging
message.

Below is an example of the logging options required to log to journald.

```bash
$ docker run \
    --log-driver=journald \
    --log-opt labels=location \
    --log-opt env=TEST \
    --env "TEST=false" \
    --label location=west \
    your/application
```

This configuration also directs the driver to include in the payload the label
location, and the environment variable TEST.  If the `--env "TEST=false"`
or `--label location=west` arguments were omitted, the corresponding key would
not be set in the journald log.

## Note regarding container names

The value logged in the `CONTAINER_NAME` field is the name of the container that
was set at startup. If you use `docker rename` to rename a container, the new
name **is not reflected** in the journal entries. Journal entries continue
to use the original name.

## Retrieve log messages with `journalctl`

Use the `journalctl` command to retrieve log messages. You can apply filter
expressions to limit the retrieved messages to those associated with a specific
container:

```bash
$ sudo journalctl CONTAINER_NAME=webserver
```

You can use additional filters to further limit the messages retrieved. The `-b`
flag only retrieves messages generated since the last system boot:

```bash
$ sudo journalctl -b CONTAINER_NAME=webserver
```

The `-o` flag specifies the format for the retried log messages. Use `-o json`
to return the log messages in JSON format.

```bash
$ sudo journalctl -o json CONTAINER_NAME=webserver
```

### View logs for a container with a TTY enabled

If TTY is enabled on a container you may see `[10B blob data]` in the output
when retrieving log messages.
The reason for that is that `\r` is appended to the end of the line and
`journalctl` doesn't strip it automatically unless `--all` is set:

```bash
$ sudo journalctl -b CONTAINER_NAME=webserver --all
```

## Retrieve log messages with the `journal` API

This example uses the `systemd` Python module to retrieve container
logs:

```python
import systemd.journal

reader = systemd.journal.Reader()
reader.add_match('CONTAINER_NAME=web')

for msg in reader:
    print '{CONTAINER_ID_FULL}: {MESSAGE}'.format(**msg)
```
