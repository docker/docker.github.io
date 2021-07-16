---
title: "Configure CI/CD for your application"
keywords: CI/CD, GitHub Actions, Go, Golang, cloud, deployment
description: Learn how to deploy your application to the Cloud.
---

{% include_relative nav.html selected="6" %}

Now, that we have configured a CI/CD pipleline, let's look at how we can deploy the application to cloud. Docker supports deploying containers on Azure ACI and AWS ECS.

## Docker and Azure ACI

The Docker Azure Integration enables developers to use native Docker commands to run applications in Azure Container Instances (ACI) when building cloud-native applications. The new experience provides a tight integration between Docker Desktop and Microsoft Azure allowing developers to quickly run applications using the Docker CLI or VS Code extension, to switch seamlessly from local development to cloud deployment.

For detailed instructions, see [Deploying Docker containers on Azure](/cloud/aci-integration/).

## Docker and AWS ECS

The Docker ECS Integration enables developers to use native Docker commands in Docker Compose CLI to run applications in Amazon EC2 Container Service (ECS) when building cloud-native applications.

The integration between Docker and Amazon ECS allows developers to use the Docker Compose CLI to set up an AWS context in one Docker command, allowing you to switch from a local context to a cloud context and run applications quickly and easily simplify multi-container application development on Amazon ECS using Compose files.

For detailed instructions, see [Deploying Docker containers on ECS](/cloud/ecs-integration/).

## Feedback

Help us improve this topic by providing your feedback. Let us know what you think by creating an issue in the [Docker Docs ](https://github.com/docker/docker.github.io/issues/new?title=[Golang%20docs%20feedback]){:target="_blank" rel="noopener" class="_"} GitHub repository. Alternatively, [create a PR](https://github.com/docker/docker.github.io/pulls){:target="_blank" rel="noopener" class="_"} to suggest updates.

<br />
