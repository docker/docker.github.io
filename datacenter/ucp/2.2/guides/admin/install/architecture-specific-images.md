---
title: Architecture-specific images
description: Learn how to use images that are specific to particular hardware architectures in Docker Universal Control Plane.
keywords: UCP, Docker Enterprise, image, IBM z, Windows
---

Docker Universal Control Plane deploys images for a number of different
hardware architectures, including IBM z Systems. Some architectures require
pulling images that have specific tags or names indicating the target
architecture.

## Tag for IBM z Systems

Append the string `-s390x` to a UCP system image tag to pull the appropriate
image for IBM z Systems. For example, you can modify the CLI command for getting
a [UCP support dump](..\..\get-support.md) to use an environment variable
that indicates the current architecture:

```bash
{% raw %}[[ $(docker info --format='{{.Architecture}}') == s390x ]] && export _ARCH='-s390x' || export _ARCH=''{% endraw %}

docker container run --rm \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --log-driver none \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }}${_ARCH} \
  support > \
  docker-support-${HOSTNAME}-$(date +%Y%m%d-%H_%M_%S).tgz
```

In this example, the environment variable is named `_ARCH`, but you can use any 
valid shell name.

## OS-specific component names

Some UCP component names depend on the node's operating system. Use the
following table to ensure that you're pulling the right images for each node.

| UCP component base name | Windows name   | IBM z Systems name |
|-------------------------|----------------|--------------------|
| ucp-agent               | ucp-agent-win  | ucp-agent-s390x    |
| ucp-dsinfo              | ucp-dsinfo-win | ucp-dsinfo-s390x   |
