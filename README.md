# zfsmon

`zfsmon` is a utility that sends alerts when a zpool becomes degraded.

## usage

***TODO***: The usage section is out of date.

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

## Install

Download the latest release and place it in your path

## Build

The Makefile will use Docker to build the binary, this is to allow developing
zfsmon on platforms that do not support OpenZFS. The build binary can be found
in `artifacts/`
```
make
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
go build -o zfsmon
```
