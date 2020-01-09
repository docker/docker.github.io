---
title: "Containerizing an application"
keywords: containers, images, dockerfiles, node, code, coding, build, push, run
description: Learn how to create a Docker image by writing a Dockerfile, and use it to run a simple container.
---

{% include_relative nav.html selected="2" %}

## Prerequisites

- Work through setup and orientation in [Part 1](index.md).

## Introduction

Now that we've got our orchestrator of choice set up in our development environment thanks to Docker Desktop,
we can begin to develop containerized applications. In general, the development workflow looks like this:

1. Create and test individual containers for each component of your application by first creating Docker images.
2. Assemble your containers and supporting infrastructure into a complete application, expressed either as a *Docker stack file* or in Kubernetes YAML.
3. Test, share and deploy your complete containerized application. 

In this stage of the tutorial, let's focus on step 1 of this workflow: creating the images that our containers will be based on. Remember, a Docker image captures the private filesystem that our containerized processes will run in; we need to create an image that contains just what our application needs to run.

> **Containerized development environments** are easier to set up than traditional development environments, once you learn how to build images as we'll discuss below. This is because a containerized development environment will isolate all the dependencies your app needs inside your Docker image; there's no need to install anything other than Docker on your development machine. In this way, you can easily develop applications for different stacks without changing anything on your development machine.

## Setting Up

1.  Clone an example project from GitHub (if you don't have git installed, see the [install instructions](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) first):

    ```shell
    git clone -b v1 https://github.com/docker-training/node-bulletin-board
    cd node-bulletin-board/bulletin-board-app
    ```

    This is a simple bulletin board application, written in node.js. In this example, let's imagine you wrote this app, and are now trying to containerize it.

2.  Have a look at the file called `Dockerfile`. Dockerfiles describe how to assemble a private filesystem for a container, and can also contain some metadata describing how to run a container based on this image. The bulletin board app Dockerfile looks like this:

    ```dockerfile
    FROM node:6.11.5    

    WORKDIR /usr/src/app
    COPY package.json .
    RUN npm install    
    COPY . .

    CMD [ "npm", "start" ]    
    ```

    Writing a Dockerfile is the first step to containerizing an application. You can think of these Dockerfile commands as a step-by-step recipe on how to build up our image. This one takes the following steps:

    - Start `FROM` the pre-existing `node:6.11.5` image. This is an *official image*, built by the node.js vendors and validated by Docker to be a high-quality image containing the node 6.11.5 interpreter and basic dependencies.
    - Use `WORKDIR` to specify that all subsequent actions should be taken from the directory `/usr/src/app` *in your image filesystem* (never the host's filesystem).
    - `COPY` the file `package.json` from your host to the present location (`.`) in your image (so in this case, to `/usr/src/app/package.json`)
    - `RUN` the command `npm install` inside your image filesystem (which will read `package.json` to determine your app's node dependencies, and install them)
    - `COPY` in the rest of your app's source code from your host to your image filesystem.

    You can see that these are much the same steps you might have taken to set up and install your app on your host - but capturing these as a Dockerfile allows us to do the same thing inside a portable, isolated Docker image.

    The steps above built up the filesystem of our image, but there's one more line in our Dockerfile. The `CMD` directive is our first example of specifying some metadata in our image that describes how to run a container based off of this image. In this case, it's saying that the containerized process that this image is meant to support is `npm start`.

    What you see above is a good way to organize a simple Dockerfile; always start with a `FROM` command, follow it with the steps to build up your private filesystem, and conclude with any metadata specifications. There are many more Dockerfile directives than just the few we see above; for a complete list, see the [Dockerfile reference](https://docs.docker.com/engine/reference/builder/).

## Build and Test Your Image

Now that we have some source code and a Dockerfile, it's time to build our first image, and make sure the containers launched from it work as expected.

> **Windows users**: this example uses Linux containers. Make sure your environment is running Linux containers by right-clicking on the Docker logo in your system tray, and clicking 'Switch to Linux containers...' if the option appears. Don't worry - everything you'll learn in this tutorial works the exact same way for Windows containers.

1.  Make sure you're in the directory `node-bulletin-board/bulletin-board-app` in a terminal or powershell, and build your bulletin board image:

    ```script
    docker image build -t bulletinboard:1.0 .
    ```

    You'll see Docker step through each instruction in your Dockerfile, building up your image as it goes. If successful, the build process should end with a message `Successfully tagged bulletinboard:1.0`.

    > **Windows Users:** you may receive a message titled 'SECURITY WARNING' at this step, noting the read, write and execute permissions being set for files added to your image; we aren't handling any sensitive information in this example, so feel free to disregard this warning in this example.

2.  Start a container based on your new image:

    ```script
    docker container run --publish 8000:8080 --detach --name bb bulletinboard:1.0
    ```

    We used a couple of common flags here:

    - `--publish` asks Docker to forward traffic incoming on the host's port 8000, to the container's port 8080 (containers have their own private set of ports, so if we want to reach one from the network, we have to forward traffic to it in this way; otherwise, firewall rules will prevent all network traffic from reaching your container, as a default security posture).
    - `--detach` asks Docker to run this container in the background.
    - `--name` lets us specify a name with which we can refer to our container in subsequent commands, in this case `bb`.

    Also notice, we didn't specify what process we wanted our container to run. We didn't have to, since we used the `CMD` directive when building our Dockerfile; thanks to this, Docker knows to automatically run the process `npm start` inside our container when it starts up.

3.  Visit your application in a browser at `localhost:8000`. You should see your bulletin board application up and running. At this step, we would normally do everything we could to ensure our container works the way we expected; now would be the time to run unit tests, for example.

4.  Once you're satisfied that your bulletin board container works correctly, delete it:

    ```script
    docker container rm --force bb
    ```
The notation for associating a local image with a repository on a registry is
`username/repository:tag`. The tag is optional, but recommended, since it is
the mechanism that registries use to give Docker images a version. Give the
repository and tag meaningful names for the context, such as
`get-started:part2`. This puts the image in the `get-started` repository and
tags it as `part2`.

## Conclusion


At this point, we've performed a simple containerization of an application, and confirmed that our app runs successfully in its container. The next step will be to write the Kubernetes yaml that describes how to run and manage these containers on Kubernetes which we'll study in Part 3 of this tutorial, or to write the stack file that will let us do the same on Docker Swarm, which we discuss in Part 4. 

[On to Part 3 >>](part3.md){: class="button outline-btn" style="margin-bottom: 30px; margin-right: 100%"}

## CLI References

Further documentation for all CLI commands used in this article are available here:

 - [docker image *](https://docs.docker.com/engine/reference/commandline/image/)
 - [docker container *](https://docs.docker.com/engine/reference/commandline/container/)
 - [Dockerfile reference](https://docs.docker.com/engine/reference/builder/)
