---
description: Home page for Docker's documentation
keywords: Docker, documentation, manual, guide, reference, api, samples
landing: true
title: Docker Documentation
notoc: true
notags: true
---
{% assign page.title = site.name %}

## Docs Hackathon, April 17-22nd, 2017

<a href="/hackathon/"><img src="docs-hackathon-2.png" alt="Docker Docs Hackathon, April 17-22nd, 2017" style="max-width: 100%"></a>

Fix docs bugs to claim the points, and cash in your points for prizes in [the swag store](http://www.cafepress.com/dockerdocshackathon). Every 10 points is worth $1 USD in store credit. Happening all DockerCon week, from April 17-21, 2017.

[Hackathon details](/hackathon/){: class="button outline-btn" style="margin:20px"}[View available bugs on GitHub](https://github.com/docker/docker.github.io/milestone/9){: class="button outline-btn" style="margin:20px"} [Visit the rewards store](http://www.cafepress.com/dockerdocshackathon){: class="button outline-btn" style="margin:20px"}

## Introduction to Docker


Docker packages your app with its dependencies, freeing you from worrying about
your system configuration, and making your app more portable.

{% if site.edge == true %}
{% capture ce-edge-section %}

## Docker CE Edge

The Docker CE Edge channel provides monthly releases which allow you to try
new features of Docker and verify bug fixes quickly. Edge releases are only
supported for one month, and a given Edge release will not receive any updates
once a new edge release is available.

Stable releases are not published to the Edge channel, so Linux repository users
still need to subscribe to the Stable channel as well.

Commercial support is not available for Docker CE.

For information about all Docker release channels and expectations about
support, see [Docker channels](/engine/installation/#docker-channels).

<!-- This button toggles the div below, and hides itself when clicked -->
<a id="ce-edge-readmore-btn" onclick="$(this).hide(); ga('send', 'event', 'ce-edge-readmore', 'click', 'CE engagement');" data-target="#ce-edge-readmore" data-toggle="collapse" class="button outline-btn collapse in">Read more about Docker CE Edge releases</a>

<div markdown="1" id="ce-edge-readmore" class="collapse" data-target="#ce-edge-readmore-btn" data-toggle="collapse">

This page lists features that are only available in Docker CE Edge releases.
Where applicable, the API and CLI reference documentation has been updated to
reflect these features, but **full documentation for a given feature may not be
available until a Docker CE Stable release incorporates the feature**.

### Docker CE Edge new features

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#1704">17.04</a></li>
  <!--<li><a data-toggle="tab" data-target="#1705">17.05</a></li>-->
</ul>
<div markdown="1" class="tab-content">
<div markdown="1" id="1704" class="tab-pane fade in active">

#### Docker CE Edge 17.04

The following major features and changes are included in Docker CE Edge 17.04.
Continue reading, or go straight to [API and CLI](#api-and-cli),
[Daemon](#daemon), [Dockerfile](#dockerfile), [Services](#services), or
[Stacks](#stacks).

[Read the full release notes](https://github.com/docker/docker/releases/tag/v17.04.0-ce){: target="_blank" class="_" }

##### API and CLI

- Add `--device-cgroup-rule` flag to give containers access to devices that appear
  after the container is started. {% include github-pr.md pr=22563 %}

- Allow swarm nodes to join with `--availability=drain` to prevent them from
  taking non-manager workloads. {% include github-pr.md pr=24993 %}

- Add `publish` and `expose` filters to `docker ps`, so that containers can be
  filtered by port or port range for TCP or UDP protocols {% include github-pr.md pr=27557 %}

- Add `--no-trunc` and `--format` flags to the `docker service ls` command, and
  as well as the ability to specify the default format for `docker service ls`
  using the `ServicesFormat` option to the Docker CLI. Also add a
  `docker stack services` command. {% include github-pr.md pr=28199 %}

- Add ability to filter plugins by whether they are enabled or disabled in
  `docker plugin ls` output. {% include github-pr.md pr=28627 %}

- Add `mode` option to `--log-opts` flag for both `docker` and `dockerd`. If set
  to `non-blocking`, and the log buffer fills up, log messages will be lost, but
  the container will not block. The `max-buffer-size` option controls the
  maximum size of the ring buffer. Defaults to `blocking`, which will cause the
  container to block if messages cannot be logged. See
  [Options for all drivers](/engine/admin/logging/overview.md#options-for-all-drivers).
  {% include github-pr.md pr=28762 %}

- It is no longer possible to inadvertently pull images on an architecture where
  they will not run. {% include github-pr.md pr=29001 %}

- It is now possible to create AWS log groups when using the AWS logging driver.
  See [`awslogs-create-group`](engine/admin/logging/awslogs.md#awslogs-create-group).
  {% include github-pr.md pr=29504 %}

- Add the ability to filter `docker network ls` output by creation time, using
  the `{% raw %}{{CreatedAt}}{% endraw %}` format specifier.
  {% include github-pr.md pr=29900 %}

- Named but untagged images are now removed if you run `docker image prune` if
  `--dangling-only` is set to `true`. {% include github-pr.md pr=30330 %}

- Add `--add-host` flag to `docker build`, which will add entries to the
  `/etc/hosts` file of a container created from that image. The `/etc/hosts`
  file is not saved within the image itself. {% include github-pr.md pr=30383 %}

- Prevent `docker network ls` from pulling all the endpoints, to reduce
  impact on the network. {% include github-pr.md pr=30673 %}

- Windows-specific commands and options no longer show in command help text on
  non-Windows clients. {% include github-pr.md pr=30780 %}

- When you specify an IP address when running `docker network connect`, the
  IP address is now checked for validity. {% include github-pr.md pr=30807 %}

- Add the ability to customize bind-mount consistency to be more appropriate
  for some platforms and workloads. Options are `consistent` (the default),
  `cached`, or `delegated`. {% include github-pr.md pr=31047 %}

##### Daemon

- Docker Daemon logging settings no longer affect the `docker build` command.
  {% include github-pr.md pr=29552 %}

- Add a `registry-mirrors` configuration option for the Docker daemon, which
  replaces the daemon's registry mirrors with a new set of registry mirrors.
  {% include github-pr.md pr=29650 %}

- Add the ability to specify the default shared memory size for the Docker
  daemon, using the `--default-shm-size` or the `default-shm-size` key in
  `daemon.json`. {% include github-pr.md pr=29692 %}

- Add a `no-new-privileges` configuration option for the Docker daemon, which
  prevents unprivileged containers from gaining new privileges.
  {% include github-pr.md pr=29984 %}

- If a Docker client communicates with an older daemon and attempts to perform
  an operation not supported by the daemon, an error is printed, which shows
  the API versions of both the client and daemon.
  {% include github-pr.md pr=30187 %}

- The Docker daemon no longer depends upon `sqlite`. This change means that it
  is not possible to upgrade the Docker daemon from version 1.9 directly to the
  latest version. It is recommended to upgrade from one major version to the
  next, in sequence. {% include github-pr.md pr=30208 %}

##### Dockerfile

- Using the pattern `**/` in a Dockerfile now (correctly) behaves the same as
  `**`. {% include github-pr.md pr=29043 %}

- Time values less than 1 second are no longer allowed in health-check options
  in the Dockerfile. {% include github-pr.md pr=31177 %}

##### Services

- When a service is updated with both `--secret-add` and `--secret-rm` in the
  same operation, the order of operations is now changed so that the
  `--secret-rm` always occurs first. {% include github-pr.md pr=29802 %}

- Add the ability to create or update a service to be read-only using the
  `--read-only` flag. {% include github-pr.md pr=30162 %}

- Docker now updates swarm nodes if the swarm configuration is updated.
  {% include github-pr.md pr=30259 %}

- Add topology-aware placement preferences for Swarm services. This feature
  allows services to be balanced over nodes based on a particular user-defined
  property, such as which datacenter or rack they are located in.
  See [Control service scale and placement](/engine/swarm/services.md#control-service-scale-and-placement).
  {% include github-pr.md pr=30725 %}

- Add the ability to customize the stop signal which will be sent to nodes, when
  creating or updating a service. {% include github-pr.md pr=30754 %}

- Add the ability to address a secret by name or prefix, as well as ID, when
  updating it. {% include github-pr.md pr=30856 %}

- Add the ability to roll back to a previous version of a service if an
  updated service fails to deploy. Several flags are available at service
  creation or update,to control the rollback action, failure threshold,
  monitoring delay, rollback delay, and parallelism.
  {% include github-pr.md pr=31108 %}

- Add the ability to specify the stream when using the Docker service logs API.
  {% include github-pr.md pr=31313 %}

- Add `--tail` and `--since` flags to `docker service logs` command, to filter
  the logs by time or to show the tail of the logs and show new content as it
  is logged. {% include github-pr.md pr=31500 %}

- Add a `--verbose` flag to the `docker inspect` command. For swarm networks,
  this flag shows all nodes and services attached to the network.
  {% include github-pr.md pr=31710 %}

##### Stacks

- Compose file version 3.2 is now supported. This includes support for different
  types of endpoints and expands the options you can use when specifying mounts.
  {% include github-pr.md pr=31795 %}

</div> <!-- 17.04 -->
<!--<div id="1705" class="tab-pane fade">TAB 2 CONTENT</div>-->
</div> <!-- tab-content -->
</div> <!-- ce-edge-readmore -->
{% endcapture %} <!-- from line 13 -->
{{ ce-edge-section | markdownify }}
{% endif %}

<div class="row">
<div markdown="1" class="col-xs-12 col-sm-12 col-md-12 col-lg-6 block">
### Learn Docker basics

Get started learning Docker concepts, tools, and commands. The examples show you
how to build, push, and pull Docker images, and run them as containers. This
tutorial stops short of teaching you how to deploy applications.

[Start the basic tutorial](/engine/getstarted/){: class="button outline-btn"}
</div>

<div markdown="1" class="col-xs-12 col-sm-12 col-md-12 col-lg-6 block">
### Define and deploy apps in Swarm Mode

Learn how to relate containers to each other, define them as services, and
configure an application stack ready to deploy at scale in a production
environment. Highlights Compose Version 3 new features and swarm mode.

[Start the application tutorial](/engine/getstarted-voting-app/){: class="button outline-btn"}
</div>
</div>

## Components

<div class="component-container">
    <!--start row-->
    <div class="row">
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="docker-for-mac/"> <img src="../images/apple_48.svg" alt="Docker for Mac"> </a>
                </div>
                <h3 id="docker-for-mac"><a href="docker-for-mac/">Docker for Mac</a></h3>
                <p>A native application using the macOS sandbox security model which delivers all Docker tools to your Mac.</p>
            </div>
        </div>
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="docker-for-windows/"> <img src="../images/windows_48.svg" alt="Docker for Windows"> </a>
                </div>
                <h3 id="docker-for-windows"><a href="docker-for-windows/">Docker for Windows</a></h3>
                <p>A native Windows application which delivers all Docker tools to your Windows computer.</p>
            </div>
        </div>
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="engine/installation/linux/ubuntu/"> <img src="../images/linux_48.svg" alt="Docker for Linux"> </a>
                </div>
                <h3 id="docker-for-linux"><a href="engine/installation/linux/ubuntu/">Docker for Linux</a></h3>
                <p>Install Docker on a computer which already has a Linux distribution installed.</p>
            </div>
        </div>
    <!--organic row 2-->
    <!--editions-->
    <div class="col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                <a href="engine/installation/"> <img src="../images/apple_48.svg" alt="Docker Editions"> </a>
            </div>
            <h3 id="editions"><a href="engine/installation/">Docker Editions</a></h3>
            <p>Platform matrix and superset of installers for Docker for desktops, servers, or cloud providers.</p>
        </div>
    </div>
    <!--compose-->
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="compose/overview/"> <img src="../images/compose_48.svg" alt="Docker Compose"> </a>
                </div>
                <h3 id="compose"><a href="compose/overview/">Docker Compose</a></h3>
                <p>Define application stacks built using multiple containers, services, and swarm configurations.</p>
        </div>
    </div>
    <!--machine-->
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="machine/overview/"> <img src="../images/machine_48.svg" alt="Docker Trusted Registry"> </a>
                </div>
                <h3 id="machine"><a href="machine/overview/">Docker Machine</a></h3>
                <p>Automate container provisioning on your network or in the cloud. Available for Windows, macOS, or Linux.</p>
        </div>
    </div>
    <!--organic row 3-->
    <!--cloud-->
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="docker-cloud/"> <img src="../images/cloud_48.svg" alt="Docker Cloud"> </a>
                </div>
                <h3 id="docker-cloud"><a href="docker-cloud/">Docker Cloud</a></h3>
                <p>A hosted service for building, testing, and deploying Docker images to your hosts.</p>
            </div>
        </div>
    <!--UCP-->
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="datacenter/ucp/2.1/guides/"> <img src="../images/UCP_48.svg" alt="Docker Universal Control Plane"> </a>
                </div>
                <h3 id="ucp"><a href="datacenter/ucp/2.1/guides/">Docker Universal Control Plane</a></h3>
                <p>(UCP) Manage a cluster of on-premise Docker hosts like a single machine with this enterprise product.</p>
            </div>
        </div>
    <!--DTR-->
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="datacenter/dtr/2.2/guides/"> <img src="../images/dtr_48.svg" alt="Docker Trusted Registry"> </a>
                </div>
                <h3 id="dtr"><a href="datacenter/dtr/2.2/guides/">Docker Trusted Registry</a></h3>
                <p>(DTR) An enterprise image storage solution you can install behind a firewall to manage images and access.</p>
        </div>
    </div>
    <!-- end real row-->
    </div>
<!-- end component-container 2-->
</div>
