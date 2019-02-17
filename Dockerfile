FROM ubuntu:bionic

# Get build tools and dependency libs
RUN apt-get update
RUN apt-get install -y \
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
RUN go build

