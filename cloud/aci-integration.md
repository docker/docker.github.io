---
title: Deploying Docker containers on Azure
description: Deploying Docker containers on Azure
keywords: Docker, Azure, Integration, ACI, context, Compose, cli, deploy, containers, cloud
redirect_from:
  - /engine/context/aci-integration/
toc_min: 1
toc_max: 2
---

## Overview

The Docker Azure Integration enables developers to use native Docker commands to run applications in Azure Container Instances (ACI) when building cloud-native applications. The new experience provides a tight integration between Docker Desktop and Microsoft Azure allowing developers to quickly run applications using the Docker CLI or VS Code extension, to switch seamlessly from local development to cloud deployment.

In addition, the integration between Docker and Microsoft developer technologies allow developers to use the Docker CLI to:

- Easily log into Azure
- Set up an ACI context in one Docker command allowing you to switch from a local context to a cloud context and run applications quickly and easily
- Simplify single container and multi-container application development using the Compose specification, allowing a developer to invoke fully Docker-compatible commands seamlessly for the first time natively within a cloud container service

Also see the [full list of container features supported by ACI](aci-container-features.md) and [full list of compose features supported by ACI](aci-compose-features.md).

## Prerequisites

To deploy Docker containers on Azure, you must meet the following requirements:

1. Download and install Docker Desktop Stable version 2.3.0.5 or later, or Edge version 2.3.2.0 or later.

   - [Download for Mac](https://desktop.docker.com/mac/edge/Docker.dmg){: target="_blank" rel="noopener" class="_"}
   - [Download for Windows](https://desktop.docker.com/win/edge/Docker%20Desktop%20Installer.exe){: target="_blank" rel="noopener" class="_"}

   Alternatively, install the [Docker Compose CLI for Linux](#install-the-docker-compose-cli-on-linux).

2. Ensure you have an Azure subscription. You can get started with an [Azure free account](https://aka.ms/AA8r2pj){: target="_blank" rel="noopener" class="_"}.

## Run Docker containers on ACI

Docker not only runs containers locally, but also enables developers to seamlessly deploy Docker containers on ACI using `docker run` or deploy multi-container applications defined in a Compose file using the `docker-compose up` command.

The following sections contain instructions on how to deploy your Docker containers on ACI.
Also see the [full list of container features supported by ACI](aci-container-features.md).

### Log into Azure

Run the following commands to log into Azure:

```console
docker login azure
```

This opens your web browser and prompts you to enter your Azure login credentials.
If the Docker CLI cannot open a browser, it will fall back to the [Azure device code flow](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-device-code){:target="_blank" rel="noopener" class="_"} and lets you connect manually.
Note that the [Azure command line](https://docs.microsoft.com/en-us/cli/azure/){:target="_blank" rel="noopener" class="_"} login is separated from the Docker CLI Azure login.

Alternatively, you can log in without interaction (typically in
scripts or continuous integration scenarios), using an Azure Service
Principal, with `docker login azure --client-id xx --client-secret yy --tenant-id zz`

> **Note**
>
> Logging in through the Azure Service Provider obtains an access token valid
> for a short period (typically 1h), but it does not allow you to automatically
> and transparently refresh this token. You must manually re-login
> when the access token has expired when logging in with a Service Provider.

You can also use the `--tenant-id` option alone to specify a tenant, if
you have several ones available in Azure.

### Create an ACI context

After you have logged in, you need to create a Docker context associated with ACI to deploy containers in ACI.
Creating an ACI context requires an Azure subscription, a [resource group](https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/manage-resource-groups-portal), and a region.
For example, let us create a new context called `myacicontext`:

```console
docker context create aci myacicontext
```

This command automatically uses your Azure login credentials to identify your subscription IDs and resource groups. You can then interactively select the subscription and group that you would like to use. If you prefer, you can specify these options in the CLI using the following flags: `--subscription-id`,
`--resource-group`, and `--location`.

If you don't have any existing resource groups in your Azure account, the `docker context create aci myacicontext` command creates one for you. You don’t have to specify any additional options to do this.

After you have created an ACI context, you can list your Docker contexts by running the `docker context ls` command:

```console
NAME                TYPE                DESCRIPTION                               DOCKER ENDPOINT                KUBERNETES ENDPOINT   ORCHESTRATOR
myacicontext        aci                 myResourceGroupGTA@eastus
default *           moby              Current DOCKER_HOST based configuration   unix:///var/run/docker.sock                          swarm
```

### Run a container

Now that you've logged in and created an ACI context, you can start using Docker commands to deploy containers on ACI.

There are two ways to use your new ACI context. You can use the `--context` flag with the Docker command to specify that you would like to run the command using your newly created ACI context.

```console
docker --context myacicontext run -p 80:80 nginx
```

Or, you can change context using `docker context use` to select the ACI context to be your focus for running Docker commands. For example, we can use the `docker context use` command to deploy an ngnix container:

```console
docker context use myacicontext
docker run -p 80:80 nginx
```

After you've switched to the `myacicontext` context, you can use `docker ps` to list your containers running on ACI.

In the case of the demonstration nginx container started above, the result of the ps command will display in column "PORTS" the IP address and port on which the container is running. For example, it may show `52.154.202.35:80->80/tcp`, and you can view the nginx welcome page by browsing `http://52.154.202.35`.

To view logs from your container, run:

```console
docker logs <CONTAINER_ID>
```

To execute a command in a running container, run:

```console
docker exec -t <CONTAINER_ID> COMMAND
```

To stop and remove a container from ACI, run:

```console
docker stop <CONTAINER_ID>
docker rm <CONTAINER_ID>
```

You can remove containers using `docker rm`. To remove a running container, you must use the `--force` flag, or stop the container using `docker stop` before removing it.

> **Note**
>
> The semantics of restarting a container on ACI are different to those when using a local Docker context for local development. On ACI, the container will be reset to its initial state and started on a new node. This includes the container's filesystem so all state that is not stored in a volume will be lost on restart.

## Running Compose applications

You can also deploy and manage multi-container applications defined in Compose files to ACI using the `docker-compose` command.
All containers in the same Compose application are started in the same container group. Service discovery between the containers works using the service name specified in the Compose file.
Name resolution between containers is achieved by writing service names in the `/etc/hosts` file that is shared automatically by all containers in the container group.

Also see the [full list of compose features supported by ACI](aci-compose-features.md).

1. Ensure you are using your ACI context. You can do this either by specifying the `--context myacicontext` flag or by setting the default context using the command `docker context use myacicontext`.

2. Run `docker-compose up` and `docker-compose down` to start and then stop a full Compose application.

By default, `docker-compose up` uses the `docker-compose.yaml` file in the current folder. You can specify the working directory using the --workdir flag or specify the Compose file directly using the `--file` flag.

You can also specify a name for the Compose application using the `--project-name` flag during deployment. If no name is specified, a name will be derived from the working directory.

Containers started as part of Compose applications will be displayed along with single containers when using `docker ps`. Their container ID will be of the format: `<COMPOSE-PROJECT>_<SERVICE>`.
These containers cannot be stopped, started, or removed independently since they are all part of the same ACI container group.
You can view each container's logs with `docker logs`. You can list deployed Compose applications with `docker-compose ls`. This will list only compose applications, not single containers started with `docker run`. You can remove a Compose application with `docker-compose down`.

> **Note**
>
> The current Docker Azure integration does not allow fetching a combined log stream from all the containers that make up the Compose application.

## Updating applications

From a deployed Compose application, you can update the application by re-deploying it with the same project name: `docker-compose up --project-name PROJECT`.

Updating an application means the ACI node will be reused, and the application will keep the same IP address that was previously allocated to expose ports, if any. ACI has some limitations on what can be updated in an existing application (you will not be able to change CPU/memory reservation for example), in these cases, you need to deploy a new application from scratch.

Updating is the default behavior if you invoke `docker-compose up` on an already deployed Compose file, as the Compose project name is derived from the directory where the Compose file is located by default. You need to explicitly execute `docker-compose down` before running `docker-compose up` again in order to totally reset a Compose application.

## Releasing resources

Single containers and Compose applications can be removed from ACI with
the `docker prune` command. The `docker prune` command removes deployments
that are not currently running. To remove running depoyments, you can specify
`--force`. The `--dry-run` option lists deployments that are planned for
removal, but it doesn't actually remove them.

```console
$ ./bin/docker --context acicontext prune --dry-run --force
Resources that would be deleted:
my-application
Total CPUs reclaimed: 2.01, total memory reclaimed: 2.30 GB
```

## Exposing ports

Single containers and Compose applications can optionally expose ports.
For single containers, this is done using the `--publish` (`-p`) flag of the `docker run` command : `docker run -p 80:80 nginx`.

For Compose applications, you must specify exposed ports in the Compose file service definition:

```yaml
services:
  nginx:
    image: nginx
    ports:
      - "80:80"
```

> **Note**
>
> ACI does not allow port mapping (that is, changing port number while exposing port). Therefore, the source and target ports must be the same when deploying to ACI.
>
> All containers in the same Compose application are deployed in the same ACI container group. Different containers in the same Compose application cannot expose the same port when deployed to ACI.

By default, when exposing ports for your application, a random public IP address is associated with the container group supporting the deployed application (single container or Compose application).
This IP address can be obtained when listing containers with `docker ps` or using `docker inspect`.

### DNS label name

In addition to exposing ports on a random IP address, you can specify a DNS label name to expose your application on an FQDN of the form: `<NAME>.region.azurecontainer.io`.

You can set this name with the `--domainname` flag when performing a `docker run`, or by using the `domainname` field in the Compose file when performing a `docker-compose up`:

```yaml
services:
  nginx:
    image: nginx
    domainname: "myapp"
    ports:
      - "80:80"
```

> **Note**
>
> The domain of a Compose application can only be set once, if you specify the
> `domainname` for several services, the value must be identical.
>
> The FQDN `<DOMAINNAME>.region.azurecontainer.io` must be available.

## Using Azure file share as volumes in ACI containers

You can deploy containers or Compose applications that use persistent data
stored in volumes. Azure File Share can be used to support volumes for ACI
containers.

Using an existing Azure File Share with storage account name `mystorageaccount`
and file share name `myfileshare`, you can specify a volume in your deployment `run`
command as follows:

```console
docker run -v mystorageaccount/myfileshare:/target/path myimage
```

The runtime container will see the file share content in `/target/path`.

In a Compose application, the volume specification must use the following syntax
in the Compose file:

```yaml
myservice:
  image: nginx
  volumes:
    - mydata:/mount/testvolumes

volumes:
  mydata:
    driver: azure_file
    driver_opts:
      share_name: myfileshare
      storage_account_name: mystorageaccount
```

> **Note**
>
> The volume short syntax in Compose files cannot be used as it is aimed at volume definition for local bind mounts. Using the volume driver and driver option syntax in Compose files makes the volume definition a lot more clear.

In single or multi-container deployments, the Docker CLI will use your Azure login to fetch the key to the storage account, and provide this key with the container deployment information, so that the container can access the volume.
Volumes can be used from any file share in any storage account you have access to with your Azure login. You can specify `rw` (read/write) or `ro` (read only) when mounting the volume (`rw` is the default).

### Managing Azure volumes

To create a volume that you can use in containers or Compose applications when
using your ACI Docker context, you can use the `docker volume create` command,
and specify an Azure storage account name and the file share name:

```console
$ docker --context aci volume create test-volume --storage-account mystorageaccount
[+] Running 2/2
 ⠿ mystorageaccount  Created                         26.2s
 ⠿ test-volume       Created                          0.9s
mystorageaccount/test-volume
```

By default, if the storage account does not already exist, this command
creates a new storage account using the Standard LRS as a default SKU, and the
resource group and location associated with your Docker ACI context.

If you specify an existing storage account, the command creates a new
file share in the existing account:

```console
$ docker --context aci volume create test-volume2 --storage-account mystorageaccount
[+] Running 2/2
 ⠿ mystorageaccount   Use existing                    0.7s
 ⠿ test-volume2       Created                         0.7s
mystorageaccount/test-volume2
```

Alternatively, you can create an Azure storage account or a file share using the Azure
portal, or the `az` [command line](https://docs.microsoft.com/en-us/azure/storage/files/storage-how-to-use-files-cli).

You can also list volumes that are available for use in containers or Compose applications:

```console
$ docker --context aci volume ls
ID                                 DESCRIPTION
mystorageaccount/test-volume       Fileshare test-volume in mystorageaccount storage account
mystorageaccount/test-volume2      Fileshare test-volume2 in mystorageaccount storage account
```

To delete a volume and the corresponding Azure file share, use the `volume rm` command:

```console
$ docker --context aci volume rm mystorageaccount/test-volume
mystorageaccount/test-volume
```

This permanently deletes the Azure file share and all its data.

When deleting a volume in Azure, the command checks whether the specified file share
is the only file share available in the storage account. If the storage account is
created with the `docker volume create` command, `docker volume rm` also
deletes the storage account when it does not have any file shares.
If you are using a storage account created without the `docker volume create` command
(through Azure portal or with the `az` command line for example), `docker volume rm`
does not delete the storage account, even when it has zero remaining file shares.

## Environment variables

When using `docker run`, you can pass the environment variables to ACI containers using the `--env` flag.
For Compose applications, you can specify the environment variables in the Compose file with the `environment` or `env-file` service field, or with the `--environment` command line flag.

## Health checks

You can specify a container health checks using either the `--healthcheck-` prefixed flags with `docker run`, or in a Compose file with the `healthcheck` section of the service.

Health checks are converted to ACI `LivenessProbe`s. ACI runs the health check command periodically, and if it fails, the container will be terminated.

Health checks must be used in addition to restart policies to ensure the container is then restarted on termination. The default restart policy for `docker run` is `no` which will not restart the container. The default restart policy for Compose is `any` which will always try restarting the service containers.

Example using `docker run`:

```console
docker --context acicontext run -p 80:80 --restart always --health-cmd "curl http://localhost:80" --health-interval 3s  nginx
```

Example using Compose files:

```yaml
services:
  web:
    image: nginx
    deploy:
      restart_policy:
        condition: on-failure
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:80"]
      interval: 10s
```

## Private Docker Hub images and using the Azure Container Registry

You can deploy private images to ACI that are hosted by any container registry. You need to log into the relevant registry using `docker login` before running `docker run` or `docker-compose up`. The Docker CLI will fetch your registry login for the deployed images and send the credentials along with the image deployment information to ACI.
In the case of the Azure Container Registry, the command line will try to automatically log you into ACR from your Azure login. You don't need to manually login to the ACR registry first, if your Azure login has access to the ACR.

## Using ACI resource groups as namespaces

You can create several Docker contexts associated with ACI. Each context must be associated with a unique Azure resource group. This allows you to use Docker contexts as namespaces. You can switch between namespaces using `docker context use <CONTEXT>`.

When you run the `docker ps` command, it only lists containers in your current Docker context. There won’t be any contention in container names or Compose application names between two Docker contexts.

## Install the Docker Compose CLI on Linux

The Docker Compose CLI adds support for running and managing containers on Azure Container Instances (ACI).

### Install Prerequisites

- [Docker 19.03 or later](https://docs.docker.com/get-docker/)

### Install script

You can install the new CLI using the install script:

```console
curl -L https://raw.githubusercontent.com/docker/compose-cli/main/scripts/install/install_linux.sh | sh
```

### Manual install

You can download the Docker ACI Integration CLI from the
[latest release](https://github.com/docker/compose-cli/releases/latest){: target="_blank" rel="noopener" class="_"} page.

You will then need to make it executable:

```console
chmod +x docker-aci
```

To enable using the local Docker Engine and to use existing Docker contexts, you
must have the existing Docker CLI as `com.docker.cli` somewhere in your
`PATH`. You can do this by creating a symbolic link from the existing Docker
CLI:

```console
ln -s /path/to/existing/docker /directory/in/PATH/com.docker.cli
```

> **Note**
>
> The `PATH` environment variable is a colon-separated list of
> directories with priority from left to right. You can view it using
> `echo $PATH`. You can find the path to the existing Docker CLI using
> `which docker`. You may need root permissions to make this link.

On a fresh install of Ubuntu 20.04 with Docker Engine
[already installed](https://docs.docker.com/engine/install/ubuntu/):

```console
$ echo $PATH
/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin
$ which docker
/usr/bin/docker
$ sudo ln -s /usr/bin/docker /usr/local/bin/com.docker.cli
```

You can verify that this is working by checking that the new CLI works with the
default context:

```console
$ ./docker-aci --context default ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
$ echo $?
0
```

To make this CLI with ACI integration your default Docker CLI, you must move it
to a directory in your `PATH` with higher priority than the existing Docker CLI.

Again, on a fresh Ubuntu 20.04:

```console
$ which docker
/usr/bin/docker
$ echo $PATH
/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin
$ sudo mv docker-aci /usr/local/bin/docker
$ which docker
/usr/local/bin/docker
$ docker version
...
 Azure integration  0.1.4
...
```

### Supported commands

After you have installed the Docker ACI Integration CLI, run `--help` to see the current list of commands.

### Uninstall

To remove the Docker Azure Integration CLI, you need to remove the binary you downloaded and `com.docker.cli` from your `PATH`. If you installed using the script, this can be done as follows:

```console
sudo rm /usr/local/bin/docker /usr/local/bin/com.docker.cli
```

## Feedback

Thank you for trying out Docker Azure Integration. Your feedback is very important to us. Let us know your feedback by creating an issue in the [compose-cli](https://github.com/docker/compose-cli){: target="_blank" rel="noopener" class="_"} GitHub repository.
