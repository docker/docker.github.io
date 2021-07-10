---
description: Instructions for installing Docker Engine on Debian
keywords: requirements, apt, installation, debian, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/debian/
- /engine/installation/linux/raspbian/
- /engine/installation/linux/debian/
- /engine/installation/linux/docker-ce/debian/
- /install/linux/docker-ce/debian/
title: Install Docker Engine on Debian
toc_max: 4
---

To get started with Docker Engine on Debian, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#installation-methods).

## Prerequisites

### OS requirements

To install Docker Engine, you need the 64-bit version of one of these Debian or
Raspbian versions:

- Debian Bullseye 11 (testing)
- Debian Buster 10 (stable)
- Raspbian Bullseye 11 (testing)
- Raspbian Buster 10 (stable)

Docker Engine is supported on `x86_64` (or `amd64`), `armhf`, and `arm64` architectures.

### Uninstall old versions

Older versions of Docker were called `docker`, `docker.io`, or `docker-engine`.
If these are installed, uninstall them:

```console
$ sudo apt-get remove docker docker-engine docker.io containerd runc
```

It's OK if `apt-get` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker Engine package is now called `docker-ce`.

## Installation methods

You can install Docker Engine in different ways, depending on your needs:

- Most users
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach, except for Raspbian.

- Some users download the DEB package and
  [install it manually](#install-from-a-package) and manage
  upgrades completely manually. This is useful in situations such as installing
  Docker on air-gapped systems with no access to the internet.

- In testing and development environments, some users choose to use automated
  [convenience scripts](#install-using-the-convenience-script) to install Docker.
  This is currently the only approach for Raspbian.

### Install using the repository

Before you install Docker Engine for the first time on a new host machine, you need
to set up the Docker repository. Afterward, you can install and update Docker
from the repository.

> **Raspbian users cannot use this method!**
>
> For Raspbian, installing using the repository is not yet supported. You must
> instead use the [convenience script](#install-using-the-convenience-script).
{:.important}

#### Set up the repository

{% assign download-url-base = "https://download.docker.com/linux/debian" %}

1.  Update the `apt` package index and install packages to allow `apt` to use a
    repository over HTTPS:

    ```console
    $ sudo apt-get update

    $ sudo apt-get install \
        apt-transport-https \
        ca-certificates \
        curl \
        gnupg \
        lsb-release
    ```

2.  Add Docker's official GPG key:

    ```console
    $ curl -fsSL {{ download-url-base }}/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    ```

3.  Use the following command to set up the **stable** repository. To add the
    **nightly** or **test** repository, add the word `nightly` or `test` (or both)
    after the word `stable` in the commands below. [Learn about **nightly** and **test** channels](index.md).

    > **Other distributions**
    >
    > The `lsb_release -cs` sub-command below returns the name of your
    > Debian distribution, such as `stretch`. Sometimes, in a distribution
    > like BunsenLabs Linux, you might need to change `$(lsb_release -cs)`
    > to your parent Debian distribution. For example, if you are using
    >  `BunsenLabs Linux Helium`, you could use `stretch`.
    > Docker does not offer any guarantees on untested and unsupported Debian
    > derivatives.
    {:.important}

    <ul class="nav nav-tabs">
      <li class="active"><a data-toggle="tab" data-target="#x86_64_repo">x86_64 / amd64</a></li>
      <li><a data-toggle="tab" data-target="#armhf_repo">armhf</a></li>
      <li><a data-toggle="tab" data-target="#arm64_repo">arm64</a></li>
    </ul>
    <div class="tab-content">
    <div id="x86_64_repo" class="tab-pane fade in active" markdown="1">

    ```console
    $ echo \
      "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] {{ download-url-base }} \
      $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    ```

    </div>
    <div id="armhf_repo" class="tab-pane fade" markdown="1">

    ```console
    $ echo \
      "deb [arch=armhf signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] {{ download-url-base }} \
      $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    ```

    </div>
    <div id="arm64_repo" class="tab-pane fade" markdown="1">

    ```console
    $ echo \
      "deb [arch=arm64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] {{ download-url-base }} \
      $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    ```

    </div>
    </div> <!-- tab-content -->

#### Install Docker Engine

> This procedure works for Debian on `x86_64` / `amd64`, `armhf`, `arm64`, and Raspbian.
{:.note}

1. Update the `apt` package index, and install the _latest version_ of Docker
   Engine and containerd, or go to the next step to install a specific version:

    ```console
    $ sudo apt-get update
    $ sudo apt-get install docker-ce docker-ce-cli containerd.io
    ```

    > Got multiple Docker repositories enabled?
    >
    > If you have multiple Docker repositories enabled, installing
    > or updating without specifying a version in the `apt-get install` or
    > `apt-get update` command always installs the highest possible version,
    > which may not be appropriate for your stability needs.
    {:.tip}

2.  To install a _specific version_ of Docker Engine, list the available versions
    in the repo, then select and install:

    a. List the versions available in your repo:

    ```console
    $ apt-cache madison docker-ce

      docker-ce | 5:18.09.1~3-0~debian-stretch | {{ download-url-base }} stretch/stable amd64 Packages
      docker-ce | 5:18.09.0~3-0~debian-stretch | {{ download-url-base }} stretch/stable amd64 Packages
      docker-ce | 18.06.1~ce~3-0~debian        | {{ download-url-base }} stretch/stable amd64 Packages
      docker-ce | 18.06.0~ce~3-0~debian        | {{ download-url-base }} stretch/stable amd64 Packages
    ```

    b. Install a specific version using the version string from the second column,
       for example, `5:18.09.1~3-0~debian-stretch `.

    ```console
    $ sudo apt-get install docker-ce=<VERSION_STRING> docker-ce-cli=<VERSION_STRING> containerd.io
    ```

3.  Verify that Docker Engine is installed correctly by running the `hello-world`
    image.

    ```console
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker Engine is installed and running. The `docker` group is created but no users
are added to it. You need to use `sudo` to run Docker commands.
Continue to [Linux postinstall](linux-postinstall.md) to allow non-privileged
users to run Docker commands and for other optional configuration steps.

#### Upgrade Docker Engine

To upgrade Docker Engine, first run `sudo apt-get update`, then follow the
[installation instructions](#install-using-the-repository), choosing the new
version you want to install.

### Install from a package

If you cannot use Docker's repository to install Docker Engine, you can download the
`.deb` file for your release and install it manually. You need to download
a new file each time you want to upgrade Docker.

1.  Go to [`{{ download-url-base }}/dists/`]({{ download-url-base }}/dists/){: target="_blank" rel="noopener" class="_" },
    choose your Debian version, then browse to `pool/stable/`, choose `amd64`,
    `armhf`, or `arm64`, and download the `.deb` file for the Docker Engine
    version you want to install.

    > Feeling adventurous?
    >
    > To install a **test** (pre-release) or **nightly** package,
    > change the word `stable` in the above URL to `test` or `nightly`.
    > [Learn about **test** and **nightly** channels](index.md).
    {:.tip}

2.  Install Docker Engine, changing the path below to the path where you downloaded
    the Docker package.

    ```console
    $ sudo dpkg -i /path/to/package.deb
    ```

    The Docker daemon starts automatically.

3.  Verify that Docker Engine is installed correctly by running the `hello-world`
    image.

    ```console
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker Engine is installed and running. The `docker` group is created but no users
are added to it. You need to use `sudo` to run Docker commands.
Continue to [Post-installation steps for Linux](linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker Engine

To upgrade Docker Engine, download the newer package file and repeat the
[installation procedure](#install-from-a-package), pointing to the new file.

{% include install-script.md %}

## Uninstall Docker Engine

1.  Uninstall the Docker Engine, CLI, and Containerd packages:

    ```console
    $ sudo apt-get purge docker-ce docker-ce-cli containerd.io
    ```

2.  Images, containers, volumes, or customized configuration files on your host
    are not automatically removed. To delete all images, containers, and
    volumes:

    ```console
    $ sudo rm -rf /var/lib/docker
    $ sudo rm -rf /var/lib/containerd
    ```

You must delete any edited configuration files manually.

## Next steps

- Continue to [Post-installation steps for Linux](linux-postinstall.md).
- Review the topics in [Develop with Docker](../../develop/index.md) to learn how to build new applications using Docker.
