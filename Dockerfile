FROM golang:1.19.3-alpine3.16 as builder
LABEL maintainer="dru90i"

COPY . /go/src/gitlab-runner-exporter/
RUN apk --update add --virtual build-deps go git \ 
 && cd /go/src/gitlab-runner-exporter \
 && GOPATH=/go go get \
 && GOPATH=/go go build -o /bin/gitlab_runner_exporter \
 && apk del --purge build-deps \
 && rm -rf /go/bin /go/pkg /var/cache/apk/*

FROM alpine:latest

EXPOSE 9191
RUN apk update
RUN addgroup exporter \
 && adduser -S -G exporter exporter \
 && apk --no-cache add ca-certificates

COPY --from=builder /bin/gitlab_runner_exporter /bin/gitlab_runner_exporter

USER exporter

ENTRYPOINT [ "/bin/gitlab_runner_exporter" ]
