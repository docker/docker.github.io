---
description: Instructions for installing Docker Engine - Community on Debian
keywords: requirements, apt, installation, debian, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/debian/
- /engine/installation/linux/raspbian/
- /engine/installation/linux/debian/
- /engine/installation/linux/docker-ce/debian/
title: Get Docker Engine - Community for Debian
toc_max: 4
---

To get started with Docker Engine - Community on Debian, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker-ce).

## Prerequisites
 

### OS requirements

To install Docker Engine - Community, you need the 64-bit version of one of these Debian or
Raspbian versions:

- Buster 10
- Stretch 9 (stable) / Raspbian Stretch

Docker Engine - Community is supported on `x86_64` (or `amd64`), `armhf`, and `arm64` architectures.

### Uninstall old versions

Older versions of Docker were called `docker`, `docker.io`, or `docker-engine`.
If these are installed, uninstall them:

```bash
$ sudo apt-get remove docker docker-engine docker.io containerd.io docker-ce runc
```

It's OK if `apt-get` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker Engine - Community package is now called `docker-ce`.

## Install Docker Engine - Community

You can install Docker Engine - Community in different ways, depending on your needs:

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

Before you install Docker Engine - Community for the first time on a new host machine, you need
to set up the Docker repository. Afterward, you can install and update Docker
from the repository.

