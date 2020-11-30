<!-- This file is included in Docker Engine - Community or EE installation docs for Linux. -->

### Install using the convenience script

Docker provides convenience scripts at [get.docker.com](https://get.docker.com/)
and [test.docker.com](https://test.docker.com/) for installing edge and
testing versions of Docker Engine - Community into development environments quickly and
non-interactively. The source code for the scripts is in the
[`docker-install` repository](https://github.com/docker/docker-install).
**Using these scripts is not recommended for production
environments**, and you should understand the potential risks before you use
them:

- The scripts require `root` or `sudo` privileges to run. Therefore,
  you should carefully examine and audit the scripts before running them.
- The scripts attempt to detect your Linux distribution and version and
  configure your package management system for you. In addition, the scripts do
  not allow you to customize any installation parameters. This may lead to an
  unsupported configuration, either from Docker's point of view or from your own
  organization's guidelines and standards.
- The scripts install all dependencies and recommendations of the package
  manager without asking for confirmation. This may install a large number of
  packages, depending on the current configuration of your host machine.
- The script does not provide options to specify which version of Docker to install,
  and installs the latest version that is released in the "edge" channel.
- Do not use the convenience script if Docker has already been installed on the
  host machine using another mechanism.

This example uses the script at [get.docker.com](https://get.docker.com/) to
install the latest release of Docker Engine - Community on Linux. To install the latest
testing version, use [test.docker.com](https://test.docker.com/) instead. In
each of the commands below, replace each occurrence of `get` with `test`.

> **Warning**:
>
Always examine scripts downloaded from the internet before
> running them locally.
{:.warning}

```bash
$ curl -fsSL https://get.docker.com -o get-docker.sh
$ sudo sh get-docker.sh

<output truncated>
```

If you would like to use Docker as a non-root user, you should now consider
adding your user to the "docker" group with something like:

```bash
  sudo usermod -aG docker $USER
```

Remember to log out and back in for this to take effect!

> **Warning**:
>
> Adding a user to the "docker" group grants them the ability to run containers
> which can be used to obtain root privileges on the Docker host. Refer to
> [Docker Daemon Attack Surface](/engine/security/#docker-daemon-attack-surface)
> for more information.
{:.warning}

Docker Engine - Community is installed. It starts automatically on `DEB`-based distributions. On
`RPM`-based distributions, you need to start it manually using the appropriate
`systemctl` or `service` command. As the message indicates, non-root users can't
run Docker commands by default.

> **Note**:
>
> To install Docker without root privileges, see
> [Run the Docker daemon as a non-root user (Rootless mode)](/engine/security/rootless/).
>
> Rootless mode is currently available as an experimental feature.

#### Upgrade Docker after using the convenience script

If you installed Docker using the convenience script, you should upgrade Docker
using your package manager directly. There is no advantage to re-running the
convenience script, and it can cause issues if it attempts to re-add
repositories which have already been added to the host machine.
