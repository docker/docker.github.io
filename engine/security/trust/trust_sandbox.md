---
description: Play in a trust sandbox
keywords: trust, security, root,  keys, repository, sandbox
title: Play in a content trust sandbox
redirect_from:
- /security/trust/trust_sandbox/
---

This page explains how to set up and use a sandbox for experimenting with trust.
The sandbox allows you to configure and try trust operations locally without
impacting your production images.

Before working through this sandbox, you should have read through the
[trust overview](index.md).

### Prerequisites

These instructions assume you are running in Linux or macOS. You can run
this sandbox on a local machine or on a virtual machine. You need to
have privileges to run docker commands on your local machine or in the VM.

This sandbox requires you to install two Docker tools: Docker Engine >= 1.10.0
and Docker Compose >= 1.6.0. To install the Docker Engine, choose from the
[list of supported platforms](../../install/index.md). To install
Docker Compose, see the
[detailed instructions here](../../../compose/install.md).

## What is in the sandbox?

If you are just using trust out-of-the-box you only need your Docker Engine
client and access to the Docker Hub. The sandbox mimics a
production trust environment, and sets up these additional components.

| Container       | Description                                                                                                                                 |
|-----------------|---------------------------------------------------------------------------------------------------------------------------------------------|
| trustsandbox    | A container with the latest version of Docker Engine and with some preconfigured certificates. This is your sandbox where you can use the `docker` client to test trust operations. |
| Registry server | A local registry service.                                                                                                                 |
| Notary server   | The service that does all the heavy-lifting of managing trust                                                                               |

This means you run your own content trust (Notary) server and registry.
If you work exclusively with the Docker Hub, you would not need these components.
They are built into the Docker Hub for you. For the sandbox, however, you build
your own entire, mock production environment.

Within the `trustsandbox` container, you interact with your local registry rather
than the Docker Hub. This means your everyday image repositories are not used.
They are protected while you play.

When you play in the sandbox, you also create root and repository keys. The
sandbox is configured to store all the keys and files inside the `trustsandbox`
container. Since the keys you create in the sandbox are for play only,
destroying the container destroys them as well.

By using a docker-in-docker image for the `trustsandbox` container, you also
don't pollute your real Docker daemon cache with any images you push and pull.
The images are stored in an anonymous volume attached to this container,
and can be destroyed after you destroy the container.

## Build the sandbox

In this section, you use Docker Compose to specify how to set up and link together
the `trustsandbox` container, the Notary server, and the Registry server.


1. Create a new `trustsandbox` directory and change into it.

        $ mkdir trustsandbox
        $ cd trustsandbox

2. Create a file called `docker-compose.yml` with your favorite editor.  For example, using vim:

        $ touch docker-compose.yml
        $ vim docker-compose.yml

