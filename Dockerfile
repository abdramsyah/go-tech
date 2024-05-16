FROM golang:alpine AS builder

ARG SSH_PRIVATE_KEY
RUN mkdir -p /go/src
ADD . /go/src

#Build Source
WORKDIR /go/src

RUN apk add --no-cache --update; \
    apk add git openssh; \
    apk add tzdata; \
    mkdir -p /root/.ssh; \
    chmod 600 /root/.ssh; \
    echo "${SSH_PRIVATE_KEY}" | tr ',' '\n' > /root/.ssh/id_rsa; \
    chmod 600 /root/.ssh/id_rsa
RUN cat /root/.ssh/id_rsa


#Final Build Image
FROM alpine:latest
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/src/main-app /app/main-app

WORKDIR /app

RUN mkdir params; \
    mkdir -p file/temp; \
    mkdir -p file/rbac; \
    mkdir -p file/storage

COPY ./migrations/sql/ migrations/sql

ENTRYPOINT ["/app/main-app"]
