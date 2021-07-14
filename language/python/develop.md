---
title: "Use containers for development"
keywords: python, local, development, run,
description: Learn how to develop your application locally.
---

{% include_relative nav.html selected="3" %}

## Prerequisites

Work through the steps to build an image and run it as a containerized application in [Run your image as a container](run-containers.md).

## Introduction

In this module, we’ll walk through setting up a local development environment for the application we built in the previous modules. We’ll use Docker to build our images and Docker Compose to make everything a whole lot easier.

## Run a database in a container

First, we’ll take a look at running a database in a container and how we use volumes and networking to persist our data and allow our application to talk with the database. Then we’ll pull everything together into a Compose file which allows us to setup and run a local development environment with one command. Finally, we’ll take a look at connecting a debugger to our application running inside a container.

Instead of downloading MySQL, installing, configuring, and then running the MySQL database as a service, we can use the Docker Official Image for MySQL and run it in a container.

Before we run MySQL in a container, we'll create a couple of volumes that Docker can manage to store our persistent data and configuration. Let’s use the managed volumes feature that Docker provides instead of using bind mounts. You can read all about [Using volumes](../../storage/volumes.md) in our documentation.

Let’s create our volumes now. We’ll create one for the data and one for configuration of MySQL.

```shell
$ docker volume create mysql
$ docker volume create mysql_config
```

Now we’ll create a network that our application and database will use to talk to each other. The network is called a user-defined bridge network and gives us a nice DNS lookup service which we can use when creating our connection string.

```shell
$ docker network create mysqlnet
```

Now we can run MySQL in a container and attach to the volumes and network we created above. Docker pulls the image from Hub and runs it for you locally.
In the following command, option `-v` is for starting the container with volumes. For more information, see [Docker volumes](../../storage/volumes.md).

```shell
$ docker run --rm -d -v mysql:/var/lib/mysql \
  -v mysql_config:/etc/mysql -p 3306:3306 \
  --network mysqlnet \
  --name mysqldb \
  -e MYSQL_ROOT_PASSWORD=p@ssw0rd1 \
  mysql
```

Now, let’s make sure that our MySQL database is running and that we can connect to it. Connect to the running MySQL database inside the container using the following command and enter "p@ssw0rd1" when prompted for the password:

```shell
$ docker exec -ti mysqldb mysql -u root -p
Enter password:
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 8
Server version: 8.0.23 MySQL Community Server - GPL

Copyright (c) 2000, 2021, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql>
```

### Connect the application to the database

In the above command, we logged in to the MySQL database by passing the ‘mysql’ command to the `mysqldb` container. Press CTRL-D to exit the MySQL interactive terminal.

Next, we'll update the sample application we created in the [Build images](build-images.md#sample-application) module. To see the directory structure of the Python app, see [Python application directory structure](build-images.md#directory-structure).

Okay, now that we have a running MySQL, let’s update the `app.py` to use MySQL as a datastore. Let’s also add some routes to our server. One for fetching records and one for inserting records.

```shell
import mysql.connector
import json
from flask import Flask

app = Flask(__name__)

@app.route('/')
def hello_world():
  return 'Hello, Docker!'

@app.route('/widgets')
def get_widgets() :
  mydb = mysql.connector.connect(
    host="mysqldb",
    user="root",
    password="p@ssw0rd1",
    database="inventory"
  )
  cursor = mydb.cursor()


  cursor.execute("SELECT * FROM widgets")

  row_headers=[x[0] for x in cursor.description] #this will extract row headers

  results = cursor.fetchall()
  json_data=[]
  for result in results:
    json_data.append(dict(zip(row_headers,result)))

  cursor.close()

  return json.dumps(json_data)

@app.route('/initdb')
def db_init():
  mydb = mysql.connector.connect(
    host="mysqldb",
    user="root",
    password="p@ssw0rd1"
  )
  cursor = mydb.cursor()

  cursor.execute("DROP DATABASE IF EXISTS inventory")
  cursor.execute("CREATE DATABASE inventory")
  cursor.close()

  mydb = mysql.connector.connect(
    host="mysqldb",
    user="root",
    password="p@ssw0rd1",
    database="inventory"
  )
  cursor = mydb.cursor()

  cursor.execute("DROP TABLE IF EXISTS widgets")
  cursor.execute("CREATE TABLE widgets (name VARCHAR(255), description VARCHAR(255))")
  cursor.close()

  return 'init database'

if __name__ == "__main__":
  app.run(host ='0.0.0.0')
```

