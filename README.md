# zfsmon

`zfsmon` is a utility that sends slack alerts when a zpool becomes degraded.

## usage
Available flags:
`zfsmon --help`
```
zfs monitoring daemon

Usage:
  zfsmon [flags]

Flags:
      --alert-file string   hook url (default "/tmp/zfsmon")
      --channel string      slack channel
  -h, --help                help for zfsmon
      --no-alert            do not send alerts
      --print               print the health report
      --url string          hook url (default "/opt/zfsmon/alerts.dat")

```

Print result but take no action:
```
zfsmon --print --no-alert
```

Print result and send alert if necessary:
```
zfsmon --print --channel alerts --url https://myslackwebhook.com/myhook
```

Run from cron every five minutes:
```
*/5 * * * * sudo /usr/local/bin/zfsmon
```

#### notes
zfsmon should be run with `sudo` or with root. Recent versions of ZFSonLinux
do not require root for basic read only operations, however, zfsmon has not been tested
this way.

## install
An install script is provided. It will place the binary in your path and setup
a five minute cronjob. Pass a channel and hookurl as arguments
```
./deploy.sh <channel> <hookurl>
```


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
