# Set to the version for this archive
ARG VER=vXX

# This image comes from the onbuild.Dockerfile file in the docs-builder branch
# https://github.com/docker/docker.github.io/blob/docs-builder/onbuild.Dockerfile
FROM docs/docker.github.io:docs-builder-onbuild AS builder

# Reset the docs:onbuild image, which is based on nginx:alpine
# This image comes from the Dockerfile in the nginx-onbuild branch
# https://github.com/docker/docker.github.io/blob/nginx-onbuild/Dockerfile
FROM docs/docker.github.io:nginx-onbuild