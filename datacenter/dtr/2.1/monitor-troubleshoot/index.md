<!--[metadata]>
+++
aliases = ["/docker-trusted-registry/monitor-troubleshoot/monitor/"]
title = "Monitor DTR"
description = "Learn how to monitor your DTR installation."
keywords = ["docker, registry, monitor, troubleshoot"]
[menu.main]
parent="dtr_menu_monitor_troubleshoot"
identifier="dtr_monitor"
weight=0
+++
<![end-metadata]-->

# Monitor DTR

Docker Trusted Registry is a Dockerized application. To monitor it, you can
use the same tools and techniques you're already using to monitor other
containerized applications running on your cluster. One way to monitor
DTR is using the monitoring capabilities of Docker Universal Control Plane.

In your browser, log in to **Docker Universal Control Plane** (UCP), and
navigate to the **Applications** page.

To make it easier to find DTR, use the search box to **search for the
DTR application**. If you have DTR set up for high-availability, then all the
DTR nodes are displayed.

![](../images/monitor-1.png)

**Click on the DTR application** to see all of its containers, and if they're
running. **Click on a container** to see its details, like configurations,
resources, and logs.

![](../images/monitor-2.png)


## Where to go next

* [Troubleshoot DTR](troubleshoot.md)
* [DTR architecture](../architecture.md)
