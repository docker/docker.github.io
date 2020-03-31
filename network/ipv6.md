---
title: IPv6 with docker overview
description: IPv6 networking with docker, limitations, common bugs and usage examples.
keywords: network, ipv6
---

Before you can use IPv6 in Docker containers or swarm services, you need to
[enable IPv6 support in the Docker daemon](../config/daemon/ipv6.md). Afterward, you can choose to use
either IPv4, IPv6 or both with any container, service or network.
The default bridge must have a IPv6 prefix assigned or docker fails to start.

## With only a native IPv4 connection

If you only have a native IPv4 connection, you can establish a tunnel using [tunnelbroker.net](https://tunnelbroker.net)

1. Register on tunnelbroker.net
2. On the left click on `Create Regular Tunnel`
3. Enter the public IPv4 address of the host instance
4. Click the assignment button at `Routed /48`
5. Go to the `Example Configurations` tab and apply that config to your instance.
6. If you need reverse dns, enter the ip of a server where you want to manage your zone on.
7. Put that copy that /48 address into the fixed-cidr-v6 field, but replace the /48 with /64, so that the default bridge gets some addresses that your containers could pick up.
8. Create custom networks with other subnets of your prefix.

## With a /127 or /128 on the host

You don't have any addresses for your containers, either use the above method to get a prefix, or if you only have one containers:

1. Set `"fixed-cidr-v6": "fe80::/10"` to assign link local addresses to your default bridge. Containers on that bridge can communicate with each other, but not with the internet.
2. Run your container that should have a IPv6 address with `--net=host`

## With a /64 (not recomended)

You could subnet the /64 into /80's, but that could cause problems with applications and routing.
Your provider could detect this also as NDP Exhaustion Attack if you spawn many containers very quickly (e. g. after reboot/restart).
Another possibility is to use the first methode (the same as if you only got a native IPv4), but it could cause your IPv6 default route to change, if your provider has bad IPv6 peering.

### Macvlan

Depending on the network infrastructure of your cloud provider, it could be possible to get another address out of the /64 assigned to your container if you use macvlan instead of a bridge.
If that does not fit your design, or it does not work, you'll have to use NDP Proxy or tunneling a prefix using a tunnel provider instead.

### NDP Proxy

1. Ask your provider if they could assign you a bigger prefix.
2. Change the IPv6 assignment on the interface of the host to a smaller subnet, so that you have other networks to use with the containers.
3. Edit `/etc/sysctl.conf` and add `net.ipv6.conf.all.proxy_ndp = 1`
4. Everytime you launch a docker container you need to add an entry into the NDP Proxy by executing `ip -6 neigh add proxy 2001:db8:1::22 dev eth0` (where 2001:db8:1::22 is the ip docker assigned to the container)
5. You could install a deamon for this [docker-ndp-proxy](https://github.com/pommes/docker-ndp-daemon) (if you have issues, report there)
6. Create new networks with your subnets.

### Known issues

#### Ratelimiting

Providers could rate limit your ndp proxys registrations, as they think you're performing a NDP Exhaustion Attack.

#### Routing issues

Say, you click an instance on FooHosting and FooHosting assigns you a /64 prefix.

The FooHosting router thinks it has a /64 on the one side, with hosts inside it.
When you configure a smaller prefix on your box in order to be able to subnet it,
your host thinks it has a /65 on the interface facing FooHosting and another /65 on some other interface (facing your containers).
Now you try to ping google.
Therefore a packet from a host in your newly created subnet tries to reach out to google.

What happens:

1. The container puts in its source ip and the destination ip of google.
2. That package traverses your host, your host sees, that it should go to google. Looks up its destination ip within its routing table and forwards it to its (default) gateway.
3. The package arrives at the FooHosting router, it sees a package coming from one ip within the /64, so everything is fine, it will process that package as your host did.
4. That all repeats until the package reaches google.
5. Google creates a response and sends that back to your server, basically flipping the source and destination ip in the package.
6. The package after traversing some routers finding its way back through the internet reaches the FooHosting router.
7. The FooHosting router looks into its routing table, thinks that everything within the /64 is just right on that interface. So it looks into its arp table (actually it's NDP, but try to stick to ipv4 namings), sees that there is no entry for that IP.
8. The FooHosting router sends out a arp request.
9. And now nothing happens. So the FooHosting router thinks the package is for a non existing host and fails to deliver it, as it cannot fill in the mac address of the destination.
10. The FooHosting router crafts a icmp message that it was unable to deliver a package and sends that back to google.

To resolve this, check you NDP Proxy setup. If that works, try to ping your own wan ip. If that fails, you need to enable routing within your kernel.
Place `net.ipv6.conf.all.forwarding=1` into `/etc/sysctl.conf`.

## With a prefix bigger than /64

Assuming your provider assigned you `2001:0db80::/48`, your setup is:

1. Change your interface address to `2001:0db80::/64`
2. Add `{ ipv6: true, "fixed-cidr-v6": "2001:0db80:1::/64" }` to your `/etc/docker/daemon.json`, this will assign that address to your default bridge and enable ipv6 support in docker.
3. Create custom networks with your other subnets `2001:0db80:2::/64`, `2001:0db80:3::/64`, ..., `2001:0db80:ffff::/64`.
4. If you don't want to have an IPv6 prefix assigned to your default bridge, change the fixed-cidr-v6 within `/etc/docker/daemon.json` to `fe80::/10`.

## Firewalling

Your containers will allways be reachable by there IPv6 address, even though you did not explicitely export ports.
Therefore secure your containers as if they would allways be run with `--net=host`.
Your Firewall rules need to follow [RFC4890](https://tools.ietf.org/html/rfc4890).

Instructions for using ip6tables can be found [here](https://www.tldp.org/HOWTO/Linux+IPv6-HOWTO/ch18s03.html).

## IPv6 Bootcamp

<div><iframe width="768" height="432" src="https://media.ccc.de/v/froscon2018-2242-ipv6_im_jahre_2018/oembed" frameborder="0" allowfullscreen></iframe><a rel="license" href="http://creativecommons.org/licenses/by/4.0/"><img alt="Creative Commons License" style="border-width:0" src="https://i.creativecommons.org/l/by/4.0/88x31.png" /></a><br /><span xmlns:dct="http://purl.org/dc/terms/" href="http://purl.org/dc/dcmitype/MovingImage" property="dct:title" rel="dct:type">IPv6 im Jahre 2018</span> by <a xmlns:cc="http://creativecommons.org/ns#" href="https://media.ccc.de/v/froscon2018-2242-ipv6_im_jahre_2018" property="cc:attributionName" rel="cc:attributionURL">Falk Stern and Maximilian Wilhelm</a> is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by/4.0/">Creative Commons Attribution 4.0 International License</a></div>