> **Raspbian users cannot use this method!**
>
> For Raspbian, installing using the repository is not yet supported. You must
> instead use the [convenience script](#install-using-the-convenience-script).

#### Set up the repository

{% assign download-url-base = "https://download.docker.com/linux/debian" %}

1.  Update the `apt` package index:

    ```bash
    $ sudo apt-get update
    ```

2.  Install packages to allow `apt` to use a repository over HTTPS:

    ```bash
    $ sudo apt-get install \
        apt-transport-https \
        ca-certificates \
        curl \
        gnupg2 \
        software-properties-common
    ```

3.  Add Docker's official GPG key:

    ```bash
    $ curl -fsSL {{ download-url-base }}/gpg | sudo apt-key add -
    ```

    Verify that you now have the key with the fingerprint
    `9DC8 5822 9FC7 DD38 854A  E2D8 8D81 803C 0EBF CD88`, by searching for the
    last 8 characters of the fingerprint.

    ```bash
    $ sudo apt-key fingerprint 0EBFCD88

    pub   4096R/0EBFCD88 2017-02-22
          Key fingerprint = 9DC8 5822 9FC7 DD38 854A  E2D8 8D81 803C 0EBF CD88
    uid                  Docker Release (CE deb) <docker@docker.com>
    sub   4096R/F273FCD8 2017-02-22
    ```

4.  Use the following command to set up the **stable** repository. To add the
    **nightly** or **test** repository, add the word `nightly` or `test` (or both)
    after the word `stable` in the commands below. [Learn about **nightly** and **test** channels](/install/index.md).

    > **Note**: The `lsb_release -cs` sub-command below returns the name of your
    > Debian distribution, such as `helium`. Sometimes, in a distribution
    > like BunsenLabs Linux, you might need to change `$(lsb_release -cs)`
    > to your parent Debian distribution. For example, if you are using
    >  `BunsenLabs Linux Helium`, you could use `stretch`. Docker does not offer any guarantees on untested
    > and unsupported Debian distributions.

    <ul class="nav nav-tabs">
      <li class="active"><a data-toggle="tab" data-target="#x86_64_repo">x86_64 / amd64</a></li>
      <li><a data-toggle="tab" data-target="#armhf_repo">armhf</a></li>
      <li><a data-toggle="tab" data-target="#arm64_repo">arm64</a></li>
    </ul>
    <div class="tab-content">
    <div id="x86_64_repo" class="tab-pane fade in active" markdown="1">

    ```bash
    $ sudo add-apt-repository \
       "deb [arch=amd64] {{ download-url-base }} \
       $(lsb_release -cs) \
       stable"
    ```

    </div>
    <div id="armhf_repo" class="tab-pane fade" markdown="1">

    ```bash
    $ sudo add-apt-repository \
       "deb [arch=armhf] {{ download-url-base }} \
       $(lsb_release -cs) \
       stable"
    ```

    </div>
    <div id="arm64_repo" class="tab-pane fade" markdown="1">

    ```bash
    $ sudo add-apt-repository \
       "deb [arch=arm64] {{ download-url-base }} \
       $(lsb_release -cs) \
       stable"
    ```

    </div>
    </div> <!-- tab-content -->

#### Install Docker Engine - Community

> **Note**: This procedure works for Debian on `x86_64` / `amd64`, Debian ARM,
> or Raspbian.

1.  Update the `apt` package index.

    ```bash
    $ sudo apt-get update
    ```

2.  Install the _latest version_ of Docker Engine - Community and containerd, or go to the next step to install a specific version:

    ```bash
    $ sudo apt-get install docker-ce docker-ce-cli containerd.io
    ```

    > Got multiple Docker repositories?
    >
    > If you have multiple Docker repositories enabled, installing
    > or updating without specifying a version in the `apt-get install` or
    > `apt-get update` command always installs the highest possible version,
    > which may not be appropriate for your stability needs.

3.  To install a _specific version_ of Docker Engine - Community, list the available versions in the repo, then select and install:

    a. List the versions available in your repo:

    ```bash
    $ apt-cache madison docker-ce

      docker-ce | 5:18.09.1~3-0~debian-stretch | {{ download-url-base }} stretch/stable amd64 Packages
      docker-ce | 5:18.09.0~3-0~debian-stretch | {{ download-url-base }} stretch/stable amd64 Packages
      docker-ce | 18.06.1~ce~3-0~debian        | {{ download-url-base }} stretch/stable amd64 Packages
      docker-ce | 18.06.0~ce~3-0~debian        | {{ download-url-base }} stretch/stable amd64 Packages
      ...
    ```

    b. Install a specific version using the version string from the second column,
       for example, `5:18.09.1~3-0~debian-stretch `.

    ```bash
    $ sudo apt-get install docker-ce=<VERSION_STRING> docker-ce-cli=<VERSION_STRING> containerd.io
    ```

4.  Verify that Docker Engine - Community is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker Engine - Community is installed and running. The `docker` group is created but no users
are added to it. You need to use `sudo` to run Docker commands.
Continue to [Linux postinstall](/install/linux/linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker Engine - Community

To upgrade Docker Engine - Community, first run `sudo apt-get update`, then follow the
[installation instructions](#install-docker-ce), choosing the new version you want
to install.

### Install from a package

If you cannot use Docker's repository to install Docker Engine - Community, you can download the
`.deb` file for your release and install it manually. You need to download
a new file each time you want to upgrade Docker.

1.  Go to [`{{ download-url-base }}/dists/`]({{ download-url-base }}/dists/){: target="_blank" class="_" },
    choose your Debian version, browse to `pool/stable/`, choose `amd64`,
    `armhf`, or `arm64` and download the `.deb` file for the Docker Engine - Community version
    you want to install.

    > **Note**: To install a **nightly**  package, change the word
    > `stable` in the  URL to `nightly`.
    > [Learn about **nightly** and **test** channels](/install/index.md).

2.  Install Docker Engine - Community, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo dpkg -i /path/to/package.deb
    ```

    The Docker daemon starts automatically.

3.  Verify that Docker Engine - Community is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker Engine - Community is installed and running. The `docker` group is created but no users
are added to it. You need to use `sudo` to run Docker commands.
Continue to [Post-installation steps for Linux](/install/linux/linux-postinstall.md)
to allow non-privileged users to run Docker commands and for other optional
configuration steps.

#### Upgrade Docker Engine - Community

To upgrade Docker Engine - Community, download the newer package file and repeat the
[installation procedure](#install-from-a-package), pointing to the new file.

{% include install-script.md %}

## Uninstall Docker Engine - Community

1.  Uninstall the Docker Engine - Community package:

    ```bash
    $ sudo apt-get purge docker-ce
    ```

2.  Images, containers, volumes, or customized configuration files on your host
    are not automatically removed. To delete all images, containers, and
    volumes:

    ```bash
    $ sudo rm -rf /var/lib/docker
    ```

You must delete any edited configuration files manually.

## Next steps

- Continue to [Post-installation steps for Linux](/install/linux/linux-postinstall.md)

- Continue with the [User Guide](/get-started/index.md).

