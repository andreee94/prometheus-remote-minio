FROM golang:latest AS builder
WORKDIR /go/src/app
COPY ./src .
RUN go build -o /usr/bin/prometheus-remote-minio .

###############################################

# FROM ubuntu:16.04
# COPY --from=builder /usr/bin/prometheus-remote-minio /usr/bin/prometheus-remote-minio
# RUN apt-get update && apt-get install -y \
#     ca-certificates \
#  && rm -rf /var/lib/apt/lists/*

# ENTRYPOINT ["/usr/bin/prometheus-remote-minio"]

FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN apk --no-cache add libc6-compat

RUN mkdir /tmp/buffer

COPY --from=builder /usr/bin/prometheus-remote-minio /usr/bin/prometheus-remote-minio


# RUN apt-get update && apt-get install -y \
#     ca-certificates \
#  && rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["/usr/bin/prometheus-remote-minio"]
# ENTRYPOINT ["/bin/sh"]
