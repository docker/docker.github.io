---
description: Overview of Dockerfile commands
keywords: dockerfile
title: Dockerfile reference overview
notoc: true
---
## Description

Directives and commands you can use in a Dockerfile.

| Command | Description |
|-|-|
| [BuildKit](../builder/#buildkit) | Toolkit for converting source code to build artifacts |
| [Parser directives](../builder/#parser-directives) | Affect the way in which subsequent lines in a Dockerfile are handled |
| [· syntax](../builder/#syntax) | Define the location of the Dockerfile syntax that is used to build the Dockerfile |
| [· escape](../builder/#escape) | Set the character used to escape characters in a Dockerfile |
| [Environment replacement](../builder/#environment-replacement) | Environment variable syntax |
| [.dockerignore file](../builder/#.dockerignore-file) | Modify the context to exclude files and directories that match patterns in it |
| [FROM](../builder/#from) | Initialize a new build stage and sets the Base Image for subsequent instructions |
| [RUN](../builder/#run) | Execute any commands in a new layer on top of the current image and commit the results |
| [CMD](../builder/#cmd) | Provide defaults for an executing container |
| [LABEL](../builder/#label) | Add metadata to an image |
| [MAINTAINER (deprecated)](../builder/#maintainer-deprecated) | Set the Author field of the generated images |
| [EXPOSE](../builder/#expose) | Inform Docker that the container listens on the specified network ports at runtime |
| [ENV](../builder/#env) | Set values for environment variables |
| [ADD](../builder/#add) | Copy new files, directories or a remote file URL from context to image |
| [COPY](../builder/#copy) | Copy new files or directories from context to image |
| [ENTRYPOINT](../builder/#entrypoint) | Command that's launched when starting container |
| [VOLUME](../builder/#volume) | Create a mount point with the specified name |
| [USER](../builder/#user) | Set the user name to use when running the image |
| [WORKDIR](../builder/#workdir) | Set the working directory |
| [ARG](../builder/#arg) | Define a variable that users can pass at build-time to the builder |
| [· Default values](../builder/#default-values) | Optional default value for ARG instruction |
| [· Scope](../builder/#scope) | Where a variable definition comes into effect |
| [· Predefined ARGs](../builder/#predefined-args) | Predefined ARG variables |
| [· Platform ARGs](../builder/#automatic-platform-args-in-the-global-scope) | Set of ARG variables with information on the platform |
| [ONBUILD](../builder/#onbuild) | Instruction, executed when the image is used as the base for another build |
| [STOPSIGNAL](../builder/#stopsignal) | Set the system call signal that will be sent to the container to exit |
| [HEALTHCHECK](../builder/#healthcheck) | Tell Docker how to test a container to check that it is still working |
| [SHELL](../builder/#shell) | Default shell used for the shell form of commands |