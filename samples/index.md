---
title: Samples
---

{% assign labsbase = "https://github.com/docker/labs/tree/master" %}

## Tutorial labs

Learn how to develop and ship containerized applications, by walking through a
sample that exhibits canonical practices. These labs are from the [Docker Labs
repository]({{ labsbase }}).

| Sample | Description |
| ------ | ----------- |
| [Docker for Beginners]({{ labsbase }}/beginner/){: target="_blank"} | A good "Docker 101" course. |
| [Docker Swarm mode]({{ labsbase}}/swarm-mode){: target="_blank"} | Use Docker for natively managing a cluster of Docker Engines called a swarm. |
| [Configuring developer tools and programming languages]({{ labsbase }}/developer-tools/README.md){: target="_blank"} | How to set-up and use common developer tools and programming languages with Docker. |
| [Live Debugging Java with Docker]({{ labsbase }}/developer-tools/java-debugging){: target="_blank"} | Java developers can use Docker to build a development environment where they can run, test, and live debug code running within a container. |
| [Docker for Java Developers]({{ labsbase }}/developer-tools/java/){: target="_blank"} | Offers Java developers an intro-level and self-paced hands-on workshop with Docker. |
| [Live Debugging a Node.js application in Docker]({{ labsbase }}/developer-tools/nodejs-debugging){: target="_blank"} | Node developers can use Docker to build a development environment where they can run, test, and live debug code running within a container. |
| [Dockerizing a Node.js application]({{ labsbase }}/developer-tools/nodejs/porting/){: target="_blank"} | This tutorial starts with a simple Node.js application and details the steps needed to Dockerize it and ensure its scalability. |
| [Docker for ASP.NET and Windows containers]({{ labsbase }}/windows/readme.md){: target="_blank"} | Docker supports Windows containers, too! Learn how to run ASP.NET, SQL Server, and more in these tutorials. |
| [Docker Security]({{ labsbase }}/security/README.md){: target="_blank"} | How to take advantage of Docker security features. |
| [Building a 12-factor application with Docker]({{ labsbase}}/12factor){: target="_blank"} | Use Docker to create an app that conforms to Heroku's "12 factors for cloud-native applications." |

## Sample applications

Run popular software using Docker.

| Sample | Description |
| ------ | ----------- |
| [apt-cacher-ng](/engine/examples/apt-cacher-ng/) | Run a Dockerized apt-cacher-ng instance. |
| [.Net Core application](/engine/examples/dotnetcore/) | Run a Dockerized ASP.NET Core application. |
| [ASP.NET Core + SQL Server on Linux](/compose/aspnet-mssql-compose/) | Run a Dockerized ASP.NET Core + SQL Server environment. |
| [CouchDB](/engine/examples/couchdb_data_volumes/) | Run a Dockerized CouchDB instance. |
| [Django + PostgreSQL](/compose/django/) | Run a Dockerized Django + PostgreSQL environment. |
| [PostgreSQL](/engine/examples/postgresql_service/) | Run a Dockerized PostgreSQL instance. |
| [Rails + PostgreSQL](/compose/rails/) | Run a Dockerized Rails + PostgreSQL environment. |
| [Riak](/engine/examples/running_riak_service/) | Run a Dockerized Riak instance. |
| [Wordpress](compose/wordpress/) | Run a Dockerized Wordpress environment. |
| [SSHd](/engine/examples/running_ssh_service/) | Run a Dockerized SSHd instance. |

## Library references

The following table provides a list of popular official Docker images. For detailed documentation, select the specific image name. 

| Image name | Description |
| ---------- | ----------- |
{% for page in site.samples %}| [{{ page.title }}](https://hub.docker.com/_/{{ page.title }}){: target="_blank"} | {{ page.description | strip }} |
{% endfor %}
