---
description: How to install Docker EE for Windows Server
keywords: Windows, Windows Server, install, download, ucp, Docker EE
title: Install Docker Enterprise Edition for Windows Server
redirect_from:
- /docker-ee-for-windows/install/
- /engine/installation/windows/docker-ee/
---

Docker Enterprise Edition for Windows Server (*Docker EE*) enables native
Docker containers on Windows Server. Windows Server 2016 and later versions are supported. The Docker EE installation package
includes everything you need to run Docker on Windows Server.
This topic describes pre-install considerations, and how to download and
install Docker EE.

> Release notes
>
> You can [get release notes for all versions here](/release-notes/)

## Install Docker EE

Docker EE for Windows requires Windows Server 2016 or later. See
[What to know before you install](#what-to-know-before-you-install) for a
full list of prerequisites.

1.  Open a PowerShell command prompt, and type the following commands.

    ```PowerShell
    Install-Module DockerMsftProvider -Force
    Install-Package Docker -ProviderName DockerMsftProvider -Force
    ```

2.  Check if a reboot is required, and if yes, restart your instance:

    ```PowerShell
    (Install-WindowsFeature Containers).RestartNeeded
    ```
    If the output of this command is **Yes**, then restart the server with:

    ```PowerShell
    Restart-Computer
    ```

3.  Test your Docker EE installation by running the `hello-world` container.

    ```PowerShell
    docker container run hello-world:nanoserver

    Unable to find image 'hello-world:nanoserver' locally
    nanoserver: Pulling from library/hello-world
    bce2fbc256ea: Pull complete
    3ac17e2e6106: Pull complete
    8cac44e17f16: Pull complete
    5e160e4d8db3: Pull complete
    Digest: sha256:25eac12ba40f7591969085ab3fb9772e8a4307553c14ea72d0e6f98b2c8ced9d
    Status: Downloaded newer image for hello-world:nanoserver

    Hello from Docker!
    This message shows that your installation appears to be working correctly.
    ```

### (optional) Make sure you have all required updates

Some advanced Docker features, such as swarm mode, require the fixes included in
[KB4015217](https://support.microsoft.com/en-us/help/4015217/windows-10-update-kb4015217)
(or a later cumulative patch).

```PowerShell
sconfig
```

Select option `6) Download and Install Updates`.

## Install a specific version

To install a specific Docker version, you can use the
`MaximumVersion`,`MinimumVersion` or `RequiredVersion` flags. For example:

```PowerShell
Install-Package -Name docker -ProviderName DockerMsftProvider -Force -RequiredVersion 17.06.2-ee-5
...
Name                           Version          Source           Summary
----                           -------          ------           -------
Docker                         17.06.2-ee-5       Docker           Contains Docker EE for use with Windows Server 2016...
```

## Update Docker EE

To update Docker EE on Windows Server 2016:

```PowerShell
Install-Package -Name docker -ProviderName DockerMsftProvider -Update -Force
```

If Docker Universal Control Plane (UCP) is installed, run the
[UCP installation script for Windows](/datacenter/ucp/2.2/guides/admin/configure/join-windows-worker-nodes/#run-the-windows-node-setup-script).

Start the Docker service:

```PowerShell
Start-Service Docker
```

## What to know before you install

* **What the Docker EE for Windows install includes**: The installation
provides [Docker Engine](/engine/userguide/intro.md) and the
[Docker CLI client](/engine/reference/commandline/cli.md).

## About Docker EE containers and Windows Server

Looking for information on using Docker EE containers?

* [Getting Started with Windows Containers (Lab)](https://github.com/docker/labs/blob/master/windows/windows-containers/README.md)
provides a tutorial on how to set up and run Windows containers on Windows 10
or Windows Server 2016. It shows you how to use a MusicStore application with
Windows containers.

* [Setup - Windows Server 2016 (Lab)](https://github.com/docker/labs/blob/master/windows/windows-containers/Setup-Server2016.md)
describes environment setup in detail.

* Docker Container Platform for Windows Server [articles and blog
posts](https://www.docker.com/microsoft/) on the Docker website.

## Where to go next

* [Getting started](/docker-for-windows/index.md) provides an overview of
Docker for Windows, basic Docker command examples, how to get help or give
feedback, and links to all topics in the Docker for Windows guide.

* [FAQs](/docker-for-windows/faqs.md) provides answers to frequently asked
questions.

* [Release Notes](/docker-for-windows/release-notes.md) lists component
updates, new features, and improvements associated with Stable and Edge
releases.

* [Learn Docker](/learn.md) provides general Docker tutorials.

* [Windows Containers on Windows Server](https://docs.microsoft.com/en-us/virtualization/windowscontainers/quick-start/quick-start-windows-server)
is the official Microsoft documentation.
