---
description: Create a Docker image by layering your ASP.NET Core app on debian for Linux Containers or with Windows Nano Server containers using a Dockerfile.
keywords: dockerize, dockerizing, dotnet, .NET, Core, article, example, platform, installation, containers, images, image, dockerfile, build, asp.net, asp.net core
title: Dockerize an ASP.NET Core application
redirect_from:
- /engine/examples/dotnetcore/
---

## Introduction

This example demonstrates how to dockerize an ASP.NET Core application.

## Why build ASP.NET Core?

- [Open-source](https://github.com/aspnet/home)
- Develop and run your ASP.NET Core apps cross-platform on Windows, MacOS, and
  Linux
- Great for modern cloud-based apps, such as web apps, IoT apps, and mobile
  backends
- ASP.NET Core apps can run on [.NET Core](https://www.microsoft.com/net/core/platform)
  or on the full [.NET Framework](https://www.microsoft.com/net/framework)
- Designed to provide an optimized development framework for apps that are
  deployed to the cloud or run on-premises
- Modular components with minimal overhead retain flexibility while
constructing your solutions

## Prerequisites

This example assumes you already have an ASP.NET Core app
on your machine. If you are new to ASP.NET you can follow a
[simple tutorial](https://www.asp.net/get-started) to initialize a project or
clone our [ASP.NET Docker Sample](https://github.com/dotnet/dotnet-docker/tree/master/samples/aspnetapp).

## Create a Dockerfile for an ASP.NET Core application

### Method 1:

1.  Create a `Dockerfile` in your project folder.
2.  Add the text below to your `Dockerfile` for either Linux or
    [Windows Containers](https://docs.microsoft.com/virtualization/windowscontainers/about/).
    The tags below are multi-arch meaning they pull either Windows or
    Linux containers depending on what mode is set in
    [Docker Desktop for Windows](../docker-for-windows/index.md). Read more on
    [switching containers](../docker-for-windows/index.md#switch-between-windows-and-linux-containers).
3.  The `Dockerfile` assumes that your application is called `aspnetapp`. Change
   the `Dockerfile` to use the DLL file of your project.

```dockerfile
# syntax=docker/dockerfile:1
FROM mcr.microsoft.com/dotnet/sdk:5.0 AS build-env
WORKDIR /app

# Copy csproj and restore as distinct layers
COPY *.csproj ./
RUN dotnet restore

# Copy everything else and build
COPY ../engine/examples ./
RUN dotnet publish -c Release -o out

# Build runtime image
FROM mcr.microsoft.com/dotnet/aspnet:5.0
WORKDIR /app
COPY --from=build-env /app/out .
ENTRYPOINT ["dotnet", "aspnetapp.dll"]
```

4.  To make your build context as small as possible add a [`.dockerignore`
   file](/engine/reference/builder/#dockerignore-file)
   to your project folder and copy the following into it.

```dockerignore
bin/
obj/
```
### Method 2 (build app outside Docker container):

1.  Create a `Dockerfile` in your project folder.
2.  Add the text below to your `Dockerfile` for either Linux or
    [Windows Containers](https://docs.microsoft.com/virtualization/windowscontainers/about/).
    The tags below are multi-arch meaning they pull either Windows or
    Linux containers depending on what mode is set in
    [Docker Desktop for Windows](../docker-for-windows/index.md). Read more on
    [switching containers](../docker-for-windows/index.md#switch-between-windows-and-linux-containers).
3.  The `Dockerfile` assumes that your application is called `aspnetapp`. Change the `Dockerfile` to use the DLL file of your project. This method assumes that your project is already built and it copies the build artifacts from the publish folder. Refer to the Microsoft documentation on [Containerize a .Net Core app](https://docs.microsoft.com/en-us/dotnet/core/docker/build-container?tabs=windows#create-the-dockerfile){: target="blank" rel="noopener" class=“"}.

    The `docker build` step here will be much faster than method 1, as all the artifacts are built outside of the `docker build` step and the size of the base     image is much smaller compared to the build base image.

    This method is preferred for CI tools like Jenkins, Azure DevOps, GitLab CI, etc. as you can use the same artifacts in multiple deployment models if Docker     isn't the only deployment model being used. Additionally, you'll be able to run unit tests and publish code coverage reports, or use custom plugins on the     artifacts built by the CI.

      ```dockerfile
      # syntax=docker/dockerfile:1
      FROM mcr.microsoft.com/dotnet/aspnet:5.0
      COPY bin/Release/netcoreapp3.1/publish/ App/
      WORKDIR /App
      ENTRYPOINT ["dotnet", "aspnetapp.dll"]
      ```
  
4.  To make your build context as small as possible add a [`.dockerignore`
   file](/engine/reference/builder/#dockerignore-file)
   to your project folder.
   
## Build and run the Docker image

1.  Open a command prompt and navigate to your project folder.
2.  Use the following commands to build and run your Docker image:

```console
$ docker build -t aspnetapp .
$ docker run -d -p 8080:80 --name myapp aspnetapp
```

## View the web page running from a container

* Go to [localhost:8080](http://localhost:8080) to access your app in a web browser.
* If you are using the Nano [Windows Container](../docker-for-windows/index.md)
  and have not updated to the Windows Creator Update there is a bug affecting how
  [Windows 10 talks to Containers via "NAT"](https://github.com/Microsoft/Virtualization-Documentation/issues/181#issuecomment-252671828)
  (Network Address Translation). You must hit the IP of the container
  directly. You can get the IP address of your container with the following
  steps:
  1.  Run `docker inspect -f "{% raw %}{{ .NetworkSettings.Networks.nat.IPAddress }}{% endraw %}" myapp`
  2.  Copy the container IP address and paste into your browser.
  (For example, `172.16.240.197`)

## Further reading

  - [ASP.NET Core](https://docs.microsoft.com/aspnet/core/)
  - [Microsoft ASP.NET Core on Docker Hub](https://hub.docker.com/_/microsoft-dotnet-sdk/)
  - [Building Docker Images for ASP.NET Core](https://docs.microsoft.com/aspnet/core/host-and-deploy/docker/building-net-docker-images)
  - [Docker Tools for Visual Studio](https://docs.microsoft.com/dotnet/articles/core/docker/visual-studio-tools-for-docker)
