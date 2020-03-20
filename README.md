# zfsmon

`zfsmon` is a utility that sends alerts when a zpool becomes degraded.

## usage

Available flags:
`zfsmon --help`
```
zfs monitoring util

Usage:
  zfsmon [flags]

Flags:
      --alert-type string      alert system to use
      --daemon                 enable daemon mode
  -h, --help                   help for zfsmon
      --log-level string       logging level [error, warning, info, trace] (default "error")
      --no-alert               do not send alerts
      --slack-channel string   slack channel
      --slack-url string       hook url
      --state-file string      path for the state file (default "/tmp/zfsmon")
```

### alert-type

Available alert types:
- Terminal: Prints alerts to your terminal. Requires `log-level=info`
- Slack: Sends alerts to a slack channel. Requires `slack-channel` and `slack-url`

### Examples

#### Root's Crontab w/ Slack

Run every five minutes. `$HOOK` is retrieved from `/root/.zfsmon`
```
. /root/.zfsmon
*/5 * * * * /usr/local/bin/zfsmon --slack-channel=joe_testing --state-file /root/alert --slack-url=$HOOK --alert-type slack
```

#### Console with Slack and Info logging

Info logging can be used to get full output
```
root@zfs:/home/teamit# ./zfsmon --slack-channel=joe_testing --state-file alert --slack-url=$HOOK --alert-type slack --log-level=info

ERROR: 2020/03/19 22:17:13 logger.go:51: open alert: no such file or directory
INFO: 2020/03/19 22:17:13 logger.go:41: device '/home/teamit/mirror-1-0' in pool 'fake-mirror' is healthy. Status: ONLINE
INFO: 2020/03/19 22:17:13 logger.go:41: device '/home/teamit/mirror-1-1' in pool 'fake-mirror' is healthy. Status: ONLINE
WARNING: 2020/03/19 22:17:13 logger.go:46: device '/home/teamit/mirror-0-0' in pool 'fake' is not healthy. Status: CANT_OPEN
INFO: 2020/03/19 22:17:13 logger.go:41: host: zfs | zpool fake is not in a healthy state, got status: CANT_OPEN
INFO: 2020/03/19 22:17:13 logger.go:41: device '/home/teamit/mirror-0-1' in pool 'fake' is healthy. Status: ONLINE

```

#### Console with Trace logging

Very verbose output can be enabled with `log-level=trace`
```
> zfsmon --alert-type terminal --log-level trace

TRACE: 2020/03/19 22:09:52 logger.go:36: zfsmon config: {"hostname":"zfs","daemon_mode":false,"state":{"file":"alert"},"pools":[{"name":"fake-mirror","state":7,"devices":[{"name":"mirror","type":"mirror","state":7,"Devices":[{"name":"/home/teamit/mirror-1-0","type":"file","state":7},{"name":"/home/teamit/mirror-1-1","type":"file","state":7}]}]}],"alert_config":{"no_alert":false},"alert_state":{"/home/teamit/mirror-0-0":"CANT_OPEN"}}
TRACE: 2020/03/19 22:09:52 logger.go:36: checking pools
TRACE: 2020/03/19 22:09:52 logger.go:36: checking pool 'fake-mirror'
TRACE: 2020/03/19 22:09:52 logger.go:36: device 'mirror' has type 'mirror'
TRACE: 2020/03/19 22:09:52 logger.go:36: checking device '/home/teamit/mirror-1-0' in pool 'fake-mirror'
INFO: 2020/03/19 22:09:52 logger.go:41: device '/home/teamit/mirror-1-0' in pool 'fake-mirror' is healthy. Status: ONLINE
TRACE: 2020/03/19 22:09:52 logger.go:36: checking device '/home/teamit/mirror-1-1' in pool 'fake-mirror'
INFO: 2020/03/19 22:09:52 logger.go:41: device '/home/teamit/mirror-1-1' in pool 'fake-mirror' is healthy. Status: ONLINE
```

## Install

Download the latest release and place it in your path.

It is recommended to run `zfsmon` with `root` or `sudo`. Recent release of OpenZFS
does not require root privileges, however, this code as only been tested with root.
I hope to change this going forward. If you have success using a non root user, please
file an issue with your results.

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
