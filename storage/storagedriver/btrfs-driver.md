---
description: Learn how to optimize your use of Btrfs driver.
keywords: 'container, storage, driver, Btrfs '
title: Use the BTRFS storage driver
redirect_from:
- /engine/userguide/storagedriver/btrfs-driver/
---

Btrfs is a next generation copy-on-write filesystem that supports many advanced
storage technologies that make it a good fit for Docker. Btrfs is included in
the mainline Linux kernel.

Docker's `btrfs` storage driver leverages many Btrfs features for image and
container management. Among these features are block-level operations, thin
provisioning, copy-on-write snapshots, and ease of administration. You can
easily combine multiple physical block devices into a single Btrfs filesystem.

This article refers to Docker's Btrfs storage driver as `btrfs` and the overall
Btrfs Filesystem as Btrfs.

> **Note**: The `btrfs` storage driver is only supported on Docker Engine - Community on Ubuntu or Debian.

## Prerequisites

`btrfs` is supported if you meet the following prerequisites:

- **Docker Engine - Community**: For Docker Engine - Community, `btrfs` is only recommended on Ubuntu or Debian.

- Changing the storage driver makes any containers you have already
  created inaccessible on the local system. Use `docker save` to save containers,
  and push existing images to Docker Hub or a private repository, so that you
  not need to re-create them later.

- `btrfs` requires a dedicated block storage device such as a physical disk. This
  block device must be formatted for Btrfs and mounted into `/var/lib/docker/`.
  The configuration instructions below walk you through this procedure. By
  default, the SLES `/` filesystem is formatted with BTRFS, so for SLES, you do
  not need to use a separate block device, but you can choose to do so for
  performance reasons.

- `btrfs` support must exist in your kernel. To check this, run the following
  command:

  ```bash
  $ grep btrfs /proc/filesystems

  btrfs
  ```

- To manage BTRFS filesystems at the level of the operating system, you need the
  `btrfs` command. If you do not have this command, install the `btrfsprogs`
  package (SLES) or `btrfs-tools` package (Ubuntu).

## Configure Docker to use the btrfs storage driver

This procedure is essentially identical on SLES and Ubuntu.

1.  Stop Docker.

2.  Copy the contents of `/var/lib/docker/` to a backup location, then empty
    the contents of `/var/lib/docker/`:

    ```bash
    $ sudo cp -au /var/lib/docker /var/lib/docker.bk
    $ sudo rm -rf /var/lib/docker/*
    ```

