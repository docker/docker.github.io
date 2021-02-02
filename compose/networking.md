---
description: How Compose sets up networking between containers
keywords: documentation, docs, docker, compose, orchestration, containers, networking
title: Networking in Compose
---

> This page applies to Compose file formats [version 2](compose-file/compose-file-v2.md) and [higher](compose-file/index.md). Networking features are not supported for Compose file [version 1 (legacy)](compose-file/compose-file-v1.md).

By default Compose sets up a single
[network](../engine/reference/commandline/network_create.md) for your app. Each
container for a service joins the default network and is both *reachable* by
other containers on that network, and *discoverable* by them at a hostname
identical to the container name.

> **Note**
>
> Your app's network is given a name based on the "project name",
> which is based on the name of the directory it lives in. You can override the
> project name with either the [`--project-name` flag](reference/overview.md)
> or the [`COMPOSE_PROJECT_NAME` environment variable](reference/envvars.md#compose_project_name).

For example, suppose your app is in a directory called `myapp`, and your `docker-compose.yml` looks like this:

```yaml
version: "{{ site.compose_file_v3 }}"
services:
  web:
    build: .
    ports:
      - "8000:8000"
  db:
    image: postgres
    ports:
      - "8001:5432"
```

When you run `docker-compose up`, the following happens:

1.  A network called `myapp_default` is created.
2.  A container is created using `web`'s configuration. It joins the network
    `myapp_default` under the name `web`.
3.  A container is created using `db`'s configuration. It joins the network
    `myapp_default` under the name `db`.

> **In v2.1+, overlay networks are always `attachable`**
>
> Starting in Compose file format 2.1, overlay networks are always created as
> `attachable`, and this is not configurable. This means that standalone
> containers can connect to overlay networks.
>
> In Compose file format 3.x, you can optionally set the `attachable` property
> to `false`.

Each container can now look up the hostname `web` or `db` and
get back the appropriate container's IP address. For example, `web`'s
application code could connect to the URL `postgres://db:5432` and start
using the Postgres database.

It is important to note the distinction between `HOST_PORT` and `CONTAINER_PORT`.
In the above example, for `db`, the `HOST_PORT` is `8001` and the container port is
`5432` (postgres default). Networked service-to-service
communication uses the `CONTAINER_PORT`. When `HOST_PORT` is defined,
the service is accessible outside the swarm as well.

Within the `web` container, your connection string to `db` would look like
`postgres://db:5432`, and from the host machine, the connection string would
look like `postgres://{DOCKER_IP}:8001`.

## Update containers

If you make a configuration change to a service and run `docker-compose up` to update it, the old container is removed and the new one joins the network under a different IP address but the same name. Running containers can look up that name and connect to the new address, but the old address stops working.

If any containers have connections open to the old container, they are closed. It is a container's responsibility to detect this condition, look up the name again and reconnect.

## Links

Links allow you to define extra aliases by which a service is reachable from another service. They are not required to enable services to communicate - by default, any service can reach any other service at that service's name. In the following example, `db` is reachable from `web` at the hostnames `db` and `database`:

    version: "3"
    services:

      web:
        build: .
        links:
          - "db:database"
      db:
        image: postgres

See the [links reference](compose-file/compose-file-v2.md#links) for more information.

## Multi-host networking

When deploying a Compose application on an Docker Engine with [Swarm mode enabled](../engine/swarm/index.md),
you can make use of the built-in `overlay` driver to enable multi-host communication.

Consult the [Swarm mode section](../engine/swarm/index.md), to see how to set up
a Swarm cluster, and the [Getting started with multi-host networking](../network/network-tutorial-overlay.md)
to learn about multi-host overlay networks.

## Specify custom networks

Instead of just using the default app network, you can specify your own networks with the top-level `networks` key. This lets you create more complex topologies and specify [custom network drivers](/engine/extend/plugins_network/) and options. You can also use it to connect services to externally-created networks which aren't managed by Compose.

Each service can specify what networks to connect to with the *service-level* `networks` key, which is a list of names referencing entries under the *top-level* `networks` key.

Here's an example Compose file defining two custom networks. The `proxy` service is isolated from the `db` service, because they do not share a network in common - only `app` can talk to both.

    version: "3"
    services:

      proxy:
        build: ./proxy
        networks:
          - frontend
      app:
        build: ./app
        networks:
          - frontend
          - backend
      db:
        image: postgres
        networks:
          - backend

    networks:
      frontend:
        # Use a custom driver
        driver: custom-driver-1
      backend:
        # Use a custom driver which takes special options
        driver: custom-driver-2
        driver_opts:
          foo: "1"
          bar: "2"

Networks can be configured with static IP addresses by setting the [ipv4_address and/or ipv6_address](compose-file/compose-file-v2.md#ipv4_address-ipv6_address) for each attached network.

Networks can also be given a [custom name](compose-file/compose-file-v3.md#network-configuration-reference) (since version 3.5):

    version: "3.5"
    networks:
      frontend:
        name: custom_frontend
        driver: custom-driver-1

For full details of the network configuration options available, see the following references:

- [Top-level `networks` key](compose-file/compose-file-v2.md#network-configuration-reference)
- [Service-level `networks` key](compose-file/compose-file-v2.md#networks)

## Configure the default network

Instead of (or as well as) specifying your own networks, you can also change the settings of the app-wide default network by defining an entry under `networks` named `default`:

    version: "3"
    services:

      web:
        build: .
        ports:
          - "8000:8000"
      db:
        image: postgres

    networks:
      default:
        # Use a custom driver
        driver: custom-driver-1

## Use a pre-existing network

If you want your containers to join a pre-existing network, use the [`external` option](compose-file/compose-file-v2.md#network-configuration-reference):

    networks:
      my-pre-existing-network:
          external: true

Instead of attempting to create a network called `[projectname]_default`, Compose looks for a network called `my-pre-existing-network` and connect your app's containers to it.
