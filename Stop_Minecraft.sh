#! /bin/sh
PORT=25565
PID=`netstat -apn | grep $PORT | awk '{print $7}' | cut -d/ -f 1`
kill -9 $PID