3.  Format your dedicated block device or devices as a Btrfs filesystem. This
    example assumes that you are using two block devices called `/dev/xvdf` and
    `/dev/xvdg`. Double-check the block device names because this is a
    destructive operation.

    ```bash
    $ sudo mkfs.btrfs -f /dev/xvdf /dev/xvdg
    ```

    There are many more options for Btrfs, including striping and RAID. See the
    [Btrfs documentation](https://btrfs.wiki.kernel.org/index.php/Using_Btrfs_with_Multiple_Devices).

4.  Mount the new Btrfs filesystem on the `/var/lib/docker/` mount point. You
    can specify any of the block devices used to create the Btrfs filesystem.

    ```bash
    $ sudo mount -t btrfs /dev/xvdf /var/lib/docker
    ```

    Don't forget to make the change permanent across reboots by adding an
    entry to `/etc/fstab`.

5.  Copy the contents of `/var/lib/docker.bk` to `/var/lib/docker/`.

    ```bash
    $ sudo cp -au /var/lib/docker.bk/* /var/lib/docker/
    ```

6.  Change the mode of `/var/lib/docker` to `711`.

    ```bash
    $ sudo chmod 711 /var/lib/docker
    ```

7.  Configure Docker to use the `btrfs` storage driver. This is required even
    though `/var/lib/docker/` is now using a Btrfs filesystem.
    Edit or create the file `/etc/docker/daemon.json`. If it is a new file, add
    the following contents. If it is an existing file, add the key and value
    only, being careful to end the line with a comma if it is not the final
    line before an ending curly bracket (`}`).

    ```json
    {
      "storage-driver": "btrfs"
    }
    ```

    See all storage options for each storage driver in the
    [daemon reference documentation](/engine/reference/commandline/dockerd/#storage-driver-options)

8.  Start Docker. After it is running, verify that `btrfs` is being used as the
    storage driver.

    ```bash
    $ docker info

    Containers: 0
     Running: 0
     Paused: 0
     Stopped: 0
    Images: 0
    Server Version: 17.03.1-ce
    Storage Driver: btrfs
     Build Version: Btrfs v4.4
     Library Version: 101
    <...>
    ```

9.  When you are ready, remove the `/var/lib/docker.bk` directory.

## Manage a Btrfs volume

One of the benefits of Btrfs is the ease of managing Btrfs filesystems without
the need to unmount the filesystem or restart Docker.

When space gets low, Btrfs automatically expands the volume in *chunks* of
roughly 1 GB.

To add a block device to a Btrfs volume, use the `btrfs device add` and
`btrfs filesystem balance` commands.

```bash
$ sudo btrfs device add /dev/svdh /var/lib/docker

$ sudo btrfs filesystem balance /var/lib/docker
```

> **Note**: While you can do these operations with Docker running, performance
> suffers. It might be best to plan an outage window to balance the Btrfs
> filesystem.

## How the `btrfs` storage driver works

The `btrfs` storage driver works differently from `devicemapper` or other
storage drivers in that your entire `/var/lib/docker/` directory is stored on a
Btrfs volume.

### Image and container layers on-disk

Information about image layers and writable container layers is stored in
`/var/lib/docker/btrfs/subvolumes/`. This subdirectory contains one directory
per image or container layer, with the unified filesystem built from a layer
plus all its parent layers. Subvolumes are natively copy-on-write and have space
allocated to them on-demand from an underlying storage pool. They can also be
nested and snapshotted. The diagram below shows 4 subvolumes. 'Subvolume 2' and
'Subvolume 3' are nested, whereas 'Subvolume 4' shows its own internal directory
tree.

![subvolume example](images/btfs_subvolume.jpg)

Only the base layer of an image is stored as a true subvolume. All the other
layers are stored as snapshots, which only contain the differences introduced
in that layer. You can create snapshots of snapshots as shown in the diagram
below.

![snapshots diagram](images/btfs_snapshots.jpg)

On disk, snapshots look and feel just like subvolumes, but in reality they are
much smaller and more space-efficient. Copy-on-write is used to maximize storage
efficiency and minimize layer size, and writes in the container's writable layer
are managed at the block level. The following image shows a subvolume and its
snapshot sharing data.

![snapshot and subvolume sharing data](images/btfs_pool.jpg)

For maximum efficiency, when a container needs more space, it is allocated in
*chunks* of roughly 1 GB in size.

Docker's `btrfs` storage driver stores every image layer and container in its
own Btrfs subvolume or snapshot. The base layer of an image is stored as a
subvolume whereas child image layers and containers are stored as snapshots.
This is shown in the diagram below.

![Btrfs container layers](images/btfs_container_layer.jpg)

The high level process for creating images and containers on Docker hosts
running the `btrfs` driver is as follows:

1.  The image's base layer is stored in a Btrfs *subvolume* under
    `/var/lib/docker/btrfs/subvolumes`.

2.  Subsequent image layers are stored as a Btrfs *snapshot* of the parent
    layer's subvolume or snapshot, but with the changes introduced by this
    layer. These differences are stored at the block level.

3.  The container's writable layer is a Btrfs snapshot of the final image layer,
    with the differences introduced by the running container. These differences
    are stored at the block level.

## How container reads and writes work with `btrfs`

### Reading files

A container is a space-efficient snapshot of an image. Metadata in the snapshot
points to the actual data blocks in the storage pool. This is the same as with
a subvolume. Therefore, reads performed against a snapshot are essentially the
same as reads performed against a subvolume.

### Writing files

- **Writing new files**: Writing a new file to a container invokes an allocate-on-demand
  operation to allocate new data block to the container's snapshot. The file is
  then written to this new space. The allocate-on-demand operation is native to
  all writes with Btrfs and is the same as writing new data to a subvolume. As a
  result, writing new files to a container's snapshot operates at native Btrfs
  speeds.

- **Modifying existing files**: Updating an existing file in a container is a copy-on-write
  operation (*redirect-on-write* is the Btrfs terminology). The original data is
  read from the layer where the file currently exists, and only the modified
  blocks are written into the container's writable layer. Next, the Btrfs driver
  updates the filesystem metadata in the snapshot to point to this new data.
  This behavior incurs very little overhead.

- **Deleting files or directories**: If a container deletes a file or directory
  that exists in a lower layer, Btrfs masks the existence of the file or
  directory in the lower layer. If a container creates a file and then deletes
  it, this operation is performed in the Btrfs filesystem itself and the space
  is reclaimed.

With Btrfs, writing and updating lots of small files can result in slow
performance.

## Btrfs and Docker performance

There are several factors that influence Docker's performance under the `btrfs`
storage driver.

> **Note**: Many of these factors are mitigated by using Docker volumes for
> write-heavy workloads, rather than relying on storing data in the container's
> writable layer. However, in the case of Btrfs, Docker volumes still suffer
> from these draw-backs unless `/var/lib/docker/volumes/` is **not** backed by
> Btrfs.

- **Page caching**. Btrfs does not support page cache sharing. This means that
  each process accessing the same file copies the file into the Docker hosts's
  memory. As a result, the `btrfs` driver may not be the best choice
  high-density use cases such as PaaS.

- **Small writes**. Containers performing lots of small writes (this usage
  pattern matches what happens when you start and stop many containers in a short
  period of time, as well) can lead to poor use of Btrfs chunks. This can
  prematurely fill the Btrfs filesystem and lead to out-of-space conditions on
  your Docker host. Use `btrfs filesys show` to closely monitor the amount of
  free space on your Btrfs device.

- **Sequential writes**. Btrfs uses a journaling technique when writing to disk.
  This can impact the performance of sequential writes, reducing performance by
  up to 50%.

- **Fragmentation**. Fragmentation is a natural byproduct of copy-on-write
  filesystems like Btrfs. Many small random writes can compound this issue.
  Fragmentation can manifest as CPU spikes when using SSDs or head thrashing
  when using spinning disks. Either of these issues can harm performance.

  If your Linux kernel version is 3.9 or higher, you can enable the `autodefrag`
  feature when mounting a Btrfs volume. Test this feature on your own workloads
  before deploying it into production, as some tests have shown a negative
  impact on performance.

- **SSD performance**: Btrfs includes native optimizations for SSD media.
  To enable these features, mount the Btrfs filesystem with the `-o ssd` mount
  option. These optimizations include enhanced SSD write performance by avoiding
  optimization such as *seek optimizations* which do not apply to solid-state
  media.

- **Balance Btrfs filesystems often**: Use operating system utilities such as a
  `cron` job to balance the Btrfs filesystem regularly, during non-peak hours.
  This reclaims unallocated blocks and helps to prevent the filesystem from
  filling up unnecessarily. You cannot rebalance a totally full Btrfs
  filesystem unless you add additional physical block devices to the filesystem.
  See the
  [BTRFS Wiki](https://btrfs.wiki.kernel.org/index.php/Balance_Filters#Balancing_to_fix_filesystem_full_errors).

- **Use fast storage**: Solid-state drives (SSDs) provide faster reads and
  writes than spinning disks.

- **Use volumes for write-heavy workloads**: Volumes provide the best and most
  predictable performance for write-heavy workloads. This is because they bypass
  the storage driver and do not incur any of the potential overheads introduced
  by thin provisioning and copy-on-write. Volumes have other benefits, such as
  allowing you to share data among containers and persisting even when no
  running container is using them.

## Related Information

- [Volumes](../volumes.md)
- [Understand images, containers, and storage drivers](index.md)
- [Select a storage driver](select-storage-driver.md)
