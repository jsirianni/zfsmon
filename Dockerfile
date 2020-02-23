# we use the go image for sourcing golang 1.13
# and the ubuntu image becuase the native repos
# contain the ZFS libraries required for building
# zfsmon on Linux
FROM golang:1.13 as go
FROM ubuntu:bionic

ARG version

COPY --from=go /usr/local/go /usr/local/go
RUN ln -s /usr/local/go/bin/go /usr/bin/go

RUN apt-get update
RUN apt-get install -y \
    zip \
	git \
    gcc \
	zfsutils-linux \
	libzfslinux-dev

WORKDIR /zfsmon
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o artifacts/zfsmon
WORKDIR /zfsmon/artifacts
RUN zip zfsmon-v${version}-linux-amd64.zip zfsmon
RUN sha256sum zfsmon-v${version}-linux-amd64.zip >> zfsmon-v${version}.SHA256SUMS
