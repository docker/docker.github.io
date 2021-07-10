---
description: Instructions for installing Docker as a binary. Mostly meant for hackers who want to try out Docker on a variety of environments.
keywords: binaries, installation, docker, documentation, linux
title: Install Docker Engine from binaries
redirect_from:
- /engine/installation/binaries/
- /engine/installation/linux/docker-ce/binaries/
- /install/linux/docker-ce/binaries/
- /installation/binaries/
---

If you want to try Docker or use it in a testing environment, but you're not on
a supported platform, you can try installing from static binaries. If possible,
you should use packages built for your operating system, and use your operating
system's package management system to manage Docker installation and upgrades.
Be aware that 32-bit static binary archives do not include the Docker daemon.

Static binaries for the Docker daemon binary are only available for Linux (as
`dockerd`). 
Static binaries for the Docker client are available for Linux and macOS (as `docker`).

This topic discusses binary installation for both Linux and macOS:

- [Install daemon and client binaries on Linux](#install-daemon-and-client-binaries-on-linux )
- [Install client binaries on macOS](#install-client-binaries-on-macos )

## Install daemon and client binaries on Linux

### Prerequisites

Before attempting to install Docker from binaries, be sure your host machine
meets the prerequisites:

- A 64-bit installation
- Version 3.10 or higher of the Linux kernel. The latest version of the kernel
  available for your platform is recommended.
- `iptables` version 1.4 or higher
- `git` version 1.7 or higher
- A `ps` executable, usually provided by `procps` or a similar package.
- [XZ Utils](https://tukaani.org/xz/) 4.9 or higher
- A [properly mounted](
  https://github.com/tianon/cgroupfs-mount/blob/master/cgroupfs-mount)
  `cgroupfs` hierarchy; a single, all-encompassing `cgroup` mount
  point is not sufficient. See Github issues
  [#2683](https://github.com/moby/moby/issues/2683),
  [#3485](https://github.com/moby/moby/issues/3485),
  [#4568](https://github.com/moby/moby/issues/4568)).

#### Secure your environment as much as possible

##### OS considerations

Enable SELinux or AppArmor if possible.

It is recommended to use AppArmor or SELinux if your Linux distribution supports
either of the two. This helps improve security and blocks certain
types of exploits. Review the documentation for your Linux distribution for
instructions for enabling and configuring AppArmor or SELinux.

> Security Warning
>
> If either of the security mechanisms is enabled, do not disable it as a
> work-around to make Docker or its containers run. Instead, configure it
> correctly to fix any problems.
{:.warning}

##### Docker daemon considerations

- Enable `seccomp` security profiles if possible. See
  [Enabling `seccomp` for Docker](../security/seccomp.md).

- Enable user namespaces if possible. See the
  [Daemon user namespace options](/engine/reference/commandline/dockerd/#daemon-user-namespace-options).

### Install static binaries

1.  Download the static binary archive. Go to
    [https://download.docker.com/linux/static/stable/](https://download.docker.com/linux/static/stable/)
    (or change `stable` to `nightly` or `test`),
    choose your hardware platform, and download the `.tgz` file relating to the
    version of Docker Engine you want to install.

2.  Extract the archive using the `tar` utility. The `dockerd` and `docker`
    binaries are extracted.

    ```console
    $ tar xzvf /path/to/<FILE>.tar.gz
    ```

3.  **Optional**: Move the binaries to a directory on your executable path, such
    as `/usr/bin/`. If you skip this step, you must provide the path to the
    executable when you invoke `docker` or `dockerd` commands.

    ```console
    $ sudo cp docker/* /usr/bin/
    ```

4.  Start the Docker daemon:

    ```console
    $ sudo dockerd &
    ```

    If you need to start the daemon with additional options, modify the above
    command accordingly or create and edit the file `/etc/docker/daemon.json`
    to add the custom configuration options.

5.  Verify that Docker is installed correctly by running the `hello-world`
    image.

    ```console
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

## Install client binaries on macOS

The macOS binary includes the Docker client only. It does not include the
`dockerd` daemon.

1.  Download the static binary archive. Go to
    [https://download.docker.com/mac/static/stable/x86_64/](https://download.docker.com/mac/static/stable/x86_64/),
    (or change `stable` to `nightly` or `test`),
    and download the `.tgz` file relating to the version of Docker Engine you want
    to install.

2.  Extract the archive using the `tar` utility. The `docker` binary is
    extracted.

    ```console
    $ tar xzvf /path/to/<FILE>.tar.gz
    ```

3.  Clear the extended attributes to allow it run.
    
    In case executing `docker/docker` you get the error message: *'docker' is*
    *damaged and cannot be opened. You should move it to the bin.*
    
    Apple takes care about our security. Hence, we need to remove the security 
    mechanism preventing us running the executable.
    
    ```console
    $ sudo xattr -rc docker
    ```
    
    Now, when you run the following command, you can see the Docker CLI usage instructions:
    
    ```console
    $ docker/docker
    ```

4.  **Optional**: Move the binary to a directory on your executable path, such
    as `/usr/local/bin/`. If you skip this step, you must provide the path to the
    executable when you invoke `docker` or `dockerd` commands.

    ```console
    $ sudo cp docker/docker /usr/local/bin/
    ```

5.  Verify that Docker is installed correctly by running the `hello-world`
    image. The value of `<hostname>` is a hostname or IP address running the
    Docker daemon and accessible to the client.

    ```console
    $ sudo docker -H <hostname> run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

## Upgrade static binaries

To upgrade your manual installation of Docker Engine, first stop any
`dockerd` or `dockerd.exe`  processes running locally, then follow the
regular installation steps to install the new version on top of the existing
version.

## Next steps

- Continue to [Post-installation steps for Linux](linux-postinstall.md).
- Take a look at the [Get started](../../get-started/index.md) training modules to learn  how to build an image and run it as a containerized application.
- Review the topics in [Develop with Docker](../../develop/index.md) to learn how to build new applications using Docker.
