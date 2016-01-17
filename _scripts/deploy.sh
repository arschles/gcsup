#!/usr/bin/env bash
#
# Build and push Docker images to Docker Hub and quay.io.
#
# This script was adapted from https://github.com/deis/builder

make -C .. build docker-build docker-push
