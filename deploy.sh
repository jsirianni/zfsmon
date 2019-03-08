#!/bin/sh

if [ -z "$1" ]
    then
        echo "you need to pass a channel and hook url"
        echo "example: ./deploy.sh alerts https://slack.com/myhookurl"
        exit 1
fi

if [ -z "$2" ]
    then
        echo "you need to pass a channel and hook url"
        echo "example: ./deploy.sh alerts https://slack.com/myhookurl"
        exit 1
fi

wget https://github.com/jsirianni/zfsmon/releases/download/0.1.0/zfsmon-v0.1.0-linux-amd64.zip && \
    unzip zfsmon-v0.1.0-linux-amd64.zip && \
    chmod +x zfsmon && \
    mv zfsmon /usr/local/bin && \
    rm zfsmon-v0.1.0-linux-amd64.zip && \
    crontab -l | grep 'zfsmon' || (crontab -l 2>/dev/null; echo "*/5 * * * * /usr/local/bin/zfsmon --channel $1 --url $2 >> /dev/null 2>&1") | crontab - && \
    echo "" && crontab -l && echo "done. . ."
