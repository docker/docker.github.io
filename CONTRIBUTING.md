## Contributing

We value your documentation contributions, and we want to make it as easy
as possible to work in this repository. One of the first things to decide is
which branch to base your work on. If you get confused, just ask and we will
help. If a reviewer realizes you have based your work on the wrong branch, we'll
let you know so that you can rebase it.

>**Note**: To contribute code to Docker projects, see the
[Contribution guidelines](opensource/).

### Files not edited here

Files and directories listed in the `path:` keys in
[`.NOT_EDITED_HERE.yaml`](/_data/not_edited_here.yaml) are maintained in other
repositories and should not be edited in this one. Pull requests against these
files will be rejected. Make your edits to the files in the repository and path
in the `source:` key in the YAML file.

### Quickstart

If you spot a problem while reading the documentation and want to try to fix it
yourself, click the **Edit this page** link at the bottom of that page. The
page will open in the Github editor, which means you don't need to know a lot
about Git, or even about Markdown.

When you save, you will be prompted to create a fork if you don't already have
one, and to create a branch in your fork and submit the pull request. We hope
you give it a try!

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
work will be visible on [https://docs.docker.com/](/).

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

- **[vnext-opensource](https://github.com/docker/docker.github.io/tree/vnext-opensource):**
  docs for upcoming features in the [docker/opensource](https://github.com/docker/opensource/)
  project

- **[vnext-swarm](https://github.com/docker/docker.github.io/tree/vnext-swarm):**
  docs for upcoming features in the [docker/swarm](https://github.com/docker/swarm/)
  project

- **[vnext-toolbox](https://github.com/docker/docker.github.io/tree/vnext-toolbox):**
  docs for upcoming features in the [docker/toolbox](https://github.com/docker/toolbox/)
  project

- **[vnext-kitematic](https://github.com/docker/docker.github.io/tree/vnext-kitematic):**
  docs for upcoming features in the [docker/kitematic](https://github.com/docker/kitematic/)
  project

## Collaborate on a pull request

Unless the PR author specifically disables it, you can push commits into another
contributor's PR. You can do it from the command line by adding and fetching
their remote, checking out their branch, and adding commits to it. Even easier,
you can add commits from the Github web UI, by clicking the pencil icon for a
given file in the **Files** view.

If a PR consists of multiple small addendum commits on top of a more significant
one, the commit will usually be "squash-merged", so that only one commit is
merged in. On occasion this is not appropriate and all commits will be kept
separate when merging.

## Pull request guidelines

Help us review your PRs more quickly by following these guidelines.

- Try not to touch a large number of files in a single PR if possible.

- Don't change whitespace or line wrapping in parts of a file you are not
  editing for other reasons. Make sure your text editor is not configured to
  automatically reformat the whole file when saving.

- A Netlify test runs for each PR that is against one of our long-lived
  branches like `master` and the `vnext` branches, and deploys the result of
  your PR to a staging site. The URL will be available at the bottom of the PR
  in the **Conversation** view. Check the staging site for problems and fix them
  if necessary. Reviewers will check the staging site too.

If you can think of other ways we could streamline the review process, let us
know.

## Style guide

If you have questions about how to write for Docker's documentation, have a look
at the [style guide](/opensource/doc-style/). The style
guide provides guidance about grammar, syntax, formatting, styling, language, or
tone. If something isn't clear in the guide, submit an issue to let us
know or submit a pull request to help us improve it.
hai
a
