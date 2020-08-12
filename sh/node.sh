#!/bin/bash
source /etc/profile
runId=`ps aux | grep bacd | grep -v grep | grep -v log| awk '{print $2}' `
if [ "$1"  = "start" ]
then
    nohup ./bacd start > ~/bacd.log 2>&1 &
    echo "start..."
elif  [ "$1"  = "stop" ]
then
    if [ "$runId" = "" ]
    then
        echo "need not stop ,not run"
    else
        echo "stop..."$runId
        kill -15 $runId
    fi
elif  [ "$1"  = "restart" ]
then
    kill -15 $runId
    nohup ./bacd start > ~/bacd.log 2>&1 &
    echo "restart"
else
    echo "error command"
fi


echo "操作之后状态.."
runId=`ps aux | grep bacd | grep -v grep | grep -v log| awk '{print $2}'`
if [ "$runId" = "" ]
then
    echo "no run"
else
    echo "run succ"$runId
fi