---
description: How to install Docker Machine
keywords: machine, orchestration, install, installation, docker, documentation, uninstall Docker Machine, uninstall
title: Install Docker Machine
---

On macOS and Windows, Machine is installed along with other Docker products when
you install the [Docker for Mac](/docker-for-mac/index.md),
[Docker for Windows](/docker-for-windows/index.md), or
[Docker Toolbox](/toolbox/overview.md).

If you want only Docker Machine, you can install the Machine binaries directly
by following the instructions in the next section. You can find the latest
versions of the binaries on the [docker/machine release
page](https://github.com/docker/machine/releases/){: target="_blank" class="_" }
on GitHub.

## Install Machine directly

1.  Install [Docker](/engine/installation/index.md){: target="_blank" class="_" }.

2.  Download the Docker Machine binary and extract it to your PATH.

    If you are running **macOS**:

    ```console
    $ base=https://github.com/docker/machine/releases/download/v{{site.machine_version}} &&
      curl -L $base/docker-machine-$(uname -s)-$(uname -m) >/usr/local/bin/docker-machine &&
      chmod +x /usr/local/bin/docker-machine
    ```

    If you are running **Linux**:

    ```console
    $ base=https://github.com/docker/machine/releases/download/v{{site.machine_version}} &&
      curl -L $base/docker-machine-$(uname -s)-$(uname -m) >/tmp/docker-machine && chmod +x /tmp/docker-machine &&
      sudo mv /tmp/docker-machine /usr/local/bin/docker-machine
    ```

    If you are running **Windows** with [Git BASH](https://git-for-windows.github.io/){: target="_blank" class="_"}:

    ```console
    $ base=https://github.com/docker/machine/releases/download/v{{site.machine_version}} &&
      mkdir -p "$HOME/bin" &&
      curl -L $base/docker-machine-Windows-x86_64.exe > "$HOME/bin/docker-machine.exe" &&
      chmod +x "$HOME/bin/docker-machine.exe"
    ```

    > The above command works on Windows only if you use a
    terminal emulator such as [Git BASH](https://git-for-windows.github.io/){: target="_blank" class="_"}, which supports Linux commands like `chmod`.
    {: .important}

    Otherwise, download one of the releases from the [docker/machine release
    page](https://github.com/docker/machine/releases/){: target="_blank" class="_" } directly.

3.  Check the installation by displaying the Machine version:

        $ docker-machine version
        docker-machine version {{site.machine_version}}, build 9371605

## Install bash completion scripts

The Machine repository supplies several `bash` scripts that add features such
as:

-   command completion
-   a function that displays the active machine in your shell prompt
-   a function wrapper that adds a `docker-machine use` subcommand to switch the
    active machine

Confirm the version and save scripts to `/etc/bash_completion.d` or
`/usr/local/etc/bash_completion.d`:

```bash
base=https://raw.githubusercontent.com/docker/machine/v{{site.machine_version}}
for i in docker-machine-prompt.bash docker-machine-wrapper.bash docker-machine.bash
do
  sudo wget "$base/contrib/completion/bash/${i}" -P /etc/bash_completion.d
done
```

Then you need to run `source
/etc/bash_completion.d/docker-machine-prompt.bash` in your bash
terminal to tell your setup where it can find the file
`docker-machine-prompt.bash` that you previously downloaded.

To enable the `docker-machine` shell prompt, add
`$(__docker_machine_ps1)` to your `PS1` setting in `~/.bashrc`.

```
PS1='[\u@\h \W$(__docker_machine_ps1)]\$ '
```

You can find additional documentation in the comments at the [top of
each
script](https://github.com/docker/machine/tree/master/contrib/completion/bash){:
target="_blank" class="_"}.

## How to uninstall Docker Machine

To uninstall Docker Machine:

*  Optionally, remove the machines you created.

   To remove each machine individually: `docker-machine rm <machine-name>`

   To remove all machines: `docker-machine rm -f $(docker-machine ls
   -q)` (you might need to use `-force` on Windows).

   Removing machines is an optional step because there are cases where
   you might want to save and migrate existing machines to a [Docker
   for Mac](/docker-for-mac/index.md) or [Docker Desktop for
   Windows](/docker-for-windows/index.md) environment, for example.

*  Remove the executable: `rm $(which docker-machine)`


>**Note**: As a point of information, the `config.json`, certificates,
and other data related to each virtual machine created by `docker-machine`
is stored in `~/.docker/machine/machines/` on Mac and Linux and in
`~\.docker\machine\machines\` on Windows. We recommend that you do not edit or
remove those files directly as this only affects information for the Docker
CLI, not the actual VMs, regardless of whether they are local or on remote
servers.

## Where to go next

-  [Docker Machine overview](overview.md)
-  Create and run a Docker host on your [local system using virtualization](get-started.md)
-  Provision multiple Docker hosts [on your cloud provider](get-started-cloud.md)
-  [Docker Machine driver reference](/machine/drivers/index.md)
-  [Docker Machine subcommand reference](/machine/reference/index.md)
