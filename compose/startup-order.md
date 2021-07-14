---
description: How to control service startup and shutdown order in Docker Compose
keywords: documentation, docs, docker, compose, startup, shutdown, order
title: Control startup and shutdown order in Compose
notoc: true
---

You can control the order of service startup and shutdown with the
[depends_on](compose-file/compose-file-v3.md#depends_on) option. Compose always starts and stops
containers in dependency order, where dependencies are determined by
`depends_on`, `links`, `volumes_from`, and `network_mode: "service:..."`.

However, for startup Compose does not wait until a container is "ready" (whatever that means
for your particular application) - only until it's running. There's a good
reason for this.

The problem of waiting for a database (for example) to be ready is really just
a subset of a much larger problem of distributed systems. In production, your
database could become unavailable or move hosts at any time. Your application
needs to be resilient to these types of failures.

To handle this, design your application to attempt to re-establish a connection to
the database after a failure. If the application retries the connection,
it can eventually connect to the database.

The best solution is to perform this check in your application code, both at
startup and whenever a connection is lost for any reason. However, if you don't
need this level of resilience, you can work around the problem with a wrapper
script:

- Use a tool such as [wait-for-it](https://github.com/vishnubob/wait-for-it),
  [dockerize](https://github.com/jwilder/dockerize), sh-compatible
  [wait-for](https://github.com/Eficode/wait-for), or [RelayAndContainers](https://github.com/jasonsychau/RelayAndContainers) template. These are small
  wrapper scripts which you can include in your application's image to
  poll a given host and port until it's accepting TCP connections.

  For example, to use `wait-for-it.sh` or `wait-for` to wrap your service's command:

  ```yaml
  version: "2"
  services:
    web:
      build: .
      ports:
        - "80:8000"
      depends_on:
        - "db"
      command: ["./wait-for-it.sh", "db:5432", "--", "python", "app.py"]
    db:
      image: postgres
  ```

  > **Tip**
  >
  > There are limitations to this first solution. For example, it doesn't verify
  > when a specific service is really ready. If you add more arguments to the
  > command, use the `bash shift` command with a loop, as shown in the next
  > example.

- Alternatively, write your own wrapper script to perform a more application-specific
  health check. For example, you might want to wait until Postgres is ready to
  accept commands:

  ```bash
  #!/bin/sh
  # wait-for-postgres.sh

  set -e
  
  host="$1"
  # Shift arguments with mapping:
  # - $0 => $0
  # - $1 => <discarded>
  # - $2 => $1
  # - $3 => $2
  # - ...
  # This is done for `exec "$@"` below to work correctly
  shift
  
  # Login for user (`-U`) and once logged in execute quit ( `-c \q` )
  # If we can not login sleep for 1 sec
  until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$host" -U "postgres" -c '\q'; do
    >&2 echo "Postgres is unavailable - sleeping"
    sleep 1
  done
  
  >&2 echo "Postgres is up - executing command"
  # Print and execute all other arguments starting with `$1`
  # So `exec "$1" "$2" "$3" ...`
  exec "$@"
  ```

  You can use this as a wrapper script as in the previous example, by setting:

  ```yaml
  command: ["./wait-for-postgres.sh", "db", "python", "app.py"]
  ```


## Compose documentation

- [User guide](index.md)
- [Installing Compose](install.md)
- [Getting Started](gettingstarted.md)
- [Command line reference](reference/index.md)
- [Compose file reference](compose-file/index.md)
- [Sample apps with Compose](samples-for-compose.md)