3. Add the following to the new file.

        ```yaml
        version: "2.1"
        services:
          bootstrap:
            image: golang:1.13-alpine
            networks:
              - mdb
            depends_on:
              - mysql
            volumes:
              - ./notary/migrations:/go/migrations
            environment:
              GO111MODULE: 'on'
              MIGRATE_VER: v4.6.2
            command: |-
              sh -c 'go get -tags 'mysql' github.com/golang-migrate/migrate/v4/cli@$${MIGRATE_VER} &&
                mv /go/bin/cli /go/bin/migrate &&
                SERVICE_NAME=notary_signer ./migrations/migrate.sh &&
                SERVICE_NAME=notary_server ./migrations/migrate.sh'

          server:
            image: notary:server
            networks:
              mdb:
              sig:
              sandbox:
                aliases:
                  - notary-server
            ports:
              - "8080"
              - "4443:4443"
            volumes:
              - ./certs:/certs
            depends_on:
              - mysql
              - signer
              - bootstrap
            entrypoint: /usr/bin/env sh
            command: -c 'sleep 21s && notary-server -config server-config.json'

          signer:
            image: notary:signer
            networks:
              mdb:
              sig:
                aliases:
                  - notarysigner
            volumes:
              - ./certs:/certs
            depends_on:
              - mysql
              - bootstrap
            environment:
              NOTARY_SIGNER_DEFAULT_ALIAS: timestamp_1
              NOTARY_SIGNER_TIMESTAMP_1: testpassword
            entrypoint: /usr/bin/env sh
            command: -c 'sleep 20s && notary-signer -config signer-config.json'

          mysql:
            image: mariadb:10.4
            networks:
              - mdb
            volumes:
              - ./notary/notarysql/mysql-initdb.d:/docker-entrypoint-initdb.d
              - notary_data:/var/lib/mysql
            environment:
              TERM: dumb
              MYSQL_ALLOW_EMPTY_PASSWORD: "true"
            command: mysqld --innodb_file_per_table

          registry:
            image: registry:2.7
            environment:
              REGISTRY_HTTP_SECRET: topS3cr3t
            networks:
              - sandbox

          sandbox:
            image: docker:dind
            networks:
              - sandbox
            volumes:
              - ./certs:/certs
            depends_on:
              - server
            privileged: true
            environment:
              # DOCKER_CONTENT_TRUST: "1"
              # DOCKER_CONTENT_TRUST_SERVER: https://notary-server:4443
            entrypoint: /usr/bin/env sh
            command: |-
              -c 'cp /certs/root-ca.crt /usr/local/share/ca-certificates/root-ca.crt &&
                  update-ca-certificates &&
                  dockerd-entrypoint.sh --insecure-registry registry:5000'

          volumes:
            notary_data:
              external: false

          networks:
            mdb:
              external: false
            sig:
              external: false
            sandbox:
              external: false
        ```

4. Save and close the file.

    > **Note:** You can also uncomment the *DOCKER_CONTENT_TRUST* variables to enable Content Trust and point to your local notary straight away.

