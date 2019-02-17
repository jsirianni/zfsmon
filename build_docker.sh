#!/bin/sh

UNIXTIME=`date +'%s'`

sudo docker build . -t zfsmon:$UNIXTIME
sudo docker create --name $UNIXTIME zfsmon:$UNIXTIME >> /dev/null
sudo docker cp $UNIXTIME:/src/zfsmon/zfsmon zfsmon
sudo docker rm $UNIXTIME >> /dev/null
