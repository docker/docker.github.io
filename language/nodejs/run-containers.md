---
title: "Run your image as a container"
keywords: get started, Node JS, run, container,
description: Learn how to run the image as a container.
redirect_from:
- /get-started/nodejs/run-containers/
---

{% include_relative nav.html selected="2" %}

## Prerequisites

Work through the steps to build a Node JS image in [Build your Node image](build-images.md).

## Overview

In the previous module we created our sample application and then we created  a Dockerfile that we used to create an image. We created our image using the command `docker build`. Now that we have an image, we can run that image and see if our application is running correctly.

A container is a normal operating system process except that this process is isolated and has its own file system, its own networking, and its own isolated process tree separate from the host.

To run an image inside of a container, we use the `docker run` command. The `docker run` command requires one parameter and that is the image name. Let’s start our image and make sure it is running correctly. Execute the following command in your terminal.

```shell
$ docker run node-docker
```

When you run this command, you’ll notice that you were not returned to the command prompt. This is because our application is a REST server and will run in a loop waiting for incoming requests without return control back to the OS until we stop the container.

Let’s open a new terminal then make a GET request to the server using the curl command.

```shell
$ curl --request POST \
  --url http://localhost:8000/test \
  --header 'content-type: application/json' \
  --data '{
	"msg": "testing"
}'
curl: (7) Failed to connect to localhost port 8000: Connection refused
```

Our curl command failed because the connection to our server was refused. Meaning that we were not able to connect to localhost on port 8000. This is expected because our container is running in isolation which includes networking. Let’s stop the container and restart with port 8000 published on our local network.

To stop the container, press ctrl-c. This will return you to the terminal prompt.

To publish a port for our container, we’ll use the `--publish` flag (`-p` for short) on the docker run command. The format of the `--publish` command is `[host port]:[container port]`. So if we wanted to expose port 8000 inside the container to port 3000 outside the container, we would pass 3000:8000 to the --publish flag.

Start the container and expose port 8000 to port 8000 on the host.

```shell
$ docker run --publish 8000:8000 node-docker
```

Now let’s rerun the curl command from above. Remember to open a new terminal.

```shell
$ curl --request POST \
  --url http://localhost:8000/test \
  --header 'content-type: application/json' \
  --data '{"msg": "testing"}'
{"code":"success","payload":[{"msg":"testing","id":"dc0e2c2b-793d-433c-8645-b3a553ea26de","createDate":"2020-09-01T17:36:09.897Z"}]}
```

Success! We were able to connect to the application running inside of our container on port 8000. Switch back to the terminal where your container is running and you should see the POST request logged to the console.

`2020-09-01T17:36:09:8770  INFO: POST /test`

Press ctrl-c to stop the container.

## Run in detached mode

This is great so far, but our sample application is a web server and we should not have to have our terminal connected to the container. Docker can run your container in detached mode or in the background. To do this, we can use the `--detach` or `-d` for short. Docker will start your container the same as before but this time will “detach” from the container and return you to the terminal prompt.

```shell
$ docker run -d -p 8000:8000 node-docker
ce02b3179f0f10085db9edfccd731101868f58631bdf918ca490ff6fd223a93b
```

Docker started our container in the background and printed the Container ID on the terminal.

Again, let’s make sure that our container is running properly. Run the same curl command from above.

```shell
$ curl --request POST \
  --url http://localhost:8000/test \
  --header 'content-type: application/json' \
  --data '{
	"msg": "testing"
}'
{"code":"success","payload":[{"msg":"testing","id":"dc0e2c2b-793d-433c-8645-b3a553ea26de","createDate":"2020-09-01T17:36:09.897Z"}]}
```

## List containers

Since we ran our container in the background, how do we know if our container is running or what other containers are running on our machine? Well, we can run the `docker ps` command. Just like on Linux, to see a list of processes on your machine we would run the ps command. In the same spirit, we can run the `docker ps` command which will show us a list of containers running on our machine.

```
$ docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                    NAMES
ce02b3179f0f        node-docker         "docker-entrypoint.s…"   6 minutes ago       Up 6 minutes        0.0.0.0:8000->8000/tcp   wonderful_kalam
```

The `ps` command tells a bunch of stuff about our running containers. We can see the Container ID, The image running inside the container, the command that was used to start the container, when it was created, the status, ports that exposed and the name of the container.

You are probably wondering where the name of our container is coming from. Since we didn’t provide a name for the container when we started it, Docker generated a random name. We’ll fix this in a minute but first we need to stop the container. To stop the container, run the `docker stop` command which does just that, stops the container. You will need to pass the name of the container or you can use the container id.

```shell
$ docker stop wonderful_kalam
wonderful_kalam
```

Now rerun the `docker ps` command to see a list of running containers.

```shell
$ docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
```

