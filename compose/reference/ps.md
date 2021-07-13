---
description: Lists containers.
keywords: fig, composition, compose, docker, orchestration, cli,  ps
title: docker-compose ps
notoc: true
---

```none
Usage: ps [options] [SERVICE...]

Options:
    -q, --quiet          Only display IDs
    --services           Display services
    --filter KEY=VAL     Filter services by a property
    -a, --all            Show all stopped containers (including those created by the run command)
```

Lists containers.

```bash
$ docker-compose ps
         Name                        Command                 State             Ports
---------------------------------------------------------------------------------------------
mywordpress_db_1          docker-entrypoint.sh mysqld      Up (healthy)  3306/tcp
mywordpress_wordpress_1   /entrypoint.sh apache2-for ...   Restarting    0.0.0.0:8000->80/tcp
```

Get container ID for service.

```bash
$ docker-compose ps -q db
9dce35adb675f69cc39cc0134fb9a8646fa10923725096fe91b0f0756fa5e06c
```