We’ve added the MySQL module and updated the code to connect to the database server, created a database and table. We also created a couple of routes to save widgets and fetch widgets. We now need to rebuild our image so it contains our changes.

First, let’s add the `mysql-connector-python` module to our application using pip.

```shell
$ pip3 install mysql-connector-python
$ pip3 freeze > requirements.txt
```

Now we can build our image.

```shell
$ docker build --tag python-docker-dev .
```

Now, let’s add the container to the database network and then run our container. This allows us to access the database by its container name.

```shell
$ docker run \
  --rm -d \
  --network mysqlnet \
  --name rest-server \
  -p 5000:5000 \
  python-docker-dev
```

Let’s test that our application is connected to the database and is able to add a note.

```shell
$ curl http://localhost:5000/initdb
$ curl http://localhost:5000/widgets
```

You should receive the following JSON back from our service.

```shell
[]
```

## Use Compose to develop locally

In this section, we’ll create a [Compose file](../../compose/index.md) to start our python-docker and the MySQL database using a single command. We’ll also set up the Compose file to start the `python-docker-dev` application in debug mode so that we can connect a debugger to the running process.

Open the `python-docker` directory in your IDE or a text editor and create a new file named `docker-compose.dev.yml`. Copy and paste the following commands into the file.

```yaml
version: '3.8'

services:
 web:
  build:
   context: .
  ports:
  - 5000:5000
  volumes:
  - ./:/app
  networks: 
  - mysqlnet

 mysqldb:
  image: mysql
  ports:
  - 3306:3306
  environment:
  - MYSQL_ROOT_PASSWORD=p@ssw0rd1
  volumes:
  - mysql:/var/lib/mysql
  - mysql_config:/etc/mysql
  networks: 
  - mysqlnet

volumes:
  mysql:
  mysql_config:

networks: 
  mysqlnet:
```

This Compose file is super convenient as we do not have to type all the parameters to pass to the `docker run` command. We can declaratively do that using a Compose file.

We expose port 5000 so that we can reach the dev web server inside the container. We also map our local source code into the running container to make changes in our text editor and have those changes picked up in the container.

Another really cool feature of using a Compose file is that we have service resolution set up to use the service names. Therefore, we are now able to use “mysqldb” in our connection string. The reason we use “mysqldb” is because that is what we've named our MySQL service as in the Compose file.

Now, to start our application and to confirm that it is running properly, run the following command:

```shell
$ docker-compose -f docker-compose.dev.yml up --build
```

We pass the `--build` flag so Docker will compile our image and then starts the containers.

Now let’s test our API endpoint. Open a new terminal then make a GET request to the server using the curl commands:

```shell
$ curl http://localhost:5000/initdb
$ curl http://localhost:5000/widgets
```

You should receive the following response:

```shell
[]
```

This is because our database is empty.

## Next steps

In this module, we took a look at creating a general development image that we can use pretty much like our normal command line. We also set up our Compose file to map our source code into the running container and exposed the debugging port.

In the next module, we’ll take a look at how to set up a CI/CD pipeline using GitHub Actions. See:

[Configure CI/CD](configure-ci-cd.md){: .button .primary-btn}

## Feedback

Help us improve this topic by providing your feedback. Let us know what you think by creating an issue in the [Docker Docs](https://github.com/docker/docker.github.io/issues/new?title=[Python%20docs%20feedback]){:target="_blank" rel="noopener" class="_"} GitHub repository. Alternatively, [create a PR](https://github.com/docker/docker.github.io/pulls){:target="_blank" rel="noopener" class="_"} to suggest updates.
