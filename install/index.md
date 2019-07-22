---
title: About Docker CE
description: Lists the installation methods
keywords: docker, installation, install, docker ce, docker ee, docker editions, stable, edge
redirect_from:
- /installation/
- /engine/installation/linux/
- /engine/installation/linux/frugalware/
- /engine/installation/frugalware/
- /engine/installation/linux/other/
- /engine/installation/linux/archlinux/
- /engine/installation/linux/cruxlinux/
- /engine/installation/linux/gentoolinux/
- /engine/installation/linux/docker-ce/
- /engine/installation/linux/docker-ee/
- /engine/installation/
- /en/latest/installation/
---

Docker Community Edition (CE) is ideal for developers and small
teams looking to get started with Docker and experimenting with container-based
apps. Docker CE has three types of update channels, **stable**, **test**, and **nightly**:

* **Stable** gives you latest releases for general availability.
* **Test** gives pre-releases that are ready for testing before general availability.
* **Nightly** gives you latest builds of work in progress for the next major release.


## Releases

For the Docker CE engine, the open
repositories [Docker Engine](https://github.com/docker/engine) and
[Docker Client](https://github.com/docker/cli) apply.

Releases of Docker Engine and Docker Client for general availability
are versioned using dotted triples. The components of this triple
are `YY.mm.<patch>` where the `YY.mm` component is referred to as the
year-month release. The version numbering format is chosen to illustrate
cadence and does not guarantee SemVer, but the desired date for general
availability. The version number may have additional information, such as
beta and release candidate qualifications. Such releases are considered
"pre-releases".

The cadence of the year-month releases is every 6 months starting with
the `18.09` release. The patch releases for a year-month release take
place as needed to address bug fixes during its support cycle.

Docker CE binaries for a release are available on [download.docker.com](https://download.docker.com/)
as packages for the supported operating systems. Docker EE binaries are
available on the [Docker Hub](https://hub.docker.com/) for the supported operating systems. The
release channels are available for each of the year-month releases and
allow users to "pin" on a year-month release of choice. The release
channel also receives patch releases when they become available.

### Nightly builds

Nightly builds are created once per day from the master branch. The version
number for nightly builds take the format:

    0.0.0-YYYYmmddHHMMSS-abcdefabcdef

where the time is the commit time in UTC and the final suffix is the prefix
of the commit hash, for example `0.0.0-20180720214833-f61e0f7`.

These builds allow for testing from the latest code on the master branch. No
qualifications or guarantees are made for the nightly builds.

The release channel for these builds is called `nightly`.

### Pre-releases

In preparation for a new year-month release, a branch is created from
the master branch with format `YY.mm` when the milestones desired by
Docker for the release have achieved feature-complete. Pre-releases
such as betas and release candidates are conducted from their respective release
branches. Patch releases and the corresponding pre-releases are performed
from within the corresponding release branch.

While pre-releases are done to assist in the stabilization process, no
guarantees are provided.

Binaries built for pre-releases are available in the test channel for
the targeted year-month release using the naming format `test-YY.mm`,
for example `test-18.09`.

### General availability

Year-month releases are made from a release branch diverged from the master
branch. The branch is created with format `<year>.<month>`, for example
`18.09`. The year-month name indicates the earliest possible calendar
month to expect the release to be generally available. All further patch
releases are performed from that branch. For example, once `v18.09.0` is
released, all subsequent patch releases are built from the `18.09` branch.

Binaries built from this releases are available in the stable channel
`stable-YY.mm`, for example `stable-18.09`, as well as the corresponding
test channel.

### Relationship between CE and EE code

For a given year-month release, Docker releases both CE and EE
variants concurrently. EE is a superset of the code delivered in
CE. Docker maintains publicly visible repositories for the CE code
as well as private repositories for the EE code. Automation (a bot)
is used to keep the branches between CE and EE in sync so as features
and fixes are merged on the various branches in the CE repositories
(upstream), the corresponding EE repositories and branches are kept
in sync (downstream). While Docker and its partners make every effort
to minimize merge conflicts between CE and EE, occasionally they will
happen, and Docker will work hard to resolve them in a timely fashion.

## Next release

The activity for upcoming year-month releases is tracked in the milestones
of the repository.

## Support

Docker CE releases of a year-month branch are supported with patches
as needed for 7 months after the first year-month general availability
release. Docker EE releases are supported for 24 months after the first
year-month general availability release.

This means bug reports and backports to release branches are assessed
until the end-of-life date.

After the year-month branch has reached end-of-life, the branch may be
deleted from the repository.

### Reporting security issues

The Docker maintainers take security seriously. If you discover a security
issue, please bring it to their attention right away!

Please DO NOT file a public issue; instead send your report privately
to security@docker.com.

Security reports are greatly appreciated, and Docker will publicly thank you
for it. Docker also likes to send gifts — if you're into swag, make sure to
let us know. Docker currently does not offer a paid security bounty program
but are not ruling it out in the future.

### Supported platforms

Docker CE is available on multiple platforms. Use the following tables
to choose the best installation path for you.

#### Desktop

{% assign green-check = '![yes](/install/images/green-check.svg){: style="height: 14px; margin: 0 auto"}' %}

| Platform                                                                    |      x86_64       |
|:----------------------------------------------------------------------------|:-----------------:|
| [Docker Desktop for Mac (macOS)](/docker-for-mac/install/)                        | {{ green-check }} |
| [Docker Desktop for Windows (Microsoft Windows 10)](/docker-for-windows/install/) | {{ green-check }} |

#### Server

{% assign green-check = '![yes](/install/images/green-check.svg){: style="height: 14px; margin: 0 auto"}' %}
{% assign install-prefix-ce = '/install/linux/docker-ce' %}

| Platform                                    | x86_64 / amd64                                         | ARM                                                    | ARM64 / AARCH64                                        | IBM Power (ppc64le)                                    | IBM Z (s390x)                                          |
|:--------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|
| [CentOS]({{ install-prefix-ce }}/centos/) | [{{ green-check }}]({{ install-prefix-ce }}/centos/) |                                                        | [{{ green-check }}]({{ install-prefix-ce }}/centos/) |                                                        |                                                        |
| [Debian]({{ install-prefix-ce }}/debian/) | [{{ green-check }}]({{ install-prefix-ce }}/debian/) | [{{ green-check }}]({{ install-prefix-ce }}/debian/) | [{{ green-check }}]({{ install-prefix-ce }}/debian/) |                                                        |                                                        |
| [Fedora]({{ install-prefix-ce }}/fedora/) | [{{ green-check }}]({{ install-prefix-ce }}/fedora/) |                                                        | [{{ green-check }}]({{ install-prefix-ce }}/fedora/) |                                                        |                                                        |
| [Ubuntu]({{ install-prefix-ce }}/ubuntu/) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu/) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu/) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu/) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu/) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu/) |

### Backporting

Backports to the Docker products are prioritized by the Docker company. A
Docker employee or repository maintainer will endeavour to ensure sensible
bugfixes make it into _active_ releases.

If there are important fixes that ought to be considered for backport to
active release branches, be sure to highlight this in the PR description
or by adding a comment to the PR.

### Upgrade path

Patch releases are always backward compatible with its year-month version.

## Not covered

As a general rule, anything not mentioned in this document may change in any release.

## Exceptions

Exceptions are made in the interest of __security patches__. If a break
in release procedure or product functionality is required, it will
be communicated clearly, and the solution will be considered against
total impact.

## Get started

After setting up Docker, you can learn the basics with
[Getting started with Docker](/get-started/).
