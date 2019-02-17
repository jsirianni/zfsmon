# zfsmon
A WORK IN PROGRESS

## developing
***docker***
zfsmon requires several dependencies that are only available on Linux.
A docker image is provided to allow for cross platform development. Run the wrapper
script to build and export the binary to your working directory:
```
sudo ./build_docker.sh
```

`build_docker.sh` will not cleanup zfsmon images. To cleanup leftover images:
```
sudo docker images | grep zfsmon | awk '{print $3}' | xargs -n1 sudo docker rmi
```

***build manually***
zfsmon is developed on Ubuntu 18.04 LTS. You should install:
```
golang
git
zfsutils-linux
libzfslinux-dev
```

Retrieve the go dependencies:
```
go get github.com/spf13/cobra
go get github.com/jsirianni/go-libzfs
```

build the binary
```
go build
```
