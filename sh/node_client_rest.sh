#!/bin/bash
source /etc/profile
runId=`ps aux | grep rest-server | grep -v grep | grep -v log| awk '{print $2}' `

funRunCmd(){
      nohup ./baccli rest-server --node=127.0.0.1:26657 --chain-id=test --laddr=tcp://0.0.0.0:1317 > ~/baccli_rest.log 2>&1 &
}


if [ "$1"  = "start" ]
then
	funRunCmd
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
    sleep 2
    funRunCmd
    echo "restart"
else
    echo "error command"
fi


echo "操作之后状态.."
runId=`ps aux | grep rest-server | grep -v grep | grep -v log| awk '{print $2}'`
if [ "$runId" = "" ]
then
    echo "no run"
else
    echo "run succ"$runId
fi




