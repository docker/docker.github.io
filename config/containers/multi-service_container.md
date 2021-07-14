---
description: How to run more than one process in a container
keywords: docker, supervisor, process management
redirect_from:
- /engine/articles/using_supervisord/
- /engine/admin/using_supervisord/
- /engine/admin/multi-service_container/
title: Run multiple services in a container
---

A container's main running process is the `ENTRYPOINT` and/or `CMD` at the
end of the `Dockerfile`. It is generally recommended that you separate areas of
concern by using one service per container. That service may fork into multiple
processes (for example, Apache web server starts multiple worker processes).
It's ok to have multiple processes, but to get the most benefit out of Docker,
avoid one container being responsible for multiple aspects of your overall
application. You can connect multiple containers using user-defined networks and
shared volumes.

The container's main process is responsible for managing all processes that it
starts. In some cases, the main process isn't well-designed, and doesn't handle
"reaping" (stopping) child processes gracefully when the container exits. If
your process falls into this category, you can use the `--init` option when you
run the container. The `--init` flag inserts a tiny init-process into the
container as the main process, and handles reaping of all processes when the
container exits.

Handling such processes this way is superior to using a full-fledged init
process such as `sysvinit`, `upstart`, or `systemd` to handle process
lifecycles within your container. These "traditional" init systems are
much more heavyweight - thereby increasing the size of your container images -
and may not work as expected inside a containerized environment.

If you need to run more than one service within a container, you can accomplish
this in a few different ways.

- Put all of your commands in a wrapper script, complete with testing and
  debugging information. Run the wrapper script as your `CMD`. This is a very
  naive example. First, the wrapper script:

  ```bash
  #!/bin/bash

  # Start the first process
  ./my_first_process -D
  status=$?
  if [ $status -ne 0 ]; then
    echo "Failed to start my_first_process: $status"
    exit $status
  fi

  # Start the second process
  ./my_second_process -D
  status=$?
  if [ $status -ne 0 ]; then
    echo "Failed to start my_second_process: $status"
    exit $status
  fi

  # Naive check runs checks once a minute to see if either of the processes exited.
  # This illustrates part of the heavy lifting you need to do if you want to run
  # more than one service in a container. The container exits with an error
  # if it detects that either of the processes has exited.
  # Otherwise it loops forever, waking up every 60 seconds

  while sleep 60; do
    ps aux |grep my_first_process |grep -q -v grep
    PROCESS_1_STATUS=$?
    ps aux |grep my_second_process |grep -q -v grep
    PROCESS_2_STATUS=$?
    # If the greps above find anything, they exit with 0 status
    # If they are not both 0, then something is wrong
    if [ $PROCESS_1_STATUS -ne 0 -o $PROCESS_2_STATUS -ne 0 ]; then
      echo "One of the processes has already exited."
      exit 1
    fi
  done
  ```

  Next, the Dockerfile:

  ```dockerfile
  # syntax=docker/dockerfile:1
  FROM ubuntu:latest
  COPY my_first_process my_first_process
  COPY my_second_process my_second_process
  COPY my_wrapper_script.sh my_wrapper_script.sh
  CMD ./my_wrapper_script.sh
  ```

- If you have one main process that needs to start first and stay running but
  you temporarily need to run some other processes (perhaps to interact with
  the main process) then you can use bash's job control to facilitate that.
  First, the wrapper script:

  ```bash
  #!/bin/bash
  
  # turn on bash's job control
  set -m
  
  # Start the primary process and put it in the background
  ./my_main_process &
  
  # Start the helper process
  ./my_helper_process
  
  # the my_helper_process might need to know how to wait on the
  # primary process to start before it does its work and returns
  
  
  # now we bring the primary process back into the foreground
  # and leave it there
  fg %1
  ```

  ```dockerfile
  # syntax=docker/dockerfile:1
  FROM ubuntu:latest
  COPY my_main_process my_main_process
  COPY my_helper_process my_helper_process
  COPY my_wrapper_script.sh my_wrapper_script.sh
  CMD ./my_wrapper_script.sh
  ```

- Use a process manager like `supervisord`. This is a moderately heavy-weight
  approach that requires you to package `supervisord` and its configuration in
  your image (or base your image on one that includes `supervisord`), along with
  the different applications it manages. Then you start `supervisord`, which
  manages your processes for you. Here is an example Dockerfile using this
  approach, that assumes the pre-written `supervisord.conf`, `my_first_process`,
  and `my_second_process` files all exist in the same directory as your
  Dockerfile.

  ```dockerfile
  # syntax=docker/dockerfile:1
  FROM ubuntu:latest
  RUN apt-get update && apt-get install -y supervisor
  RUN mkdir -p /var/log/supervisor
  COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
  COPY my_first_process my_first_process
  COPY my_second_process my_second_process
  CMD ["/usr/bin/supervisord"]
  ```
