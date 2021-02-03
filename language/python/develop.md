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

Now we can run MySQL in a container and attach to the volumes and network we created above. Docker pulls the image from Hub and run it for you locally.

```shell
$ docker run -it --rm -d -v mysql:/var/lib/mysql \
  -v mysql_config:/etc/mysql -p 3306:3306 \
  --network mysqlnet \
  --name mysqldb \
  -e MYSQL_ALLOW_EMPTY_PASSWORD=true \
  mysql
```

Now, let’s make sure that our MySQL database is running and that we can connect to it. Connect to the running MySQL database inside the container using the following command:

```shell
$ docker run -it --network mysqlnet --rm mysql mysql -hmysqldb
Enter password: ********

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

In the above command, we used the same MySQL image to connect to the database but this time, we passed the ‘mysql’ command to the container with the `-h` flag containing the name of our MySQL container name. Press CTRL-D to exit the MySQL  interactive terminal.

Next, we'll update the sample application we created in the [Build images](build-images.md#sample-application) module. To see the directory structure of the Python app, see [Python application directory structure](build-images.md#directory-structure).

Okay, now that we have a running MySQL, let’s update the`app.py` to use MySQL as a datastore. Let’s also add some routes to our server. One for fetching records and one for inserting records.

```shell
import mysql.connector
import json
from flask import Flask, request

app = Flask(__name__)

@app.route('/')
def hello_world():
  return 'Hello, Docker!'


@app.route('/widgets', methods = ['POST', 'GET'])
def get_widgets() :
  mydb = mysql.connector.connect(
      host="mysqldb",
      user="root",
      password="",
      database="inventory"
    )
  cursor = mydb.cursor()
  
  # Inserts values into the database
  if request.method == "POST":
    new_widgets = request.form
    query = """
    INSERT INTO widgets (name, description)
    VALUES (%s, %s)
    """
    cursor.execute(query, (new_widgets["name"], new_widgets["description"]))
    mydb.commit()
  
  cursor.execute("SELECT * FROM widgets")
  results = cursor.fetchall() # Extract values stored in widgets table

  row_headers=[x[0] for x in cursor.description] # This will extract row headers

  json_data=[]
  for result in results:
    json_data.append(dict(zip(row_headers,result)))

  cursor.close()
  mydb.close()

  return json.dumps(json_data)


@app.route('/db')
def db_init():
  mydb = mysql.connector.connect(
    host="mysqldb",
    user="root",
    password=""
  )
  cursor = mydb.cursor()

  # Creates database, resets if it already exists
  cursor.execute("DROP DATABASE IF EXISTS inventory")
  cursor.execute("CREATE DATABASE inventory")

  cursor.close()
  mydb.close()

  mydb = mysql.connector.connect(
    host="mysqldb",
    user="root",
    password="",
    database="inventory"
  )

  cursor = mydb.cursor()

  # Creates a table for widgets
  cursor.execute("CREATE TABLE widgets (name VARCHAR(255), description VARCHAR(255))")

  cursor.close()
  mydb.close()

  return 'init database'

if __name__ == "__main__":
  app.run(host ='0.0.0.0')
```

We’ve added the MySQL module and updated the code to connect to the database server, created a database and table. We also created a couple of routes to save widgets and fetch widgets. We now need to rebuild our image so it contains our changes.

First, let’s add the `mysql-connector-python `module to our application using pip.

```shell
$ pip3 install mysql-connector-python
$ pip3 freeze > requirements.txt
```

Now we can build our image.

```shell
$ docker build --tag python-docker .
```

Now, let’s add the container to the database network and then run our container. This allows us to access the database by its container name.

```shell
$ docker run \
  -it --rm -d \
  --network mysqlnet \
  --name rest-server \
  -p 5000:5000 \
  python-docker
```

Let’s test that our application is connected to the database and is able to add a note.

```shell
$ curl http://localhost:5000/initdb
$ curl --request POST \
  --url http://localhost:5000/widgets \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data 'name=widget01' \
  --data 'description=this is a test widget'
```

You should receive the following JSON back from our service.

```shell
[{"name": "widget01", "description": "this is a test widget"}]
```

## Use Compose to develop locally

In this section, we’ll create a Compose file to start our python-docker and the MySQL database using a single command. We’ll also set up the Compose file to start the `python-docker` application in debug mode so that we can connect a debugger to the running process.

Open the `python-docker` code in your IDE or a text editor and create a new file named `docker-compose.dev.yml`. Copy and paste the following commands into the file.

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

 mysqldb:
  image: mysql
  ports:
  - 3306:3306
  environment:
  - MYSQL_ALLOW_EMPTY_PASSWORD=true
  volumes:
  - mysql:/var/lib/mysql
  - mysql_config:/etc/mysql

volumes:
  mysql:
  mysql_config:
```

This Compose file is super convenient as we do not have to type all the parameters to pass to the `docker run` command. We can declaratively do that using a Compose file.

We expose port 5000 so that we can reach the dev web server inside the container. We also map our local source code into the running container to make changes in our text editor and have those changes picked up in the container.

Another really cool feature of using a Compose file is that we have service resolution set up to use the service names. Therefore, we are now able to use “mysqldb” in our connection string. The reason we use “mysqldb” is because that is what we've named our MySQL service as in the Compose file.

Now, to start our application and to confirm that it is running properly, run the following command:

```shell
$ docker-compose -f docker-compose.dev.yml up --build
```

We pass the `--build` flag so Docker will compile our image and then starts the containers.

Now let’s test our API endpoint. Run the following curl command:

```shell
$ curl --request GET --url http://localhost:8080/widgets
```

You should receive the following response:

```shell
[]
```

This is because our database is empty.

## Next steps

In this module, we took a look at creating a general development image that we can use pretty much like our normal command line. We also set up our Compose file to map our source code into the running container and exposed the debugging port.

In the next module, we’ll take a look at how to set up a CI/CD pipeline using GitHub Actions. See:

[Configure CI/CD](configure-ci-cd.md){: .button .outline-btn}

## Feedback

Help us improve this topic by providing your feedback. Let us know what you think by creating an issue in the [Docker Docs](https://github.com/docker/docker.github.io/issues/new?title=[Python%20docs%20feedback]){:target="_blank" rel="noopener" class="_"} GitHub repository. Alternatively, [create a PR](https://github.com/docker/docker.github.io/pulls){:target="_blank" rel="noopener" class="_"} to suggest updates.
