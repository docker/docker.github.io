---
description: Describes how to format tags for.
keywords: docker, logging, driver, syslog, Fluentd, gelf, journald
redirect_from:
- /engine/reference/logging/log_tags/
- /engine/admin/logging/log_tags/
title: Customize log driver output
---

The `tag` log option specifies how to format a tag that identifies the
container's log messages. By default, the system uses the first 12 characters of
the container ID. To override this behavior, specify a `tag` option:

```bash
$ docker run --log-driver=fluentd --log-opt fluentd-address=myhost.local:24224 --log-opt tag="mailer"
```

Docker supports some special template markup you can use when specifying a tag's value:

{% raw %}
| Markup             | Description                                          |
|--------------------|------------------------------------------------------|
| `{{.ID}}`          | The first 12 characters of the container ID.         |
| `{{.FullID}}`      | The full container ID.                               |
| `{{.Name}}`        | The container name.                                  |
| `{{.ImageID}}`     | The first 12 characters of the container's image ID. |
| `{{.ImageFullID}}` | The container's full image ID.                       |
| `{{.ImageName}}`   | The name of the image used by the container.         |
| `{{.DaemonName}}`  | The name of the docker program (`docker`).           |
| `{{.Hostname}}`    | Hostname from the underlying OS                      |

{% endraw %}

For example, specifying a {% raw %}`--log-opt tag="{{.ImageName}}/{{.Name}}/{{.ID}}"`{% endraw %} value yields `syslog` log lines like:

```none
Aug  7 18:33:19 HOSTNAME hello-world/foobar/5790672ab6a0[9103]: Hello from Docker.
```

At startup time, the system sets the `container_name` field and {% raw %}`{{.Name}}`{% endraw %} in
the tags. If you use `docker rename` to rename a container, the new name is not
reflected in the log messages. Instead, these messages continue to use the
original container name.

