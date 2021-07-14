---
description: Get started with Docker Machine and a local VM
keywords: machine, virtualbox
title: Get started with Docker Machine and a local VM
---

Let's take a look at using `docker-machine` to create, use, and manage a
Docker host inside of a local virtual machine.

> **Note**
>
> Docker Machine is no longer included in the Docker Desktop installer. You can download it separately from the [docker/machine](https://github.com/docker/machine/releases/){: target="_blank" rel="noopener" class="_" } release
page on GitHub.

## Prerequisite information

With the advent of [Docker Desktop for Mac](../docker-for-mac/index.md) and
[Docker Desktop for Windows](../docker-for-windows/index.md), we recommend that you use these for your
primary Docker workflows. You can use these applications to run Docker natively
on your local system without using Docker Machine at all.

For now, however, if you want to create _multiple_ local machines, you still
need Docker Machine to create and manage machines for multi-node
experimentation.

The new solutions come with their own native virtualization solutions rather
than Oracle VirtualBox, so keep the following considerations in mind when using
Machine to create local VMs.

* **Docker Desktop for Mac** - You can use `docker-machine create` with the `virtualbox` driver to create additional local machines.

* **Docker Desktop for Windows** - You can use `docker-machine create` with the `hyperv` driver to create additional local machines.

#### If you are using Docker Desktop for Windows

Docker Desktop for Windows uses [Microsoft
Hyper-V](https://msdn.microsoft.com/en-us/virtualization/hyperv_on_windows/windows_welcome)
for virtualization, and Hyper-V is not compatible with Oracle VirtualBox.
Therefore, you cannot run the two solutions simultaneously. But you can still
use `docker-machine` to create more local VMs by using the Microsoft Hyper-V
driver.

The prerequisites are:

* Have Docker Desktop for Windows installed, and running (which requires that virtualization and Hyper-V are enabled, as described in [What to know before you install Docker Desktop for Windows](../docker-for-windows/install.md#what-to-know-before-you-install)).

* Set up the Hyper-V driver to use an external virtual network switch See
the [Docker Machine driver for Microsoft Hyper-V](drivers/hyper-v.md) topic,
which includes an [example](drivers/hyper-v.md#example) of how to do this.

#### If you are using Docker Desktop for Mac

Docker Desktop for Mac uses [HyperKit](https://github.com/docker/HyperKit/), a
lightweight macOS virtualization solution built on top of the
[Hypervisor.framework](https://developer.apple.com/reference/hypervisor).

Currently, there is no `docker-machine create` driver for HyperKit, so
use the `virtualbox` driver to create local machines. (See the
[Docker Machine driver for Oracle VirtualBox](drivers/virtualbox.md).) You can
run both HyperKit and Oracle VirtualBox on the same system.

* Make sure you have [the latest VirtualBox](https://www.virtualbox.org/wiki/Downloads){: target="_blank" rel="noopener" class="_"}
  correctly installed on your system.

##  Use Machine to run Docker containers

To run a Docker container, you:

* create a new (or start an existing) Docker virtual machine
* switch your environment to your new VM
* use the docker client to create, load, and manage containers

Once you create a machine, you can reuse it as often as you like. Like any VirtualBox VM, it maintains its configuration between uses.

The examples here show how to create and start a machine, run Docker commands, and work with containers.

## Create a machine

1. Open a command shell or terminal window.

    These command examples shows a Bash shell. For a different shell, such as C Shell, the same commands are the same except where noted.

2. Use `docker-machine ls` to list available machines.

    In this example, no machines have been created yet.

        $ docker-machine ls
        NAME   ACTIVE   DRIVER   STATE   URL   SWARM   DOCKER   ERRORS

3. Create a machine.

    Run the `docker-machine create` command, pass the appropriate driver to the
`--driver` flag and provide a machine name. If this is your first machine, name
it `default` as shown in the example. If you already have a "default" machine,
choose another name for this new machine.

    On Docker Desktop for Windows systems that support Hyper-V, use the `hyperv` driver as shown in the [Docker Machine Microsoft Hyper-V driver reference](drivers/hyper-v.md) and follow the [example](drivers/hyper-v.md#example), which shows how to use an external network switch and provides the flags for the full command. (See [prerequisites](get-started.md#prerequisite-information) above to learn more.)

            $ docker-machine create --driver virtualbox default
            Running pre-create checks...
            Creating machine...
            (staging) Copying /Users/ripley/.docker/machine/cache/boot2docker.iso to /Users/ripley/.docker/machine/machines/default/boot2docker.iso...
            (staging) Creating VirtualBox VM...
            (staging) Creating SSH key...
            (staging) Starting the VM...
            (staging) Waiting for an IP...
            Waiting for machine to be running, this may take a few minutes...
            Machine is running, waiting for SSH to be available...
            Detecting operating system of created instance...
            Detecting the provisioner...
            Provisioning with boot2docker...
            Copying certs to the local machine directory...
            Copying certs to the remote machine...
            Setting Docker configuration on the remote daemon...
            Checking connection to Docker...
            Docker is up and running!
            To see how to connect Docker to this machine, run: docker-machine env default

      This command downloads a lightweight Linux distribution ([boot2docker](https://github.com/boot2docker/boot2docker){: target="_blank" rel="noopener" class="_"}) with the Docker daemon installed, and creates and starts a VirtualBox VM with Docker running.

4. List available machines again to see your newly minted machine.

        $ docker-machine ls
        NAME      ACTIVE   DRIVER       STATE     URL                         SWARM   DOCKER   ERRORS
        default   *        virtualbox   Running   tcp://192.168.99.187:2376           v1.9.1

5. Get the environment commands for your new VM.

    As noted in the output of the `docker-machine create` command, you need to tell Docker to talk to the new machine. You can do this with the `docker-machine env` command.

        $ docker-machine env default
        export DOCKER_TLS_VERIFY="1"
        export DOCKER_HOST="tcp://172.16.62.130:2376"
        export DOCKER_CERT_PATH="/Users/<yourusername>/.docker/machine/machines/default"
        export DOCKER_MACHINE_NAME="default"
        # Run this command to configure your shell:
        # eval "$(docker-machine env default)"

6. Connect your shell to the new machine.

        $ eval "$(docker-machine env default)"

      **Note**: If you are using `fish`, or a Windows shell such as
      Powershell/`cmd.exe`, the above method does not work as described.
      Instead, see [the `env` command's documentation](reference/env.md){: target="_blank" rel="noopener" class="_"}
      to learn how to set the environment variables for your shell.

    This sets environment variables for the current shell that the Docker
    client reads which specify the TLS settings. You need to do this
    each time you open a new shell or restart your machine. (See also, how to
      [unset environment variables in the current shell](get-started.md#unset-environment-variables-in-the-current-shell).)

    You can now run Docker commands on this host.

## Run containers and experiment with Machine commands

Run a container with `docker run` to verify your set up.

1. Use `docker run` to download and run `busybox` with a simple 'echo' command.

        $ docker run busybox echo hello world
        Unable to find image 'busybox' locally
        Pulling repository busybox
        e72ac664f4f0: Download complete
        511136ea3c5a: Download complete
        df7546f9f060: Download complete
        e433a6c5b276: Download complete
        hello world

2. Get the host IP address.

    Any exposed ports are available on the Docker host’s IP address, which you can get using the `docker-machine ip` command:

        $ docker-machine ip default
        192.168.99.100

3. Run a [Nginx](https://www.nginx.com/){: target="_blank" rel="noopener" class="_"} webserver in a container with the following command:

        $ docker run -d -p 8000:80 nginx

    When the image is finished pulling, you can hit the server at port 8000 on the IP address given to you by `docker-machine ip`. For instance:

            $ curl $(docker-machine ip default):8000
            <!DOCTYPE html>
            <html>
            <head>
            <title>Welcome to nginx!</title>
            <style>
                body {
                    width: 35em;
                    margin: 0 auto;
                    font-family: Tahoma, Verdana, Arial, sans-serif;
                }
            </style>
            </head>
            <body>
            <h1>Welcome to nginx!</h1>
            <p>If you see this page, the nginx web server is successfully installed and
            working. Further configuration is required.</p>

            <p>For online documentation and support, refer to
            <a href="https://nginx.org">nginx.org</a>.<br/>
            Commercial support is available at
            <a href="https://www.nginx.com">www.nginx.com</a>.</p>

            <p><em>Thank you for using nginx.</em></p>
            </body>
            </html>

  You can create and manage as many local VMs running Docker as your local resources permit; just run `docker-machine create` again. All created machines appear in the output of `docker-machine ls`.

## Start and stop machines

If you are finished using a host for the time being, you can stop it with `docker-machine stop` and later start it again with `docker-machine start`.

        $ docker-machine stop default
        $ docker-machine start default

## Operate on machines without specifying the name

Some `docker-machine` commands assume that the given operation should be run on a machine named `default` (if it exists) if no machine name is specified.  Because using a local VM named `default` is such a common pattern, this allows you to save some typing on the most frequently used Machine commands.

For example:

          $ docker-machine stop
          Stopping "default"....
          Machine "default" was stopped.

          $ docker-machine start
          Starting "default"...
          (default) Waiting for an IP...
          Machine "default" was started.
          Started machines may have new IP addresses.  You may need to re-run the `docker-machine env` command.

          $ eval $(docker-machine env)

          $ docker-machine ip
            192.168.99.100

Commands that follow this style are:

        - `docker-machine config`
        - `docker-machine env`
        - `docker-machine inspect`
        - `docker-machine ip`
        - `docker-machine kill`
        - `docker-machine provision`
        - `docker-machine regenerate-certs`
        - `docker-machine restart`
        - `docker-machine ssh`
        - `docker-machine start`
        - `docker-machine status`
        - `docker-machine stop`
        - `docker-machine upgrade`
        - `docker-machine url`

For machines other than `default`, and commands other than those listed above, you must always specify the name explicitly as an argument.

## Unset environment variables in the current shell

You might want to use the current shell to connect to a different Docker Engine.
In such scenarios, you have the option to switch the environment for the current
shell to talk to different Docker engines.

1.  Run `env|grep DOCKER` to check whether DOCKER environment variables are set.

    ```none
    $ env | grep DOCKER
    DOCKER_HOST=tcp://192.168.99.100:2376
    DOCKER_MACHINE_NAME=default
    DOCKER_TLS_VERIFY=1
    DOCKER_CERT_PATH=/Users/<your_username>/.docker/machine/machines/default
    ```

    If it returns output (as shown in the example), you can unset the `DOCKER` environment variables.

2.  Use one of two methods to unset DOCKER environment variables in the current shell.

    * Run the `unset` command on the following `DOCKER` environment variables.

      ```none
      unset DOCKER_TLS_VERIFY
      unset DOCKER_CERT_PATH
      unset DOCKER_MACHINE_NAME
      unset DOCKER_HOST
      ```

    * Alternatively, run a shortcut command `docker-machine env -u` to show the command you need to run to unset all DOCKER variables:

      ```none
      $ docker-machine env -u
      unset DOCKER_TLS_VERIFY
      unset DOCKER_HOST
      unset DOCKER_CERT_PATH
      unset DOCKER_MACHINE_NAME
      # Run this command to configure your shell:
      # eval $(docker-machine env -u)
      ```

      Run `eval $(docker-machine env -u)` to unset all DOCKER variables in the current shell.

3. Now, after running either of the above commands, this command should return no output.

    ```
    $ env | grep DOCKER
    ```

    If you are running Docker Desktop for Mac, you can run Docker commands to talk
    to the Docker Engine installed with that app.

## Start local machines on startup

To ensure that the Docker client is automatically configured at the start of
each shell session, you can embed `eval $(docker-machine env default)` in your
shell profiles, by adding it to the `~/.bash_profile` file or the equivalent
configuration file for your shell. However, this fails if a machine called
`default` is not running. You can configure your system to start the `default`
machine automatically. The following examples show how to do this in macOS.

### If you are using Homebrew

`brew services` is available by default in any [Homebrew installation](https://brew.sh/). It is a convenience wrapper around `launchd`.

`brew services start docker-machine`

### Using launchd

`launchd` comes with macOS by default, so it is not necessary to install anything.

Create a file called `com.docker.machine.default.plist` in the
`~/Library/LaunchAgents/` directory, with the following content:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>EnvironmentVariables</key>
        <dict>
            <key>PATH</key>
            <string>/usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin</string>
        </dict>
        <key>Label</key>
        <string>com.docker.machine.default</string>
        <key>ProgramArguments</key>
        <array>
            <string>/usr/local/bin/docker-machine</string>
            <string>start</string>
            <string>default</string>
        </array>
        <key>RunAtLoad</key>
        <true/>
    </dict>
</plist>
```

You can change the `default` string above to make this `LaunchAgent` start a
different machine.

## Where to go next

-   Provision multiple Docker hosts [on your cloud provider](get-started-cloud.md)
-   [Understand Machine concepts](concepts.md)
- [Docker Machine list of reference pages for all supported drivers](drivers/index.md)
- [Docker Machine driver for Oracle VirtualBox](drivers/virtualbox.md)
- [Docker Machine driver for Microsoft Hyper-V](drivers/hyper-v.md)
- [`docker-machine` command line reference](reference/index.md)
