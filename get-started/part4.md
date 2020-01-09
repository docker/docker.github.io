---
title: "Deploying to Swarm"
keywords: swarm, swarm services, stacks
description: Learn how to describe and deploy a simple application on Docker Swarm.
---
{% include_relative nav.html selected="4" %}

## Prerequisites

- Work through containerizing an application in [Part 2](part2.md).
- Make sure that Swarm is enabled on your Docker Desktop by typing `docker system info`, and looking for a message `Swarm: active` (you might have to scroll up a little).

  If Swarm isn't running, simply type `docker swarm init` at a shell prompt to set it up.

## Introduction

Now that we've demonstrated that the individual components of our application run as stand-alone containers and shown how to deploy it using Kubernetes, let's look at how to arrange for them to be managed by Docker Swarm. Swarm provides many tools for scaling, networking, securing and maintaining your containerized applications, above and beyond the abilities of containers themselves.

In order to validate that our containerized application works well on Swarm, we'll use Docker Desktop's built in Swarm environment right on our development machine to deploy our application, before handing it off to run on a full Swarm cluster in production. The Swarm environment created by Docker Desktop is _fully featured_, meaning it has all the Swarm features your app will enjoy on a real cluster, accessible from the convenience of your development machine.

## Describing Apps Using Stack Files

Swarm never creates individual containers like we did in the previous step of this tutorial; instead, all Swarm workloads are scheduled as _services_, which are scalable groups of containers with added networking features maintained automatically by Swarm. Furthermore, all Swarm objects can and should be described in manifests called _stack files_; these YAML files describe all the components and configurations of your Swarm app, and can be used to easily create and destroy your app in any Swarm environment.

1.  Let's write a simple stack file to run and manage our bulletin board. Place the following in a file called `bb-stack.yaml`:

    ```yaml
    version: '3.7'    

    services:
      bb-app:
        image: bulletinboard:1.0
        ports:
          - "8000:8080"    
    ```


    In this Swarm YAML file, we have just one object: a `service`, describing a scalable group of identical containers. In this case, you'll get just one container (the default), and that container will be based off of your `bulletinboard:1.0` image from step 2 of this tutorial. We've furthermore asked Swarm to forward all traffic arriving at port 8000 on our development machine to port 8080 inside our bulletin board container.

    > **Kubernetes Services and Swarm Services are very different!** Despite the similar name, the two orchestrators mean very different things by the term 'service'. In Swarm, a service provides both scheduling _and_ networking facilities, creating containers and providing tools for routing traffic to them. In Kubernetes, scheduling and networking are handled separately: _deployments_ (or other controllers) handle the scheduling of containers as pods, while _services_ are responsible only for adding networking features to those pods.
  
##  Deploying and Checking Your Application

1.  Deploy your application to Swarm:

    ```shell
    docker stack deploy -c bb-stack.yaml demo
    ```

    If all goes well, Swarm will report creating all your stack objects with no complaints:

    ```shell
    Waiting for the stack to be stable and running...
    bb-app: Ready           [pod status: 1/1 ready, 0/1 pending, 0/1 failed]

    Stack demo is stable and running
    ```

    Notice that in addition to your service, Swarm also creates a Docker network by default to isolate the containers deployed as part of your stack.

2.  Make sure everything worked by showing information about your service:

    ```shell
    docker stack services demo
    ```

    If all has gone well, your service will report with 1/1 of its replicas created:

    ```shell
    ID                  NAME                MODE                REPLICAS            IMAGE               PORTS
    il7elwunymbs        demo_bb-app         replicated          1/1                 bulletinboard:1.0   *:8000->8080/tcp
    ```

    This indicates 1/1 containers you asked for as part of your services are up and running. Also, we see that port 8000 on your development machine is getting forwarded to port 8080 in your bulletin board container.

3.  Open a browser and visit your bulletin board at `localhost:8000`; you should see your bulletin board, the same as when we ran it as a stand-alone container in Step 2 of this tutorial.

4.  Once satisfied, tear down your application:

    ```shell
    docker stack rm demo
    ```

## Conclusion

At this point, we have successfully used Docker Desktop to deploy our application to a fully-featured Swarm environment on our development machine. We haven't done much with Swarm yet, but the door is now open: you can begin adding other components to your app and taking advantage of all the features and power of Swarm, right on your own machine.

In addition to deploying to Swarm, we have also described our application as a stack file. This simple text file contains everything we need to create our application in a running state; we can check it into version control and share it with our colleagues, allowing us to distribute our applications to other clusters (like the testing and production clusters that probably come after our development environments) easily.

[On to Part 5 >>](part5.md){: class="button outline-btn" style="margin-bottom: 30px; margin-right: 100%"}

## Swarm &amp; CLI References

Further documentation for all new Swarm objects and CLI commands used in this article are available here:

 - [Swarm Services](https://docs.docker.com/engine/swarm/how-swarm-mode-works/services/)
 - [Swarm Stacks](https://docs.docker.com/engine/swarm/stack-deploy/)
 - [`docker stack *`](https://docs.docker.com/engine/reference/commandline/stack/)
 - [`docker service *`](https://docs.docker.com/engine/reference/commandline/service/)