5. Clone the notary repository

    As the notary `server` and `signer` images need some certificates and a database to be in place we will grab those for brevity from the [notary github repository](https://github.com/theupdateframework/notary).

    > **NOTE:** for production notary servers you should not rely on these certificates.

    Furthermore this repository contains the database migration scripts which are used in our docker-compose setup as well.

        $ git clone https://github.com/theupdateframework/notary.git

6. Copy the test certificates required for notary server

        $ mkdir -p certs
        $ cp notary/fixtures/root-ca.crt certs
        $ cp notary/fixtures/notary-{server,signer}.{key,crt} certs

7. Verify your folder structure and content

    If you followed the steps accordingly you should now have the following folders and files

        $ tree .
        trustsandbox
        ├── ...
        ├── certs
        │   ├── notary-server.crt
        │   ├── notary-server.key
        │   ├── notary-signer.crt
        │   ├── notary-signer.key
        │   └── root-ca.crt
        ├── docker-compose.yml
        └── notary
                ├── ...
                ├── left for brevity...
                ├── ...
                ├── fixtures
                │   ├── ...
                │   ├── notary-server.crt
                │   ├── notary-server.key
                │   ├── notary-signer.crt
                │   ├── notary-signer.key
                │   ├── root-ca.crt
                │   └── ...
                ├── ...
                ├── notarysql
                │   ├── mysql-initdb.d
                │   │   ├── initial-notaryserver.sql
                │   │   └── initial-notarysigner.sql
                │   └── ...
                ├── ...
                └── ...

        103 directories, 356 files

8. Run the containers on your local system.

        $ docker-compose up -d

    The first time you run this, the docker-in-docker, Notary server, and registry
    images are downloaded from Docker Hub.

## Play in the sandbox

Now that everything is setup, you can go into your `sandbox` container and
start testing Docker content trust. From your host machine, obtain a shell
in the `sandbox` container.

    $ docker-compose exec sandbox sh
    / #

### Test some trust operations

Now, pull some images from within the `sandbox` container.

1. Download a `docker` image to test with.

        / # docker pull docker/trusttest
        docker pull docker/trusttest
        Using default tag: latest
        latest: Pulling from docker/trusttest

        b3dbab3810fc: Pull complete
        a9539b34a6ab: Pull complete
        Digest: sha256:d149ab53f8718e987c3a3024bb8aa0e2caadf6c0328f1d9d850b2a2a67f2819a
        Status: Downloaded newer image for docker/trusttest:latest

2. Tag it to be pushed to our sandbox registry:

        / # docker tag docker/trusttest registry:5000/test/trusttest:latest

3. Enable content trust.

        / # export DOCKER_CONTENT_TRUST=1

4. Identify the trust server.

        / # export DOCKER_CONTENT_TRUST_SERVER=https://notary-server:4443

    This step is only necessary because the sandbox is using its own server.
    Normally, if you are using the Docker Public Hub this step isn't necessary.

5. Pull the test image.

        / # docker pull registry:5000/test/trusttest
        Using default tag: latest
        Error: remote trust data does not exist for registry:5000/test/trusttest: notary-server:4443 does not have trust data for registry:5000/test/trusttest

      You see an error, because this content doesn't exist on the `notary-server` yet.

6. Push and sign the trusted image.

        / # docker push registry:5000/test/trusttest:latest
        The push refers to a repository [registry:5000/test/trusttest]
        5f70bf18a086: Pushed
        c22f7bc058a9: Pushed
        latest: digest: sha256:ebf59c538accdf160ef435f1a19938ab8c0d6bd96aef8d4ddd1b379edf15a926 size: 734
        Signing and pushing trust metadata
        You are about to create a new root signing key passphrase. This passphrase
        will be used to protect the most sensitive key in your signing system. Please
        choose a long, complex passphrase and be careful to keep the password and the
        key file itself secure and backed up. It is highly recommended that you use a
        password manager to generate the passphrase and keep it safe. There will be no
        way to recover this key. You can find the key in your config directory.
        Enter passphrase for new root key with ID 27ec255:
        Repeat passphrase for new root key with ID 27ec255:
        Enter passphrase for new repository key with ID 58233f9 (registry:5000/test/trusttest):
        Repeat passphrase for new repository key with ID 58233f9 (registry:5000/test/trusttest):
        Finished initializing "registry:5000/test/trusttest"
        Successfully signed "registry:5000/test/trusttest":latest

    Because you are pushing this repository for the first time, Docker creates
    new root and repository keys and asks you for passphrases with which to
    encrypt them. If you push again after this, it only asks you for repository
    passphrase so it can decrypt the key and sign again.

7. Try pulling the image you just pushed:

        / # docker pull registry:5000/test/trusttest
        Using default tag: latest
        Pull (1 of 1): registry:5000/test/trusttest:latest@sha256:ebf59c538accdf160ef435f1a19938ab8c0d6bd96aef8d4ddd1b379edf15a926
        sha256:ebf59c538accdf160ef435f1a19938ab8c0d6bd96aef8d4ddd1b379edf15a926: Pulling from test/trusttest
        Digest: sha256:ebf59c538accdf160ef435f1a19938ab8c0d6bd96aef8d4ddd1b379edf15a926
        Status: Downloaded newer image for registry:5000/test/trusttest@sha256:ebf59c538accdf160ef435f1a19938ab8c0d6bd96aef8d4ddd1b379edf15a926
        Tagging registry:5000/test/trusttest@sha256:ebf59c538accdf160ef435f1a19938ab8c0d6bd96aef8d4ddd1b379edf15a926 as registry:5000/test/trusttest:latest

8. Inspect the Docker Trust data for the image

        / # docker trust inspect --pretty registry:5000/test/trusttest
        Signatures for registry:5000/test/trusttest

        SIGNED TAG          DIGEST                                                             SIGNERS
        latest              7034d197b82fcb07299fda8b05c91d1601ce64f31bc102b1345d03a2953d210a   (Repo Admin)

        Administrative keys for registry:5000/test/trusttest

        Repository Key: 0c0015a8db104e0138de5af72171100c38df1245c21c2ac1068abb2db6326287
        Root Key:       6ac70a9c3e844f90a74eed5605a2d00ed9bdb3af75813d8fadb338d4c063e442

### Test with malicious images

What happens when data is corrupted and you try to pull it when trust is
enabled? In this section, you go into the `registry` and tamper with some
data. Then, you try and pull it.

1. Leave the `sandbox` shell and container running.

2. Open a new interactive terminal from your host, and obtain a shell into the
    `registry` container.

        $ docker-compose exec registry sh
        root@65084fc6f047:/#

3. Navigate to the docker registry folder.

        / # cd /var/lib/registry/docker/registry/v2/

4. List the layers for the `test/trusttest` image you pushed:

        /var/lib/registry/docker/registry/v2/ # ls -l repositories/test/trusttest/_layers/sha256
        total 12
        drwxr-xr-x    2 root     root          4096 Dec 23 15:24 4f4fb700ef54461cfa02571ae0db9a0dc1e0cdb5577484a6d75e68dc38e8acc1
        drwxr-xr-x    2 root     root          4096 Dec 23 15:24 cc7629d1331a7362b5e5126beb5bf15ca0bf67eb41eab994c719a45de53255cd
        drwxr-xr-x    2 root     root          4096 Dec 23 15:24 fbe94c69c9b964ffe67daa6e1654ab4f1db1dc14decc30f7c03ae43284cfa72

5. Add malicious data to one of the `trusttest` layers:

        root@65084fc6f047:/# vi blobs/sha256/cc/cc7629d1331a7362b5e5126beb5bf15ca0bf67eb41eab994c719a45de53255cd/data

    Change some date in for example the command and save the file.

6. Go back to your `trustsandbox` terminal.

7. List the `trusttest` image.

        / # docker image ls | grep trusttest
        REPOSITORY                     TAG                 IMAGE ID            CREATED             SIZE
        docker/trusttest               latest              cc7629d1331a        11 months ago       5.025 MB
        registry:5000/test/trusttest   latest              cc7629d1331a        11 months ago       5.025 MB
        registry:5000/test/trusttest   <none>              cc7629d1331a        11 months ago       5.025 MB

8. Remove the `trusttest:latest` image from our local cache.

        / # docker image rm -f cc7629d1331a
        Untagged: docker/trusttest:latest
        Untagged: registry:5000/test/trusttest:latest
        Untagged: registry:5000/test/trusttest@sha256:ebf59c538accdf160ef435f1a19938ab8c0d6bd96aef8d4ddd1b379edf15a926
        Deleted: sha256:cc7629d1331a7362b5e5126beb5bf15ca0bf67eb41eab994c719a45de53255cd
        Deleted: sha256:2a1f6535dc6816ffadcdbe20590045e6cbf048d63fd4cc753a684c9bc01abeea
        Deleted: sha256:c22f7bc058a9a8ffeb32989b5d3338787e73855bf224af7aa162823da015d44c

    Docker does not re-download images that it already has cached, but we want
    Docker to attempt to download the tampered image from the registry and reject
    it because it is invalid.

9. Pull the image again. This downloads the image from the registry, because we don't have it cached.

        / # docker pull registry:5000/test/trusttest
        Using default tag: latest
        Pull (1 of 1): registry:5000/test/trusttest:latest@sha256:7034d197b82fcb07299fda8b05c91d1601ce64f31bc102b1345d03a2953d210a
        sha256:7034d197b82fcb07299fda8b05c91d1601ce64f31bc102b1345d03a2953d210a: Pulling from test/trusttest
        fbe94c69c9b9: Downloading [>                                                  ]  36.41kB/2.296MB
        4f4fb700ef54: Download complete
        error pulling image configuration: image config verification failed for digest sha256:cc7629d1331a7362b5e5126beb5bf15ca0bf67eb41eab994c719a45de53255cd

      The pull did not complete because the trust system couldn't verify the
      image.

## More play in the sandbox

Now, you have a full Docker content trust sandbox on your local system,
feel free to play with it and see how it behaves. If you find any security
issues with Docker, feel free to send us an email at <security@docker.com>.

## Clean up your sandbox

When you are done, and want to clean up all the services you've started and any
anonymous volumes that have been created, just run the following command in the
directory where you've created your Docker Compose file:

        $ docker-compose down -v
