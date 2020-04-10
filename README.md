# Docs @ Docker
Welcome to the repo for our documentation. This is the source for	
[https://docs.docker.com/](https://docs.docker.com/).	

Feel free to send us pull requests and file issues. Our docs are completely	
open source and we deeply appreciate contributions from our community!
## Table of Contents

- [Providing feedback](#providing-feedback)
- [Contributing](#contributing)
  - [Files not edited here](#files-not-edited-here)
  - [Overall doc improvements](#overall-doc-improvements)
  - [Specific new features for a project](#specific-new-features-for-a-project)
- [Per-PR staging on GitHub](#per-pr-staging-on-github)
- [Staging the docs](#staging-the-docs)
- [Read these docs offline](#read-these-docs-offline)
- [Important files](#important-files)
- [Relative linking for GitHub viewing](#relative-linking-for-github-viewing)
  - [Testing changes and practical guidance](#testing-changes-and-practical-guidance)
  - [Per-page front-matter](#per-page-front-matter)
  - [Creating tabs](#creating-tabs)
  - [Running in-page Javascript](#running-in-page-javascript)
  - [Images](#images)
- [Beta content disclaimer](#beta-content-disclaimer)
- [Accessing unsupported archived documentation](#accessing-unsupported-archived-documentation)
- [Building archives and the live published docs](#building-archives-and-the-live-published-docs)
- [Creating a new archive](#creating-a-new-archive)
- [Copyright and license](#copyright-and-license)


## Providing feedback

We really want your feedback, and we've made it easy. You can edit a page or
request changes in the right column of every page on
[docs.docker.com](https://docs.docker.com/). You can also rate each page by
clicking a link at the footer.

**Only file issues about the documentation in this repository.** One way
to think about this is that you should file a bug here if your issue is that you
don't see something that should be in the docs, or you see something incorrect
or confusing in the docs.

- If your problem is a general question about how to configure or use Docker,
  ask in [https://forums.docker.com](https://forums.docker.com) instead.

- If you have an idea for a new feature or behavior change in a specific aspect
  of Docker, or have found a bug in part of Docker, file that issue in
  the project's code repository.

## Contributing

We value your documentation contributions, and we want to make it as easy
as possible to work in this repository. One of the first things to decide is
which branch to base your work on. If you get confused, just ask and we will
help. If a reviewer realizes you have based your work on the wrong branch, we'll
let you know so that you can rebase it.

>**Note**: To contribute code to Docker projects, see the
[Contribution guidelines](CONTRIBUTING.md).

### Files not edited here

Files and directories listed in the `path:` keys in
[`.NOT_EDITED_HERE.yaml`](.NOT_EDITED_HERE.yaml) are maintained in other
repositories and should not be edited in this one. Pull requests against these
files will be rejected. Make your edits to the files in the repository and path
in the `source:` key in the YAML file.

### Overall doc improvements

Most commits will be made against the `master` branch. This include:

- Conceptual and task-based information not specific to new features
- Restructuring / rewriting
- Doc bug fixing
- Typos and grammar errors

One quirk of this project is that the `master` branch is where the live docs are
published from, so upcoming features can't be documented there. See
[Specific new features for a project](#specific-new-features-for-a-project)
for how to document upcoming features. These feature branches will be periodically
merged with `master`, so don't worry about fixing typos and documentation bugs
there.

>Do you enjoy creating graphics? Good graphics are key to great documentation,
and we especially value contributions in this area.

### Specific new features for a project

Our docs cover many projects which release at different times. **If, and only if,
your pull request relates to a currently unreleased feature of a project, base
your work on that project's `vnext` branch.** These branches were created by
cloning `master` and then importing a project's `master` branch's docs into it
(at the time of the migration), in a way that preserved the commit history. When
a project has a release, its `vnext` branch will be merged into `master` and your
work will be visible on [https://docs.docker.com/](https://docs.docker.com/).

The following `vnext` branches currently exist:

- **[vnext-engine](https://github.com/docker/docker.github.io/tree/vnext-engine):**
  docs for upcoming features in the [docker/docker](https://github.com/moby/moby/)
  project

- **[vnext-compose](https://github.com/docker/docker.github.io/tree/vnext-compose):**
  docs for upcoming features in the [docker/compose](https://github.com/docker/compose/)
  project

- **[vnext-distribution](https://github.com/docker/docker.github.io/tree/vnext-distribution):**
  docs for upcoming features in the [docker/distribution](https://github.com/docker/distribution/)
  project

- **[vnext-swarm](https://github.com/docker/docker.github.io/tree/vnext-swarm):**
  docs for upcoming features in the [docker/swarm](https://github.com/docker/swarm/)
  project

- **[vnext-toolbox](https://github.com/docker/docker.github.io/tree/vnext-toolbox):**
  docs for upcoming features in the [docker/toolbox](https://github.com/docker/toolbox/)
  project

## Per-PR staging on GitHub

For every PR against `master` and all the long-lived branches, a staged version
of the site is built using Netlify. If the site builds, you will see
**deploy/netlify — Deploy preview ready**. Otherwise, you will see an error.
Click **Details** to review the staged site or the errors that prevented it from
building. Review the staged site and amend your commit if necessary. Reviewers
will also check the staged site before merging the PR, to protect the integrity
of [https://docs.docker.com/](https://docs.docker.com/).

## Staging the docs

You have three options:

1.  On your local machine, clone this repo and run our staging container:

    ```bash
    git clone --recursive https://github.com/docker/docker.github.io.git
    cd docker.github.io
    docker-compose up
    ```

    If you haven't got Docker Compose installed,
    [follow these installation instructions](https://docs.docker.com/compose/install/).

    The container runs in the background and incrementally rebuilds the site each
    time a file changes. You can keep your browser open to http://localhost:4000/
    and refresh to see your changes. The container runs in the foreground, but
    you can use `CTRL+C` to get the command prompt back. To stop the container,
    issue the following command:

    ```bash
    docker-compose down
    ```

2.  Install Jekyll and GitHub Pages on your local machine.

    a. Clone this repo by running:

       ```bash
       git clone --recursive https://github.com/docker/docker.github.io.git
       ```

    b. Install Ruby 2.3 or later as described in [Installing Ruby](https://www.ruby-lang.org/en/documentation/installation/).

    c. Install Bundler:

       ```bash
       gem install bundler
       ```

    d. If you use Ubuntu, install packages required for the Nokogiri HTML
       parser:

       ```bash
       sudo apt-get install ruby-dev zlib1g-dev liblzma-dev
       ```

    e. Install Jekyll and other required dependencies:

       ```bash
       bundle install
       ```

       >**Note**: You may need to install some packages manually.

    f. Change the directory to `docker.github.io`.

    g. Use the `jekyll serve` command to continuously build the HTML output.

    The `jekyll serve` process runs in the foreground, and starts a web server
    running on http://localhost:4000/ by default. To stop it, use `CTRL+C`.
    You can continue working in a second terminal and Jekyll will rebuild the
    website incrementally. Refresh the browser to preview your changes.

3. Build and run a Docker image for your working branch.

   ```bash
   $ docker build -t docker build -t docs/docker.github.io:<branch_name> .
   $ docker run --rm -it -p 4000:4000 docs/docker.github.io:<branch_name>
    ```

   After the `docker run` command, copy the URL provided in the container build output in a browser,
   http://0.0.0.0:4000, and verify your changes.

## Read these docs offline

To read the docs offline, you can use either a standalone container or a swarm service.
To see all available tags, go to
[Docker Hub](https://docs.docker.com/docker-hub/).

The following examples use the `latest` tag:

- Run a single container:

  ```bash
  docker run  -it -p 4000:4000 docs/docker.github.io:latest
  ```

- Run a swarm service:

  ```bash
  docker service create -p 4000:4000 --name localdocs --replicas 1 docs/docker.github.io:latest
  ```

  This example uses only a single replica, but you could run as many replicas as you'd like.

Either way, you can now access the docs at port 4000 on your Docker host.

## Important files

- `/_data/toc.yaml` defines the left-hand navigation for the docs
- `/js/menu.js` defines most of the docs-specific JS such as TOC generation and menu syncing
- `/css/style.scss` defines the docs-specific style rules
- `/_layouts/docs.html` is the HTML template file, which defines the header and footer, and includes all the JS/CSS that serves the docs content

## Relative linking for GitHub viewing

Feel free to link to `../foo.md` so that the docs are readable in GitHub, but keep in mind that Jekyll templating notation
`{% such as this %}` will render in raw text and not be processed. In general it's best to assume the docs are being read
directly on [https://docs.docker.com/](https://docs.docker.com/).

### Testing changes and practical guidance

If you want to test a style change, or if you want to see how to achieve a
particular outcome with Markdown, Bootstrap, JQuery, or something else, have
a look at `test.md` (which renders in the site at `/test/`).

### Per-page front-matter

The front-matter of a given page is in a section at the top of the Markdown
file that starts and ends with three hyphens. It includes YAML content. The
following keys are supported. The title, description, and keywords are required.

| Key                    | Required  | Description                             |
|------------------------|-----------|-----------------------------------------|
| title                  | yes       | The page title. This is added to the HTML output as a `<h1>` level header. |
| description            | yes       | A sentence that describes the page contents. This is added to the HTML metadata. |
| keywords               | yes       | A comma-separated list of keywords. These are added to the HTML metadata. |
| redirect_from          | no        | A YAML list of pages which should redirect to THIS page. At build time, each page listed here is created as a HTML stub containing a 302 redirect to this page. |
| notoc                  | no        | Either `true` or `false`. If `true`, no in-page TOC is generated for the HTML output of this page. Defaults to `false`. Appropriate for some landing pages that have no in-page headings.|
| toc_min                | no        | Ignored if `notoc` is set to `true`. The minimum heading level included in the in-page TOC. Defaults to `2`, to show `<h2>` headings as the minimum. |
| toc_max                | no        | Ignored if `notoc` is set to `false`. The maximum heading level included in the in-page TOC. Defaults to `3`, to show `<h3>` headings. Set to the same as `toc_min` to only show `toc_min` level of headings. |
| tree                   | no        | Either `true` or `false`. Set to `false` to disable the left-hand site-wide navigation for this page. Appropriate for some pages like the search page or the 404 page. |
| no_ratings             | no        | Either `true` or `false`. Set to `true` to disable the page-ratings applet for this page. Defaults to `false`. |
| skip_read_time             | no        | Set to `true` to disable the 'Estimated reading time' banner for this page. |

The following is an example of valid (but contrived) page metadata. The order of
the metadata elements in the front-matter is not important.

```liquid
---
description: Instructions for installing Docker on Ubuntu
keywords: requirements, apt, installation, ubuntu, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/ubuntulinux/
- /installation/ubuntulinux/
- /engine/installation/linux/ubuntulinux/
title: Get Docker for Ubuntu
toc_min: 1
toc_max: 6
skip_read_time: true
no_ratings: true
---
```

### Creating tabs

The use of tabs, as on pages like [https://docs.docker.com/engine/api/](/engine/api/), requires
the use of HTML. The tabs use Bootstrap CSS/JS, so refer to those docs for more
advanced usage. For a basic horizontal tab set, copy/paste starting from this
code and implement from there. Keep an eye on those `href="#id"` and `id="id"`
references as you rename, add, and remove tabs.

```
<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab1">TAB 1 HEADER</a></li>
  <li><a data-toggle="tab" data-target="#tab2">TAB 2 HEADER</a></li>
</ul>
<div class="tab-content">
  <div id="tab1" class="tab-pane fade in active">TAB 1 CONTENT</div>
  <div id="tab2" class="tab-pane fade">TAB 2 CONTENT</div>
</div>
```

For more info and a few more permutations, see `test.md`.

### Running in-page Javascript

If you need to run custom Javascript within a page, and it depends upon JQuery
or Bootstrap, make sure the `<script>` tags are at the very end of the page,
after all the content. Otherwise the script may try to run before JQuery and
Bootstrap JS are loaded.

> **Note**: In general, this is a bad idea.

### Images

Don't forget to remove images that are no longer used. Keep the images sorted
in the local `images/` directory, with names that naturally group related images
together in alphabetical order. For instance prefer `settings-file-share.png`
and `settings-proxies.png` to `file-share-settings.png` and
`proxies-settings.png`. You may also use numbers, especially in the case of a
sequence, e.g., `run-only-the-images-you-trust-1.svg`
`run-only-the-images-you-trust-2.png` `run-only-the-images-you-trust-3.png`.

When applicable, capture windows rather than rectangular regions. This
eliminates unpleasant background and saves the editors the need to crop.

On Mac, capture windows without shadows. To this end, once you pressed
`Command-Shift-4`, press Option while clicking on the window. To disable
shadows once for all, run:

```bash
$ defaults write com.apple.screencapture disable-shadow -bool TRUE
$ killall SystemUIServer  # restart it.
```

You can restore shadows later with `-bool FALSE`.

In order to keep the Git repository light, _please_ compress the images
(losslessly). On Mac you may use (ImageOptim)[https://imageoptim.com] for
instance. Be sure to compress the images *before* adding them to the
repository, doing it afterwards actually worsens the impact on the Git repo (but
still optimizes the bandwidth during browsing).

## Beta content disclaimer
```bash
> BETA DISCLAIMER
>
> This is beta content. It is not yet complete and should be considered a work in progress. This content is subject to change without notice.
```

## Accessing unsupported archived documentation

Supported documentation includes the current version plus the previous five versions.

If you are using a version of the documentation that is no longer supported, which means that the version number is not listed in the site dropdown list, you can still access that documentation in the following ways:

- By entering your version number and selecting it from the branch selection list for this repo
- By directly accessing the Github URL for your version. For example, https://github.com/docker/docker.github.io/tree/v1.9 for `v1.9`
- By running a container of the specific [tag for your documentation version](https://hub.docker.com/r/docs/docker.github.io/tags)
in Docker Hub. For example, run the following to access `v1.9`:

 ```bash
  docker run -it -p 4000:4000 docs/docker.github.io:v1.9
  ```

## Building archives and the live published docs

All the images described below are automatically built using Docker Hub. To
build the site manually, from scratch, including all utility and archive images,
see the [README in the publish-tools branch](https://github.com/docker/docker.github.io/blob/publish-tools/README.md).

- Some utility images are built from Dockerfiles in the `publish-tools` branch.
  See its [README](https://github.com/docker/docker.github.io/blob/publish-tools/README.md)
  for details.
- Each archive branch automatically builds an image tagged
  `docs/docker.github.io:v<VERSION>` when a change is merged into that branch.
- The `master` branch has a Dockerfile which uses the static HTML from each
  archive image, in combination with the Markdown
  files in `master` and some upstream resources which are fetched at build-time,
  to create the full site at [https://docs.docker.com/](/). All
  of the long-running branches, such as `vnext-engine`, `vnext-compose`, etc,
  use the same logic.

## Creating a new archive

When a new Docker Engine - Community Stable version is released, the previous state of `master`
is archived into a version-specific branch like `v17.09`, by doing the following:

1.  Create branch based off the commit hash before the new version was released.

    ```bash
    $ git checkout <HASH>
    $ git checkout -b v17.09
    ```

2.  Run the `_scripts/fetch-upstream-resources.sh` script. This puts static
    copies of the files in place that the `master` build typically fetches
    each build.

    ```bash
    $ _scripts/fetch-upstream/resources.sh
    ```

3.  Overwrite the `Dockerfile` with the `Dockerfile.archive` (use `cp` rather
    than `mv` so you don't inadvertently remove either file). Edit the resulting
    `Dockerfile` and set the `VER` build argument to the appropriate value, like
    `v17.09`.

    ```bash
    $ mv Dockerfile.archive Dockerfile
    $ vi Dockerfile

      < edit the variable and save >
    ```

4.  Do `git status` and add all changes, being careful not to add anything extra
    by accident. Commit your work.

    ```bash
    $ git status
    $ git add <filename>
    $ git add <filename> (etc etc etc)
    $ git commit -m "Creating archive for 17.09 docs"
    ```

5.  Make sure the archive builds.

    ```bash
    $ docker build -t docker build -t docs/docker.github.io:v17.09 .
    $ docker run --rm -it -p 4000:4000 docs/docker.github.io:v17.09
    ```

    After the `docker run` command, browse to `http://localhost:4000/` and
    verify that the archive is self-browseable.

6.  Push the branch to the upstream repository. Do not create a pull request
    as there is no reference branch to compare against.

    ```bash
    $ git push upstream v17.09
    ```

## Copyright and license

Code and documentation copyright 2017 Docker, inc, released under the Apache 2.0 license.
