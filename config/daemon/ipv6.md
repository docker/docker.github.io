---
title: Enable IPv6 support
description: How to enable IPv6 support in the Docker daemon
keywords: daemon, network, networking, ipv6
redirect_from:
- /engine/userguide/networking/default_network/ipv6/
---

Before you can use IPv6 in Docker containers or swarm services, you need to
enable IPv6 support in the Docker daemon. Afterward, you can choose to use
either IPv4, IPv6 or both with any container, service, or network. 
The default bridge must have a IPv6 prefix assigned or docker fails to start,
assigning a link local address prefix is enough.

> **Note**: IPv6 networking is only supported on Docker daemons running on Linux
> hosts.
> Also there is a difference in between if you got
> a subnetable prefix 2001:0db8::/48,
> only one network 2001:0db8::/64,
> or only a IPv4 address.
> There is a bug, that docker fails to start if you don't specify a `fixed-cidr-v6`,
> while `ipv6` is specified. It can also be `fe80::/10` for link lokal, if you don't
> have a prefix associated with the default bridge.

1.  Edit `/etc/docker/daemon.json`, set the `ipv6` key to `true` and specify the prefix for the default bridge with `fixed-cidr-v6` (assuming you got 2001:0db8::/48 assigned).

    ```json
    {
      "ipv6": true,
      "fixed-cidr-v6": "2001:0db80:1::/64"
    }
    ```

    Save the file.

2.  Reload the Docker configuration file.

    ```bash
    $ systemctl reload docker
    ```

You can now create networks with the `--ipv6` flag and assign containers IPv6
addresses using the `--ip6` flag.

> **Note**: IPv6 Addresses are reachable publicly, therefore make sure to configure your routes and firewalls.

## Next steps

- [IPv6 with docker overview](/network/ipv6.md)
- [Networking overview](/network/index.md)
- [Container networking](/config/containers/container-networking.md)

