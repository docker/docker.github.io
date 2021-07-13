---
description: Deploy jobs to a swarm
keywords: guide, jobs, swarm mode, swarm
title: Deploy jobs to a swarm
toc_max: 4
---

Docker Swarm jobs, available in Docker Engine 20.10 or newer, provide the
ability for Swarm to support one-off workloads, such as periodic batch
operations. Traditionally, [Swarm Services](services.md) are long
running workloads, defined in the Swarm in a *declarative* model. When using a
traditional Swarm Service, the Swarm will maintain its state, for example
reconciling Tasks to ensuring the number of running Tasks equals the desired
number of Tasks. This is not the case with Swarm jobs, a job will execute until
*completion*, once completed the Task will not be restarted.

There are 2 modes of a Swarm job:

- [Replicated jobs](#replicated-jobs)
- [Global jobs](#global-jobs)
  
Swarm job modes are similar to Swarm services, where by a `ReplicatedJob` is
similar to a `ReplicatedService` and a `GlobalJob` is comparable to a
`GlobalService`. A Replicated job deploys a number of parallel Tasks within the
Swarm cluster, these Tasks could be scheduled anywhere in the cluster assuming
scheduling and resource constraints are met. A Global job deploy a single Task
onto every node in the cluster. Swarm jobs can be attached to existing [Swarm
Overlay Networks](/network/overlay/) and can leverage the Swarm Objects that are
already defined in the cluster, such as [Configs](/engine/swarm/configs/),
[Secrets](/engine/swarm/secrets/) and [Volumes](/storage/volumes/).

## Replicated jobs

A `ReplicatedJob` is a desired number of parallel Tasks that will be scheduled
on to the Swarm. A [Task](/engine/swarm/how-swarm-mode-works/swarm-task-states/)
is a single running container in the cluster, a Task is the scheduled unit for
both Swarm Services and Swarm jobs. However, unlike a Swarm Service, a Swarm job
will not reconcile Tasks. Once a Task has exited successfully, it will not be
rescheduled.

A Replicated job is defined by creating a Service with the mode
`replicated-job`. By default the Swarm will schedule a single Task onto an
available node in the cluster. The Task will stay in the `Running` state until
the container exits successfully (exit code 0), at which point the job will
transition from a `Running` state into a `Completed` state.

```bash
$ docker service create \
    --name sleeper \
    --mode replicated-job  \
    alpine \
    sleep 30

$ docker service ps sleeper
ID             NAME                                IMAGE           NODE            DESIRED STATE   CURRENT STATE            ERROR     PORTS
i3mb2bslik9s   sleeper.cmzlnu1mzp3inydp80pt19d8e   alpine:latest   ip-10-0-2-175   Complete        Complete 2 seconds ago
```

A `ReplicatedJob` can have concurrency with multiple Tasks running the same
workload deployed on to the cluster at the same time. When creating a
`ReplicatedJob` the number of `Replicas` controls its concurrency. Additionally
the number of `Replicas` also sets the number of Tasks that need to be
successfully completed for the Swarm job to move into a `Completed` state.

In this example a Swarm job is created with `--replicas 2`. This will
instruct the Swarm to schedule 2 Tasks on to the Swarm to run in parallel,
assuming there are resources available to do so. Additionally the required
number of Completed Tasks for the Swarm job state to transition from `Running` to
`Completed` is also 2.

```bash
$ docker service create \
    --name concurrent-sleeper \
    --mode replicated-job  \
    --replicas 2 \
    alpine \
    sleep 30

$ docker service inspect concurrent-sleeper | jq -r '.[].Spec.Mode'
{
  "ReplicatedJob": {
    "MaxConcurrent": 2, # Number of Tasks to run in parralel
    "TotalCompletions": 2 # Number of Tasks that need to complete succesfully
  }
}
```

### Maximum Concurrency

The concurrency of Tasks can be managed separately from the required number of
completed tasks, using the `--max-concurrent` flag. When this flag is set the
number of Tasks that are ran in parallel is different from the number of
required completed Tasks. When using the `--max-concurrent` flag, the required
number of completed Tasks is still controlled by the `--replicas` flag.

In this example, the number of required completed Tasks is 4, defined by
`--replicas 4`. However leveraging `--max-concurrent` the Swarm is limited to
schedule only 2 Tasks at a time. The Swarm will wait for those first 2 Tasks to
complete before scheduling the remaining tasks.

```bash
$ docker service create \
    --name max-concurrent-sleeper \
    --mode replicated-job  \
    --max-concurrent 2 \
    --replicas 4 \
    alpine \
    sleep 30
```

Inspecting the Swarm job the concurrency of the Tasks and the required number of
completed Tasks can be seen.

```bash
$ docker service list --filter name=max-concurrent-sleeper
ID             NAME                     MODE             REPLICAS              IMAGE           PORTS
2watf8vnnllr   max-concurrent-sleeper   replicated job   2/2 (0/4 completed)   alpine:latest   

$ docker service inspect max-concurrent-sleeper | jq -r '.[].Spec.Mode'
{
  "ReplicatedJob": {
    "MaxConcurrent": 2, # Number of Tasks to run in parralel
    "TotalCompletions": 4 # Number of Tasks that need to complete succesfully
  }
}
```

### Scheduling Constraints

When using a Swarm jobs placement
[constraints](/engine/swarm/services/#placement-constraints) and
[preferences](/engine/swarm/services/#placement-preferences) can still be used
to support the scheduling of Swarm Tasks. For example to ensure Swarm jobs are
not ran on Nodes labeled `devel` the following `--contraint` can be used.

```bash
$ docker service create \
  --name constrained-sleeper \
  --mode replicated-job \
  --constraint node.labels.type!=devel \
  alpine \
  sleep 30
```

To limit the number of concurrent tasks scheduled on to each Node in a
Replicated job, the `--replicas-max-per-node` flag can be passed when scheduling
a Swarm job. In the following example `--replicas 6` defines the number of
parallel Tasks and the number of completed Tasks, however
`--replicas-max-per-node 2` will ensure the Swarm will only schedule 2 Tasks on
each node.

```bash
$ docker service create \
  --name scheduling-constraints-sleeper \
  --mode replicated-job \
  --replicas 6 \
  --replicas-max-per-node 2 \
  alpine \
  sleep 30
```

### Scaling a services

Once deployed a Swarm job can be scaled through the `docker service scale`. This
will adjust the number of `Completed` Tasks required for a Replicated job, not the
number of concurrent Tasks. 

> **Note** when a Swarm job is scaled the whole job is restarted. All Tasks that
> are currently in the `Running` or `Completed` state will be rerun.

```bash
# Deploy a Swarm job
$ docker service create \
  --name scale-sleeper \
  --mode replicated-job \
  --replicas 2 \
  alpine \
  sleep 30

# Scale the Swarm job to 4 Replicas
$ docker service scale scale-sleeper:4
```

Inspecting this service, shows the concurrency has remained the same, however
the desired number of `Completed` Tasks has increased.

```bash
$ docker service inspect scale-sleeper | jq -r '.[].Spec.Mode'
{
  "ReplicatedJob": {
    "MaxConcurrent": 2,  # Number of Tasks to run in parralel
    "TotalCompletions": 4 # Number of Tasks that need to complete succesfully
  }
}
```

## Global jobs

For global jobs, the scheduler places one task on each available node in the
cluster that meets the job's [placement
constraints](services.md#placement-constraints) and [resource
requirements](services.md#reserve-memory-or-cpus-for-a-service).

```bash
$ docker service create \
    --name global-sleeper \
    --mode global-job  \
    alpine \
    sleep 30
```

Contrary to global _services_, swarm does not update global jobs in the
`running` state when new nodes are added or removed from the cluster. Once a
global job is scheduled, the Swarm will not attempt to add / remove new tasks or
adjust the concurrency when new nodes join the cluster. 

Use `docker service update --force <sevice-name>` to adjust the concurrency of a
global job, and to reschedule a global job on all available nodes in the
cluster.

> **Note**
>
> When a job is updated, all tasks for that job are restarted, including
> tasks that are in the `running` or `completed` state.

## Event or time-based triggers for jobs

Docker does not include a built-in event or time-based trigger for jobs. jobs
are "one-off" tasks, which means that, once completed, they will not rerun, and
will remain in the service list until they are manually removed. At this time,
the triggering of jobs should be implemented outside of the Swarm cluster for
example with tools like Cron or a CI/CD pipeline.

### Example: run a job on a time-based interval with cron

In this example a Global job is initially deployed onto the cluster, the
cron-job will then rerun this Swarm job every 5 minutes, using `docker service
update`. This example assumes `crond` and `crontab` are already installed on
your system, for documentation on getting started with Cron on Ubuntu see this
guide on the [Ubuntu
documentation](https://help.ubuntu.com/community/CronHowto).

1. Create an initial Global job on the cluster.

```bash
$ docker service create \
    --name cron-sleeper \
    --mode global-job  \
    alpine \
    sleep 60
```

2. Open the crontab for the current user and add a new entry to the bottom of
   the file to trigger a rerun of the Swarm job. In this example the Swarm
   Global job will run every 5 minutes.

```console
$ crontab -e

*/5 * * * * /usr/bin/docker service update cron-sleeper --force > /dev/null
```

3. Verify that the cron-job ran and that a new execution of the job has been
   triggered. The last execution time of a job can be seen from the `docker
   service inspect` command.

```bash
$ docker service inspect cron-sleeper | jq -r .[].JobStatus
{
  "JobIteration": {
    "Index": 74
  },
  "LastExecution": "2020-12-18T17:00:01.231624312Z" # Last Run Was 17:00
}

# Wait 5 minutes for the next run of the job to be completed

$ docker service inspect cron-sleeper | jq -r .[].JobStatus
{
  "JobIteration": {
    "Index": 75
  },
  "LastExecution": "2020-12-18T17:05:01.65568426Z" # Last Run Was 17:05
}
```

## Learn More

* [Swarm administration guide](admin_guide.md)
* [Docker Engine command line reference](../reference/commandline/docker.md)
* [Swarm mode tutorial](swarm-tutorial/index.md)