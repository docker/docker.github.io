---
description: Instructions for installing Docker Engine - Community on Ubuntu
keywords: requirements, apt, installation, ubuntu, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/ubuntulinux/
- /installation/ubuntulinux/
- /engine/installation/linux/ubuntulinux/
- /engine/installation/linux/docker-ce/ubuntu/
- /install/linux/ubuntu/
- /engine/installation/linux/ubuntu/
title: Get Docker Engine - Community for Ubuntu
toc_max: 4
---

To get started with Docker Engine - Community on Ubuntu, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker-engine---community-1).

## Prerequisites

### Docker EE customers

To install Docker Enterprise Edition (Docker EE), go to
[Get Docker EE for Ubuntu](/install/linux/docker-ee/ubuntu.md)
**instead of this topic**.

To learn more about Docker EE, see
[Docker Enterprise Edition](https://www.docker.com/enterprise-edition/){: target="_blank" class="_" }.

### OS requirements

To install Docker Engine - Community, you need the 64-bit version of one of these Ubuntu
versions:

- Eoan 19.10
- Bionic 18.04 (LTS)
- Xenial 16.04 (LTS)

Docker Engine - Community is supported on `x86_64` (or `amd64`), `armhf`, `arm64`, `s390x`
(IBM Z), and `ppc64le` (IBM Power) architectures.

### Intermediate (non-LTS) Ubuntu versions

Ubuntu intermediate (non-LTS) releases are [supported for 9 months](https://ubuntu.com/about/release-cycle){: target="_blank" class="_" },
after which they reach end-of-life (EOL). The package repositories at download.docker.com
may contain packages for older Ubuntu releases that reached EOL. These packages
may be outdated, and do not receive updates (including security updates).

For production use, we recommend installing the latest Ubuntu LTS version. If you
are running a non-LTS version of Ubuntu, make sure to upgrade to the latest supported
version. Be aware that Docker's release cycle may not align with Ubuntu releases,
which means that packages for the latest Ubuntu release may not be available at
the day of release. Before installing or upgrading, verify that Docker packages
are available for your version of Ubuntu.

When upgrading your Ubuntu version on a machine that already has Docker installed,
the `apt` repository may need to be updated to match the upgraded Ubuntu version.
Refer to the [update the repository after upgrading Ubuntu versions](#update-the-repository-after-upgrading-ubuntu-versions)
section below for details.

### Uninstall old versions

Older versions of Docker were called `docker`, `docker.io`, or `docker-engine`.
If these are installed, uninstall them:

```bash
$ sudo apt-get remove docker docker-engine docker.io containerd runc
```

It's OK if `apt-get` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker Engine - Community package is now called `docker-ce`.

### Supported storage drivers

Docker Engine - Community on Ubuntu supports `overlay2`, `aufs` and `btrfs` storage drivers.
> **Note**: In Docker Engine - Enterprise, `btrfs` is only supported on SLES. See the documentation on
> [btrfs](/engine/userguide/storagedriver/btrfs-driver.md) for more details.

For new installations on version 4 and higher of the Linux kernel, `overlay2`
is supported and preferred over `aufs`. Docker Engine - Community uses the `overlay2`
storage driver by default. If you need to use `aufs` instead, you need to
configure it manually. See [aufs](/engine/userguide/storagedriver/aufs-driver.md)

## Install Docker Engine - Community

You can install Docker Engine - Community in different ways, depending on your needs:

- Most users
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach.

- Some users download the DEB package and
  [install it manually](#install-from-a-package) and manage
  upgrades completely manually. This is useful in situations such as installing
  Docker on air-gapped systems with no access to the internet.

- In testing and development environments, some users choose to use automated
  [convenience scripts](#install-using-the-convenience-script) to install Docker.

### Install using the repository

Before you install Docker Engine - Community for the first time on a new host machine, you need
to set up the Docker repository. Afterward, you can install and update Docker
from the repository.

#### Set up the repository

{% assign download-url-base = "https://download.docker.com/linux/ubuntu" %}

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
        gnupg-agent \
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
    
    pub   rsa4096 2017-02-22 [SCEA]
          9DC8 5822 9FC7 DD38 854A  E2D8 8D81 803C 0EBF CD88
    uid           [ unknown] Docker Release (CE deb) <docker@docker.com>
    sub   rsa4096 2017-02-22 [S]
    ```

4.  Use the following command to set up the **stable** repository. To add the
    **nightly** or **test** repository, add the word `nightly` or `test` (or both)
    after the word `stable` in the commands below. [Learn about **nightly** and **test** channels](/install/index.md).

    > **Note**: The `lsb_release -cs` sub-command below returns the name of your
    > Ubuntu distribution, such as `xenial`. Sometimes, in a distribution
    > like Linux Mint, you might need to change `$(lsb_release -cs)`
    > to your parent Ubuntu distribution. For example, if you are using
    >  `Linux Mint Tessa`, you could use `bionic`. Docker does not offer any guarantees on untested
    > and unsupported Ubuntu distributions.


    <ul class="nav nav-tabs">
      <li class="active"><a data-toggle="tab" data-target="#x86_64_repo">x86_64 / amd64</a></li>
      <li><a data-toggle="tab" data-target="#armhf_repo">armhf</a></li>
      <li><a data-toggle="tab" data-target="#arm64_repo">arm64</a></li>
      <li><a data-toggle="tab" data-target="#ppc64le_repo">ppc64le (IBM Power)</a></li>
      <li><a data-toggle="tab" data-target="#s390x_repo">s390x (IBM Z)</a></li>
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
    <div id="ppc64le_repo" class="tab-pane fade" markdown="1">

    ```bash
    $ sudo add-apt-repository \
       "deb [arch=ppc64el] {{ download-url-base }} \
       $(lsb_release -cs) \
       stable"
    ```

    </div>
    <div id="s390x_repo" class="tab-pane fade" markdown="1">

    ```bash
    $ sudo add-apt-repository \
       "deb [arch=s390x] {{ download-url-base }} \
       $(lsb_release -cs) \
       stable"
    ```

    </div>
    </div> <!-- tab-content -->

#### Install Docker Engine - Community

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

      docker-ce | 5:18.09.1~3-0~ubuntu-xenial | {{ download-url-base }}  xenial/stable amd64 Packages
      docker-ce | 5:18.09.0~3-0~ubuntu-xenial | {{ download-url-base }}  xenial/stable amd64 Packages
      docker-ce | 18.06.1~ce~3-0~ubuntu       | {{ download-url-base }}  xenial/stable amd64 Packages
      docker-ce | 18.06.0~ce~3-0~ubuntu       | {{ download-url-base }}  xenial/stable amd64 Packages
      ...
    ```

    b. Install a specific version using the version string from the second column,
       for example, `5:18.09.1~3-0~ubuntu-xenial`.

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
    choose your Ubuntu version, browse to `pool/stable/`, choose `amd64`,
    `armhf`, `arm64`, `ppc64el`, or `s390x`, and download the `.deb` file for the
    Docker Engine - Community version you want to install.

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

### Update the repository after upgrading Ubuntu versions

When upgrading your version of Ubuntu, the URL of the package repository may have
to be updated to match the Ubuntu version you have installed. Use the `apt-cache policy`
command to find the repository that is currently configured. The example below
shows that `apt` is configured to install packages for Ubuntu `disco` (19.04):

```bash
apt-cache policy | grep docker

 500 https://download.docker.com/linux/ubuntu disco/stable amd64 Packages
     origin download.docker.com
```

If the Ubuntu version in the output does not match the version you have installed,
the repository URL needs to be updated. To update the package repository to match
your current version, either removing all lines containing `download.docker.com`
in `/etc/apt/sources.list`, or use the `add-apt-repository --remove` command:

```bash
sudo add-apt-repository --remove "deb {{ download-url-base }} disco stable"
```

Verify that the repository was removed:

```bash
apt-cache policy | grep docker
```

After removing the old repository, continue with the steps outlined in [set up the
repository](#set-up-the-repository) to configure the correct package repository
for your Ubuntu version and to update or install the latest Docker release.


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

- Continue to [Post-installation steps for Linux](/install/linux/linux-postinstall.md).
- Review the topics in [Develop with Docker](/develop/index.md) to learn how to build new applications using Docker.
