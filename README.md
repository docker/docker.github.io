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
- [Per-PR staging on GitHub](#per-pr-staging-on-github)
- [Build and preview the docs locally](#build-and-preview-the-docs-locally)
  - [Build the docs with deployment features enabled](#build-the-docs-with-deployment-features-enabled)
- [Important files](#important-files)
- [Relative linking for GitHub viewing](#relative-linking-for-github-viewing)
  - [Testing changes and practical guidance](#testing-changes-and-practical-guidance)
  - [Per-page front-matter](#per-page-front-matter)
  - [Creating tabs](#creating-tabs)
  - [Running in-page Javascript](#running-in-page-javascript)
  - [Images](#images)
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

### Files not edited here

Some files and directories are maintained in the upstream repositories. You can find a list of such files in
[`_config_production.yml`](_config_production.yml#L52). Pull requests against these files will be rejected.

### Overall doc improvements

Pull requests should be opened against the `master` branch, this includes:

- Conceptual and task-based information not specific to new features
- Restructuring / rewriting
- Doc bug fixing
- Typos and grammar errors

> Do you enjoy creating graphics? Good graphics are key to great documentation,
> and we especially value contributions in this area.

## Per-PR staging on GitHub

For every PR against `master`, a staged version of the site is built using Netlify.
If the site builds, you will see **deploy/netlify — Deploy preview ready**.
Otherwise, you will see an error. Click **Details** to review the staged site or
the errors that prevented it from building. Review the staged site and amend your
commit if necessary. Reviewers will also check the staged site before merging the
PR, to protect the integrity of [https://docs.docker.com/](https://docs.docker.com/).

## Build and preview the docs locally

On your local machine, clone this repo:

```bash
git clone --recursive https://github.com/docker/docker.github.io.git
cd docker.github.io
```

Then build and run the documentation with [Docker Compose](https://docs.docker.com/compose/)

```bash
docker-compose up -d --build
```

> Docker Compose is included with [Docker Desktop](https://docs.docker.com/desktop/).
> If you don't have Docker Compose installed, [follow these installation instructions](https://docs.docker.com/compose/install/).

Once the container is built and running, visit [http://localhost:4000](http://localhost:4000)
in your web browser to view the docs.

To rebuild the docs after you made changes, run the `docker-compose up` command
again. This rebuilds the documentation, and updates the container with your changes:

```bash
docker-compose up -d --build
```

Once the container is built and running, visit [http://localhost:4000](http://localhost:4000)
in your web browser to view the docs.


To stop the staging container, use the `docker-compose down` command:

```bash
docker-compose down
```

### Build the docs with deployment features enabled

The default configuration for local builds of the documentation disables some
features to allow for a shorter build-time. The following options differ between
local builds, and builds that are deployed to docs.docker.com:

- search auto-completion, and generation of `js/metadata.json`
- google analytics
- page ratings
- `sitemap.xml` generation
- minification of stylesheets (css/style.css)
- adjusting "edit this page" links for content in other repositories

If you want to contribute in these areas, you can perform a "production" build
locally.

To preview the documentation with deployment features enabled, you need to set the
`JEKYLL_ENV` environment variable when building the documentation;

```bash
JEKYLL_ENV=production docker-compose up --build
```

Once the container is built and running, visit [http://localhost:4000](http://localhost:4000)
in your web browser to view the docs.

To rebuild the docs after you make changes, repeat the steps above.

<!--
TODO re-enable auto-builds, or push master builds to Docker hub
## Read these docs offline

To read the docs offline, you can use either a standalone container or a swarm service.
To see all available tags, go to [Docker Hub](https://docs.docker.com/docker-hub/).

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
-->

## Important files

- `/_data/toc.yaml` defines the left-hand navigation for the docs
- `/js/docs.js` defines most of the docs-specific JS such as TOC generation and menu syncing
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
| no_ratings             | no        | Either `true` or `false`. Set to `true` to disable the page-ratings applet for this page. Defaults to `false`. |
| skip_read_time         | no        | Set to `true` to disable the 'Estimated reading time' banner for this page. |
| sitemap                | no        | Exclude the page from indexing by search engines. When set to `false`, the page is excluded from `sitemap.xml`, and a `<meta name="robots" content="noindex"/>` header is added to the page. |

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

```html
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
(losslessly). On Mac you may use [ImageOptim](https://imageoptim.com) for
instance. Be sure to compress the images *before* adding them to the
repository, doing it afterwards actually worsens the impact on the Git repo (but
still optimizes the bandwidth during browsing).

## Copyright and license

Copyright 2013-2021 Docker, inc, released under the Apache 2.0 license.
