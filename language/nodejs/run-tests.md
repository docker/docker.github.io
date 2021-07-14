---
title: "Run your Tests using Node.js and Mocha frameworks"
keywords: Node.js, build, Mocha, test
description: How to Build and Run your Tests using Node.js and Mocha frameworks
---

{% include_relative nav.html selected="4" %}

## Prerequisites

Work through the steps to build an image and run it as a containerized application in [Use your container for development](develop.md).

## Introduction

Testing is an essential part of modern software development. Testing can mean a lot of things to different development teams. There are unit tests, integration tests and end-to-end testing. In this guide we take a look at running your unit tests in Docker.

## Create a test

Let's define a Mocha test in a `./test` directory within our application.

```shell
$ mkdir -p test
```

Save the following code in `./test/test.js`.

```javascript
var assert = require('assert');
describe('Array', function() {
  describe('#indexOf()', function() {
    it('should return -1 when the value is not present', function() {
      assert.equal([1, 2, 3].indexOf(4), -1);
    });
  });
});
```

### Running locally and testing the application

Let’s build our Docker image and confirm everything is running properly. Run the following command to build and run your Docker image in a container.

```shell
$ docker-compose -f docker-compose.dev.yml up --build
```

Now let’s test our application by POSTing a JSON payload and then make an HTTP GET request to make sure our JSON was saved correctly.

```shell
$ curl --request POST \
  --url http://localhost:8000/test \
  --header 'content-type: application/json' \
  --data '{
	"msg": "testing"
}'
```

Now, perform a GET request to the same endpoint to make sure our JSON payload was saved and retrieved correctly. The “id” and “createDate” will be different for you.

```shell
$ curl http://localhost:8000/test

{"code":"success","payload":[{"msg":"testing","id":"e88acedb-203d-4a7d-8269-1df6c1377512","createDate":"2020-10-11T23:21:16.378Z"}]}
```

## Install Mocha

Run the following command to install Mocha and add it to the developer dependencies:

```shell
$ npm install --save-dev mocha
```

## Update package.json and Dockerfile to run tests

Okay, now that we know our application is running properly, let’s try and run our tests inside of the container. We’ll use the same docker run command we used above but this time, we’ll override the CMD that is inside of our container with npm run test. This will invoke the command that is in the package.json file under the “script” section. See below.

```javascript
{
...
  "scripts": {
    "test": "mocha ./**/*.js",
    "start": "nodemon --inspect=0.0.0.0:9229 server.js"
  },
...
}
```

Below is the Docker command to start the container and run tests:

```shell
$ docker-compose -f docker-compose.dev.yml run notes npm run test
Creating node-docker_notes_run ... 

> node-docker@1.0.0 test /code
> mocha ./**/*.js



  Array
    #indexOf()
      ✓ should return -1 when the value is not present


  1 passing (11ms)

```

### Multi-stage Dockerfile for testing

In addition to running the tests on command, we can run them when we build our image, using a multi-stage Dockerfile. The following Dockerfile will run our tests and build our production image.

```dockerfile
# syntax=docker/dockerfile:1
FROM node:14.15.4 as base

WORKDIR /code

COPY package.json package.json
COPY package-lock.json package-lock.json

FROM base as test
RUN npm ci
COPY . .
CMD [ "npm", "run", "test" ]

FROM base as prod
RUN npm ci --production
COPY . .
CMD [ "node", "server.js" ]
```

We first add a label `as base` to the `FROM node:14.15.4` statement. This allows us to refer to this build stage in other build stages. Next we add a new build stage labeled test. We will use this stage for running our tests.

Now let’s rebuild our image and run our tests. We will run the same docker build command as above but this time we will add the `--target test` flag so that we specifically run the test build stage.

```shell
$ docker build -t node-docker --target test .
Sending build context to Docker daemon  22.35MB
...
Successfully built fa583b97c4dd
Successfully tagged node-docker:latest
```

Now that our test image is built, we can run it in a container and see if our tests pass.

