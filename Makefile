SHORT_NAME := "gcsup"
# dockerized development environment variables
DEV_ENV_IMAGE := quay.io/deis/go-dev:latest
DEV_ENV_WORK_DIR := /go/src/github.com/arschles/gcsup
DEV_ENV_PREFIX := docker run --rm -e GO15VENDOREXPERIMENT=1 -e GOPATH=/go -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}

IMAGE_NAME ?= quay.io/arschles/gcsup:latest

bootstrap:
	${DEV_ENV_CMD} glide up

build:
	${DEV_ENV_CMD} go build -o gcsup

test:
	${DEV_ENV_CMD} go test ${glide nv}

docker-build:
	docker build -t ${IMAGE_NAME} .

docker-push:
	docker push ${IMAGE_NAME}
