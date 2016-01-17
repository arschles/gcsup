#!/usr/bin/env bash
#
# Build and push Docker images to Docker Hub and quay.io.
#
# This script was adapted from https://github.com/deis/builder

docker login -e="$QUAY_EMAIL" -u="$QUAY_USERNAME" -p="$QUAY_PASSWORD" quay.io
make -C .. build docker-build docker-push
