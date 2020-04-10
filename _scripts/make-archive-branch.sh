#!/bin/bash

while getopts ":hv:s:" opt; do
  case ${opt} in
    v ) version="$OPTARG"
      ;;
    s ) SHA="$OPTARG"
      ;;
    \? ) echo "Usage: $0 [-h] | -v <docker-version> -s <GIT-SHA-OF-ARCHIVE>"
         echo ""
         echo "<docker-version> is in the format \"17.09\" without any \"v\" prepended."
         echo "<GIT-SHA-OF-ARCHIVE> is a Git SHA which is the last commit before work started that doesn't apply to this archive."
         break
      ;;
  esac
done

# If we can't find a version, exit gracefully
if [ -z "$version" ]; then
  echo "-v is a required argument and was not detected."
  exit 1
else
  # Do some very basic and naive checks on the format
  # We expect it to start with a number
  if ! [[ $version =~ ^[1-9].*[0-9]$ ]]; then
    echo "Invalid version. Expected numbers and dots, but got $version"
    exit 1
  fi
fi

# If we don't have a SHA for the archive, exit
if [ -z "$SHA" ]; then
  echo "-s is a required argument and was not detected."
  exit 1
fi

# Exit if we are not running from a clean master

BRANCH=$(git branch | grep '*' | awk {'print $2'})

if [ "$BRANCH" != "master" ]; then
  echo "You are on branch $BRANCH but an archive can only be created from master. Exiting."
  exit 1
fi

## Make sure our working branch is clean

BRANCH_STATUS=$(git diff --cached --exit-code)
BRANCH_STATUS_UNTRACKED=$(git ls-files --other --directory --exclude-standard | head -n 1)

if [ $BRANCH_STATUS -o "$BRANCH_STATUS_UNTRACKED" ]; then
  echo "Branch has uncommitted changes or untracked files. Exiting to protect you."
  echo "Use a fresh clone for this, not your normal working clone."
  echo "If you don't want to set up all your remotes again,"
  echo "recursively copy (cp -r on a mac) your docker.github.io directory to"
  echo "a new directory like archive-clone and run this script again from there."
  exit 1
fi

# Check out the SHA

echo "Making archive based on $SHA"
git pull -q
git checkout -q ${SHA}
status=$?

if [ $status -ne 0 ]; then
  echo "Failed to check out $SHA. Exiting."
  exit 1
fi


# Create the archive branch

git checkout -b v"$version" && echo "Created branch v$version and checked it out."

# Replace the Dockerfile, set the version
cat Dockerfile.archive |sed "s/vXX/v${version}/g" > Dockerfile

# Run _scripts/fetch_upstream_resources.sh once in local mode
bash _scripts/fetch-upstream-resources.sh -l

# Add a redirect page for each section that doesn't apply to the archives, where
# the reader should look in the live content instead
# Currently, this is:
# /samples/
# /docker-id/
# /docker-hub/
# These rely on _layout/archive-redirect.html

only_live_contents=("samples" "docker-id" "docker-hub")

for dir in "${only_live_contents[@]}"; do
  echo "Replacing contents of $dir with a redirect stub"
  # Figure out the title, which should be title case with spaces instead of dashes
  dir_title=$(echo $dir | sed 's/-/\ /g' | awk '{for(i=1;i<=NF;i++){ $i=toupper(substr($i,1,1)) substr($i,2) }}1')
  echo "dir_title is ${dir_title}"
  rm -Rf $dir/*
  cat << EOF > "$dir/index.html"
---
layout: archive-redirect
prod_title: "$dir_title"
prod_url: "$dir"
---
EOF

done

echo "You are almost done. There are FOUR manual steps left to do."
echo "1.  Edit _data/toc.yaml and remove all the entries under the following"
echo "    locations except for their roots. For instance, remove all of"
echo "    /samples/ entries except /samples/ itself."
echo "    A valid example for samples would look like:"

cat << EOF
          samples:
          - sectiontitle: Sample applications
            section:
            - path: /samples/
              title: Samples home
EOF

echo "    Do this for the following sections:"
for dir in "${only_live_contents[@]}"; do
  echo "      /$dir/"
done
echo

echo "2.  Do a local build and run to test your archive."
echo "        docker build -t archivetest ."
echo "        docker run --rm -it -p 4000:4000 archivetest"

echo "3.  After carefully reviewing the output of 'git status' to make sure no"
echo "    files are being added by mistake, Run the following to add and commit"
echo "    all your changes, including new files:"
echo "        git commit -m \"Created archive branch v$version\" -A"

echo "4.  Push your archive to the upstream remote:"
echo "        git push upstream v$version"
echo
echo "99. If you want to bail out of this operation completely,"
echo "    and get back to master, run the following:"
echo
echo "      git reset --hard; git clean -fd; git checkout master; git branch -D v$version"
echo
