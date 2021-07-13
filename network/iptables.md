---
title: Docker and iptables
description: The basics of how Docker works with iptables
keywords: network, iptables
---

On Linux, Docker manipulates `iptables` rules to provide network isolation.
While this is an implementation detail and you should not modify the rules
Docker inserts into your `iptables` policies, it does have some implications
on what you need to do if you want to have your own policies in addition to
those managed by Docker.

If you're running Docker on a host that is exposed to the Internet, you will
probably want to have iptables policies in place that prevent unauthorized
access to containers or other services running on your host. This page
describes how to achieve that, and what caveats you need to be aware of.

## Add iptables policies before Docker's rules

Docker installs two custom iptables chains named `DOCKER-USER` and `DOCKER`,
and it ensures that incoming packets are always checked by these two chains
first.

All of Docker's `iptables` rules are added to the `DOCKER` chain. Do not
manipulate this chain manually. If you need to add rules which load before
Docker's rules, add them to the `DOCKER-USER` chain. These rules are applied
before any rules Docker creates automatically.

Rules added to the `FORWARD` chain -- either manually, or by another
iptables-based firewall -- are evaluated _after_ these chains. This means that
if you expose a port through Docker, this port gets exposed no matter what
rules your firewall has configured. If you want those rules to apply even
when a port gets exposed through Docker, you _must_ add these rules to the
`DOCKER-USER` chain.

### Restrict connections to the Docker host

By default, all external source IPs are allowed to connect to the Docker host.
To allow only a specific IP or network to access the containers, insert a
negated rule at the top of the `DOCKER-USER` filter chain. For example, the
following rule restricts external access from all IP addresses except `192.168.1.1`:

```bash
$ iptables -I DOCKER-USER -i ext_if ! -s 192.168.1.1 -j DROP
```

Please note that you will need to change `ext_if` to correspond with your
host's actual external interface. You could instead allow connections from a
source subnet. The following rule only allows access from the subnet `192.168.1.0/24`:

```bash
$ iptables -I DOCKER-USER -i ext_if ! -s 192.168.1.0/24 -j DROP
```

Finally, you can specify a range of IP addresses to accept using `--src-range`
(Remember to also add `-m iprange` when using `--src-range` or `--dst-range`):

```bash
$ iptables -I DOCKER-USER -m iprange -i ext_if ! --src-range 192.168.1.1-192.168.1.3 -j DROP
```

You can combine `-s` or `--src-range` with `-d` or `--dst-range` to control both
the source and destination. For instance, if the Docker daemon listens on both
`192.168.1.99` and `10.1.2.3`, you can make rules specific to `10.1.2.3` and leave
`192.168.1.99` open.

`iptables` is complicated and more complicated rules are out of scope for this
topic. See the [Netfilter.org HOWTO](https://www.netfilter.org/documentation/HOWTO/NAT-HOWTO.html)
for a lot more information.

## Docker on a router

Docker also sets the policy for the `FORWARD` chain to `DROP`. If your Docker
host also acts as a router, this will result in that router not forwarding
any traffic anymore. If you want your system to continue functioning as a
router, you can add explicit `ACCEPT` rules to the `DOCKER-USER` chain to
allow it:

```bash
$ iptables -I DOCKER-USER -i src_if -o dst_if -j ACCEPT
```

## Prevent Docker from manipulating iptables

It is possible to set the `iptables` key to `false` in the Docker engine's configuration file at `/etc/docker/daemon.json`, but this option is not appropriate for most users.  It is not possible to completely prevent Docker from creating `iptables` rules, and creating them after-the-fact is extremely involved and beyond the scope of these instructions. Setting `iptables` to `false` will more than likely break container networking for the Docker engine.

For system integrators who wish to build the Docker runtime into other applications, explore the [`moby` project](https://mobyproject.org/).

## Setting the default bind address for containers

By default, the Docker daemon will expose ports on the `0.0.0.0` address, i.e.
any address on the host. If you want to change that behavior to only
expose ports on an internal IP address, you can use the `--ip` option to
specify a different IP address. However, setting `--ip` only changes the
_default_, it does not _restrict_ services to that IP.

## Integration with Firewalld

If you are running Docker version 20.10.0 or higher with [firewalld](https://firewalld.org){: target="blank" rel="noopener" class=“”} on your system with `--iptables` enabled, Docker automatically creates a `firewalld` zone called `docker` and inserts all the network interfaces it creates (for example, `docker0`) into the `docker` zone to allow seamless networking.

If you manually added the `docker0` network interface to a zone other than `trusted`,
you need to remove the interface from this zone. Follow the steps below fo to remove
the interface from other zones:

```bash
# Please substitute the appropriate zone and docker interface
$ firewall-cmd --get-zone-of-interface=docker0
$ firewall-cmd --zone=trusted --remove-interface=docker0 --permanent
$ firewall-cmd --reload
```

Restarting `dockerd` daemon inserts the interface into the `docker` zone.
