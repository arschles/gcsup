#!/usr/bin/env bash
#
# Build and push Docker images to Docker Hub and quay.io.
#

cd "$(dirname "$0")" || exit 1

VER=0.0.1
# note that we must check TRAVIS_PULL_REQUEST first because, for pull request builds, TRAVIS_BRANCH will contain the name of the branch that the PR is _targeting_ (not the name of the branch it is trying to merge).
# see https://docs.travis-ci.com/user/environment-variables/#Default-Environment-Variables for an explanation of TRAVIS_PULL_REQUEST and TRAVIS_BRANCH
if [[ "$TRAVIS_PULL_REQUEST" != "false" ]]; then
  VER="pr-$TRAVIS_PULL_REQUEST"
elif  [[ "$TRAVIS_BRANCH" != "master" ]]; then
  VER="branch-$TRAVIS_BRANCH"
fi

docker login -e="$QUAY_EMAIL" -u="$QUAY_USERNAME" -p="$QUAY_PASSWORD" quay.io
export IMAGE_NAME="quay.io/arschles/gcsup:$VERSION"
echo "building and pushing $IMAGE_NAME"
make build docker-build docker-push
