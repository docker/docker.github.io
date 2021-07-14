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

Lists only running containers.

```bash
$ docker-compose ps --services --status=running
mywordpress_db_1
mywordpress_wordpress_1
```

Lists only restarting containers.

```bash
$ docker-compose ps --services --status=restarting
mywordpress_wordpress_1
```

Lists only stopped containers.

```bash
$ docker-compose ps --services --status=stopped 

```

Lists only paused containers.

```bash
$ docker-compose ps --services --status=paused

```
