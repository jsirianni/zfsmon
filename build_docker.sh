#!/bin/sh

UNIXTIME=`date +'%s'`
VERSION=`cat main.go | grep "const VERSION" | cut -c 17- | tr -d '"'`

sudo docker build . -t zfsmon:$VERSION --build-arg version=${VERSION} && \
	sudo docker create --name $UNIXTIME zfsmon:$VERSION >> /dev/null && \
	sudo docker cp $UNIXTIME:/src/zfsmon/zfsmon zfsmon && \
	sudo docker cp $UNIXTIME:/src/zfsmon/zfsmon-v$VERSION-linux-amd64.zip zfsmon-v$VERSION-linux-amd64.zip && \
	sudo docker cp $UNIXTIME:/src/zfsmon/zfsmon-v$VERSION.SHA256SUMS zfsmon-v$VERSION.SHA256SUMS 

sudo docker rm $UNIXTIME >> /dev/null
