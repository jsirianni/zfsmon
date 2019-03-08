FROM ubuntu:bionic

ARG version

# Get build tools and dependency libs
RUN apt-get update
RUN apt-get install -y \
        zip \
	golang \
	git \
	zfsutils-linux \
	libzfslinux-dev

# build the binary
WORKDIR /src/zfsmon
ENV GOPATH=/                                                                                   
COPY . ./
RUN go get github.com/spf13/cobra
RUN go get github.com/jsirianni/go-libzfs
RUN GOOS=linux GOARCH=amd64 go build -o zfsmon
RUN zip zfsmon-v${version}-linux-amd64.zip zfsmon
RUN sha256sum zfsmon-v${version}-linux-amd64.zip >> zfsmon-v${version}.SHA256SUMS
