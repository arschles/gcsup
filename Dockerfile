# commented values are for alpine. https://github.com/arschles/gcsup/issues/2
# FROM alpine:3.2
FROM ubuntu-debootstrap:14.04

# RUN apk add -U libc-dev && rm -rf /var/cache/apk/*
RUN apt-get update -y && apt-get install -y ca-certificates
ADD gcsup .
RUN mv gcsup /bin
