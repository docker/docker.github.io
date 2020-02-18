---
title: Docker and iptables
description: The basics of how Docker works with iptables
keywords: network, iptables
---

On Linux, Docker manipulates `iptables` rules to provide network isolation.
This is an implementation detail, and you should not modify the rules Docker
inserts into your `iptables` policies.

## Add iptables policies before Docker's rules

All of Docker's `iptables` rules are added to the `DOCKER` chain. Do not
manipulate this table manually. If you need to add rules which load before
Docker's rules, add them to the `DOCKER-USER` chain. These rules are loaded
before any rules Docker creates automatically.

### Add a DOCKER-USER filter chain to allow persistent rules 
This can be useful if you need to pre-populate `iptables` rules that need to be in place before 
Docker runs. The following example illustrates how rules can be added to the `DOCKER-USER` chain

### Restrict connections to the Docker daemon

By default, all external source IPs are allowed to connect to the Docker daemon.
To allow only a specific IP or network to access the containers, insert a
negated rule at the top of the `DOCKER-USER` filter chain. For example, the following
rule restricts external access to all IP addresses except 192.168.1.1:

```bash
$ iptables -I DOCKER-USER -i ext_if ! -s 192.168.1.1 -j DROP
```

Please note that you will need to change `ext_if` to correspond with your
host's actual external interface. You could instead allow connections from a
source subnet. The following rule only allows access from the subnet 192.168.1.0/24:

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
192.168.1.99 and 10.1.2.3, you can make rules specific to `10.1.2.3` and leave
`192.168.1.99` open.

`iptables` is complicated and more complicated rules are out of scope for this
topic. See the [Netfilter.org HOWTO](https://www.netfilter.org/documentation/HOWTO/NAT-HOWTO.html)
for a lot more information.

### Filtering container traffic
The following example provides a set of filters and uses those filters for container and host traffic: 

```
# WAN = eth0 ; LAN = eth1

# Reset counters
:DOCKER-USER - [0:0]

# Flush
-F DOCKER-USER

# Filters :
## Activate established connexions
-A DOCKER-USER -i eth0 -m conntrack --ctstate RELATED,ESTABLISHED -j RETURN

## Allow all on https/http
-A DOCKER-USER -i eth0 -p tcp -m tcp -m conntrack --ctorigdstport 80 -j RETURN
-A DOCKER-USER -i eth0 -p tcp -m tcp -m conntrack --ctorigdstport 443 -j RETURN

## Allow 8080 from ip
-A DOCKER-USER -i eth0 -p tcp -m tcp -m conntrack --ctorigdstport 8080 -s 10.11.11.0/24 -j RETURN
-A DOCKER-USER -i eth0 -p tcp -m tcp -m conntrack --ctorigdstport 8080 -s 10.22.22.0/24 -j RETURN

# Block all external
-A DOCKER-USER -i eth0 -j DROP
-A DOCKER-USER -j RETURN

COMMIT
```
> **Note**: `--ctorigdstport` matches the destination port on the packet that initiated the connection, 
	not the destination port on the packet being filtered. Therefore, responses to requests from Docker 
	to other servers have `SPT=80`, and match `--ctorigdstport 80`.
	
	For tighter control, all rules allowing the connection should have `--ctdir` added to specifically 
	express their meaning, as shown in the following example:
	
	-A DOCKER-USER -s 1.2.3.4/32 -i eth0 -p tcp -m conntrack --ctorigdstport 80 --ctdir ORIGINAL -j ACCEPT

Load these rules with:
	
	```bash
	$ iptables-restore -n /etc/iptables.conf
	```

## Prevent Docker from manipulating iptables

It is possible to set the `iptables` key to `false` in the Docker engine's configuration file at `/etc/docker/daemon.json`, but this option is not appropriate for most users.  It is not possible to completely prevent Docker from creating `iptables` rules, and creating them after-the-fact is extremely involved and beyond the scope of these instructions. Setting `iptables` to `false` will more than likely break container networking for the Docker engine.

For system integrators who wish to build the Docker runtime into other applications, explore the [`moby` project](https://mobyproject.org/).

## Next steps

- Read [Docker Reference Architecture: Designing Scalable, Portable Docker Container Networks]
(https://success.docker.com/Architecture/Docker_Reference_Architecture%3A_Designing_Scalable%2C_Portable_Docker_Container_Networks)
