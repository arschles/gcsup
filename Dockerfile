FROM alpine:3.2

RUN apk add -U ca-certificates && rm -rf /var/cache/apk/*
ADD gcsup .
RUN mv gcsup /bin
