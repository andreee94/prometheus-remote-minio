FROM golang:1.11 AS builder
WORKDIR /go/src/app
COPY . .
RUN go build -o /usr/bin/prometheus-remote-minio .

###############################################

FROM ubuntu:16.04
COPY --from=builder /usr/bin/prometheus-remote-minio /usr/bin/prometheus-remote-minio
RUN apt-get update && apt-get install -y \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["/usr/bin/prometheus-remote-minio"]