```shell
docker run -it --rm -p 8000:8000 node-docker

> node-docker@1.0.0 test /code
> mocha ./**/*.js



  Array
    #indexOf()
      ✓ should return -1 when the value is not present


  1 passing (12ms)

```

I’ve truncated the build output but you can see that the Mocha test runner completed and all our tests passed.

This is great but at the moment we have to run two docker commands to build and run our tests. We can improve this slightly by using a RUN statement instead of the CMD statement in the test stage. The CMD statement is not executed during the building of the image but is executed when you run the image in a container. While with the RUN statement, our tests will be run during the building of the image and stop the build when they fail.

Update your Dockerfile with the highlighted line below.

```dockerfile
# syntax=docker/dockerfile:1
FROM node:14.15.4 as base

WORKDIR /code

COPY package.json package.json
COPY package-lock.json package-lock.json

FROM base as test
RUN npm ci
COPY . .
RUN npm run test

FROM base as prod
RUN npm ci --production
COPY . .
CMD [ "node", "server.js" ]
```

Now to run our tests, we just need to run the docker build command as above.

```console
$ docker build -t node-docker --target test .
Sending build context to Docker daemon  22.35MB
Step 1/8 : FROM node:14.15.4 as base
 ---> f5be1883c8e0
...
Step 6/8 : RUN npm ci
 ---> Using cache
 ---> bcc5cd4a6a1e
Step 7/8 : COPY . .
 ---> 1528ebcb73fa
Step 8/8 : RUN npm run test
 ---> Running in beadc36b293a

> node-docker@1.0.0 test /code
> mocha ./**/*.js



  Array
    #indexOf()
      ✓ should return -1 when the value is not present


  1 passing (9ms)

Removing intermediate container beadc36b293a
 ---> 445b80e59acd
Successfully built 445b80e59acd
Successfully tagged node-docker:latest
 ```

I’ve truncated the output again for simplicity but you can see that our tests are run and passed. Let’s break one of the tests and observe the output when our tests fail.

Open the test/test.js file and change line 5 as follows.

```shell
     1	var assert = require('assert');
     2	describe('Array', function() {
     3	  describe('#indexOf()', function() {
     4	    it('should return -1 when the value is not present', function() {
     5	      assert.equal([1, 2, 3].indexOf(3), -1);
     6	    });
     7	  });
     8	});
```

Now, run the same docker build command from above and observe that the build fails and the failing testing information is printed to the console.

```console
$ docker build -t node-docker --target test .
Sending build context to Docker daemon  22.35MB
Step 1/8 : FROM node:14.15.4 as base
 ---> 995ff80c793e
...
Step 8/8 : RUN npm run test
 ---> Running in b96d114a336b

> node-docker@1.0.0 test /code
> mocha ./**/*.js



  Array
    #indexOf()
      1) should return -1 when the value is not present


  0 passing (12ms)
  1 failing

  1) Array
       #indexOf()
         should return -1 when the value is not present:

      AssertionError [ERR_ASSERTION]: 2 == -1
      + expected - actual

      -2
      +-1
      
      at Context.<anonymous> (test/test.js:5:14)
      at processImmediate (internal/timers.js:461:21)



npm ERR! code ELIFECYCLE
npm ERR! errno 1
npm ERR! node-docker@1.0.0 test: `mocha ./**/*.js`
npm ERR! Exit status 1
...
```

## Next steps

In this module, we took a look at running tests as part of our Docker image build process.

In the next module, we’ll take a look at how to set up a CI/CD pipeline using GitHub Actions. See:

[Configure CI/CD](configure-ci-cd.md){: .button .primary-btn}

## Feedback

Help us improve this topic by providing your feedback. Let us know what you think by creating an issue in the [Docker Docs](https://github.com/docker/docker.github.io/issues/new?title=[Node.js%20docs%20feedback]){:target="_blank" rel="noopener" class="_"} GitHub repository. Alternatively, [create a PR](https://github.com/docker/docker.github.io/pulls){:target="_blank" rel="noopener" class="_"} to suggest updates.

<br />