## Stop, start, and name containers

Docker containers can be started, stopped and restarted. When we stop a container, it is not removed but the status is changed to stopped and the process inside of the container is stopped. When we ran the `docker ps` command, the default output is to only show running containers. If we pass the `--all` or `-a` for short, we will see all containers on our system whether they are stopped or started.

```shell
$ docker ps -a
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS                      PORTS               NAMES
ce02b3179f0f        node-docker         "docker-entrypoint.s…"   16 minutes ago      Exited (0) 5 minutes ago                        wonderful_kalam
ec45285c456d        node-docker         "docker-entrypoint.s…"   28 minutes ago      Exited (0) 20 minutes ago                       agitated_moser
fb7a41809e5d        node-docker         "docker-entrypoint.s…"   37 minutes ago      Exited (0) 36 minutes ago                       goofy_khayyam
```

If you’ve been following along, you should see several containers listed. These are containers that we started and stopped but have not been removed.

Let’s restart the container that we just stopped. Locate the name of the container we just stopped and replace the name of the container below in the restart command.

```shell
$ docker restart wonderful_kalam
```

Now, list all the containers again using the ps command.

```shell
$ docker ps --all
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS                      PORTS                    NAMES
ce02b3179f0f        node-docker         "docker-entrypoint.s…"   19 minutes ago      Up 8 seconds                0.0.0.0:8000->8000/tcp   wonderful_kalam
ec45285c456d        node-docker         "docker-entrypoint.s…"   31 minutes ago      Exited (0) 23 minutes ago                            agitated_moser
fb7a41809e5d        node-docker         "docker-entrypoint.s…"   40 minutes ago      Exited (0) 39 minutes ago                            goofy_khayyam
```

Notice that the container we just restarted has been started in detached mode and has port 8000 exposed. Also, observe the status of the container is “Up X seconds”. When you restart a container, it will be started with the same flags or commands that it was originally started with.

Let’s stop and remove all of our containers and take a look at fixing the random naming issue.

Stop the container we just started. Find the name of your running container and replace the name in the command below with the name of the container on your system.

```shell
$ docker stop wonderful_kalam
wonderful_kalam
```

Now that all of our containers are stopped, let’s remove them. When a container is removed, it is no longer running nor is it in the stopped status. However, the process inside the container has been stopped and the metadata for the container has been removed.

```shell
$ docker ps --all
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS                      PORTS                    NAMES
ce02b3179f0f        node-docker         "docker-entrypoint.s…"   19 minutes ago      Up 8 seconds                0.0.0.0:8000->8000/tcp   wonderful_kalam
ec45285c456d        node-docker         "docker-entrypoint.s…"   31 minutes ago      Exited (0) 23 minutes ago                            agitated_moser
fb7a41809e5d        node-docker         "docker-entrypoint.s…"   40 minutes ago      Exited (0) 39 minutes ago                            goofy_khayyam
```

To remove a container, simple run the `docker rm` command passing the container name. You can pass multiple container names to the command in one command.

Again, make sure you replace the containers names in the below command with the container names from your system.

```shell
$ docker rm wonderful_kalam agitated_moser goofy_khayyam
wonderful_kalam
agitated_moser
goofy_khayyam
```

Run the `docker ps --all` command again to see that all containers are gone.

Now let’s address the pesky random name issue. Standard practice is to name your containers for the simple reason that it is easier to identify what is running in the container and what application or service it is associated with. Just like good naming conventions for variables in your code makes it simpler to read. So goes naming your containers.

To name a container, we just need to pass the `--name` flag to the run command.

```shell
$ docker run -d -p 8000:8000 --name rest-server node-docker
1aa5d46418a68705c81782a58456a4ccdb56a309cb5e6bd399478d01eaa5cdda
$ docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                    NAMES
1aa5d46418a6        node-docker         "docker-entrypoint.s…"   3 seconds ago       Up 3 seconds        0.0.0.0:8000->8000/tcp   rest-server
```

Now, we can easily identify our container based on the name.

## Next steps

In this module, we took a look at running containers, publishing ports, and running containers in detached mode. We also took a look at managing containers by starting, stopping, and restarting them. We also looked at naming our containers so they are more easily identifiable. In the next module, we’ll learn how to run a database in a container and connect it to our application. See:

[How to develop your application](develop.md){: .button .primary-btn}

## Feedback

Help us improve this topic by providing your feedback. Let us know what you think by creating an issue in the [Docker Docs](https://github.com/docker/docker.github.io/issues/new?title=[Node.js%20docs%20feedback]){:target="_blank" rel="noopener" class="_"} GitHub repository. Alternatively, [create a PR](https://github.com/docker/docker.github.io/pulls){:target="_blank" rel="noopener" class="_"} to suggest updates.

<br />
