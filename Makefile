SHORT_NAME := "gcsup"
# dockerized development environment variables
REPO_PATH := github.com/arschles/${SHORT_NAME}
DEV_ENV_IMAGE := quay.io/deis/go-dev:0.2.0
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}
DEV_ENV_PREFIX := docker run --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}

IMAGE_NAME := quay.io/arschles/gcsup:latest

build:
	${DEV_ENV_PREFIX} go build -o gcsup

docker-build:
	docker build -t ${IMAGE_NAME} .

docker-push:
	docker push ${IMAGE_NAME}
