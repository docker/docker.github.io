---
description: Setup for voting app example
keywords: docker-machine, multi-container, services, swarm mode, cluster, voting app
title: Set up Dockerized machines
---

The first step in getting the voting app deployed is to set up Docker machines
for the swarm nodes. You could create these Docker hosts on different physical
machines, virtual machines, or cloud providers.

For this example, we use [Docker Machine](/machine/get-started.md) to create two
virtual machines on a single system. (See [Docker Machine
Overview](/machine/overview.md) to learn more.) We'll also verify the setup, and
run some basic commmands to interact with the machines.

## Create `manager` and `worker` machines

The Docker Machine command to create a local virtual machine is:

```
docker-machine create --driver virtualbox <host-name>
```

Here is an example of creating the manager. Create this one, then do the same for `worker`.

```
$  docker-machine create --driver virtualbox manager
Running pre-create checks...
Creating machine...
(manager) Copying /Users/victoria/.docker/machine/cache/boot2docker.iso to /Users/victoria/.docker/machine/machines/manager/boot2docker.iso...
(manager) Creating VirtualBox VM...
(manager) Creating SSH key...
(manager) Starting the VM...
(manager) Check network to re-create if needed...
(manager) Waiting for an IP...
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with boot2docker...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Checking connection to Docker...
Docker is up and running!
To see how to connect your Docker Client to the Docker Engine running on this virtual machine, run: docker-machine env manager
```

## Get the IP addresses of the machines

Use `docker-machine ls` to verify that the machines are running and to get their
IP addresses.

```
$ docker-machine ls
NAME       ACTIVE   DRIVER       STATE     URL                         SWARM   DOCKER        ERRORS
manager   *        virtualbox   Running   tcp://192.168.99.100:2376           v1.13.0-rc6
worker    -        virtualbox   Running   tcp://192.168.99.101:2376           v1.13.0-rc6
```
You now have two "Dockerized" machines, each running Docker Engine, accessible through the [Docker CLI](/engine/reference/commandline.md), and available to become swarm nodes.

## Interacting with the machines

There are several ways to interact with these machines directly on the command line or programatically. We'll cover two methods for managing the machines directly from the command line.

### Get the environment commands for your new VM

As noted in the last line of output of the `docker-machine create` command, you need to tell Docker to talk to the new machine from a current shell. You can do this with the `docker-machine env` command.

#### Manage the machines from a command shell

1. Run `docker-machine env manager` to get environment variables for that machine.

        $ docker-machine env manager
        export DOCKER_TLS_VERIFY="1"
        export DOCKER_HOST="tcp://192.168.99.100:2376"
        export DOCKER_CERT_PATH="/Users/victoriabialas/.docker/machine/machines/manager"
        export DOCKER_MACHINE_NAME="manager"
        export DOCKER_API_VERSION="1.25"
        # Run this command to configure your shell:
        # eval $(docker-machine env manager)

2. Connect your shell to the manager.

        $ eval $(docker-machine env manager)

    This sets environment variables for the current shell that the Docker client will read which specify the TLS settings. You need to do this each time you open a new shell or restart your machine.

    **Note**: If you are using `fish`, or a Windows shell such as
    Powershell/`cmd.exe` the above method will not work as described.
    Instead, see the [env command reference documentation](/machine/reference/env.md) to learn how to set the environment variables for your shell.

3. Run `docker-machine ls` again.

        $ docker-machine ls
        NAME      ACTIVE   DRIVER       STATE     URL                         SWARM   DOCKER        ERRORS
        manager   *        virtualbox   Running   tcp://192.168.99.100:2376           v1.13.0-rc6   
        worker    -        virtualbox   Running   tcp://192.168.99.101:2376           v1.13.0-rc6   

  The asterisk next `manager` indicates that the current shell is connected to that machine.

#### Log on to a machine directly

You can use the command `docker-machine ssh <machine-name>` to log onto a machine:

```
$ docker-machine ssh manager
                        ##         .
                  ## ## ##        ==
               ## ## ## ## ##    ===
           /"""""""""""""""""\___/ ===
      ~~~ {~~ ~~~~ ~~~ ~~~~ ~~~ ~ /  ===- ~~~
           \______ o           __/
             \    \         __/
              \____\_______/
 _                 _   ____     _            _
| |__   ___   ___ | |_|___ \ __| | ___   ___| | _____ _ __
| '_ \ / _ \ / _ \| __| __) / _` |/ _ \ / __| |/ / _ \ '__|
| |_) | (_) | (_) | |_ / __/ (_| | (_) | (__|   <  __/ |
|_.__/ \___/ \___/ \__|_____\__,_|\___/ \___|_|\_\___|_|

  WARNING: this is a build from test.docker.com, not a stable release.

Boot2Docker version 1.13.0-rc6, build HEAD : 5ab2289 - Wed Jan 11 23:37:52 UTC 2017
Docker version 1.13.0-rc6, build 2f2d055
```

We'll use this method later in the example.

## What's next?

In the next step, we'll [create a swarm](create-swarm.md) across these two Docker nodes.
