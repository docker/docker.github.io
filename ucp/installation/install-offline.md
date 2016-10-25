---
description: Learn how to install Docker Universal Control Plane. on a machine with
  no internet access.
keywords:
- docker, ucp, install, offline
menu:
  main:
    identifier: ucp_install_offline
    parent: mn_ucp_installation
    weight: 30
title: Install offline
---

# Install UCP offline

The procedure to install Docker Universal Control Plane on a host is the same,
whether that host has access to the internet or not.

The only difference when installing on an offline host,
is that instead of pulling the UCP images from Docker Hub, you use a
computer that is connected to the internet to download a single package with
all the images. Then you copy that package to the host where you’ll install UCP.


1.  Get the UCP package.

    Use a computer with internet access to download a single package with all
    Docker Datacenter components:

    ```none
    $ wget https://packages.docker.com/caas/ucp-1.1.3_dtr-2.0.3.tar.gz -O docker-datacenter.tar.gz
    ```

2.  Transfer the package to the offline node.

    Now that you have the UCP package in your machine, you can transfer it to the
    host where you'll be installing Docker UCP. You can use the Secure Copy command
    for this:

    ```none
    $ scp docker-datacenter.tar.gz $USER@$UCP_HOST:/tmp
    ```

3.  Login into the host where you transferred the images.

4.  Load the UCP images.

    Once the UCP package is transferred to the host, you can use the
    `docker load` command, to load the images from the tar archive. On the host
    were you are going to install UCP, run:

    ```none
    $ docker load < docker-datacenter.tar.gz
    ```

5.  Install Docker UCP.

    Now that the offline host has all the images needed to install UCP,
    you can [install Docker UCP that host](install-production.md).


## Where to go next

* [Install UCP for production](install-production.md).
* [UCP system requirements](system-requirements.md)
