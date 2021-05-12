#!/bin/bash

#clean environment
for i in `seq 1 6`
do
    rm -rf /tmp/$i/objects/*
    rm -rf /tmp/$i/temp/*
done

# prepare the distributed envrionment
for i in `seq 1 6`
do
    mkdir -p /tmp/$i/objects
    mkdir -p /tmp/$i/temp
    mkdir -p /tmp/$i/garbage
done

sudo ifconfig eno1:1 10.29.1.1/16
sudo ifconfig eno1:2 10.29.1.2/16
sudo ifconfig eno1:3 10.29.1.3/16
sudo ifconfig eno1:4 10.29.1.4/16
sudo ifconfig eno1:5 10.29.1.5/16
sudo ifconfig eno1:6 10.29.1.6/16
sudo ifconfig eno1:7 10.29.2.1/16
sudo ifconfig eno1:8 10.29.2.2/16

# rabbitmq env
# wget localhost:15672/cli/rabbitmqadmin #rabbitmq访问
python3 rabbitmqadmin declare exchange name=apiServers type=fanout
python3 rabbitmqadmin declare exchange name=dataServers type=fanout
# sudo rabbitmqctl add_user test test #首次运行需要创建用户和密码
# sudo rabbitmqctl set_permissions -p / test ".*" ".*" ".*" #修改访问权限

# start test env
export RABBITMQ_SERVER=amqp://test:test@localhost:5672
export ES_SERVER=localhost:9200

# start elasticsearch 
# 注意：部署的时候一台机器上开一次即可，多开会出现bug
sudo /usr/share/elasticsearch/bin/elasticsearch > /dev/null &

LISTEN_ADDRESS=10.29.1.1:12346 STORAGE_ROOT=/tmp/1 go run ../dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.2:12346 STORAGE_ROOT=/tmp/2 go run ../dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.3:12346 STORAGE_ROOT=/tmp/3 go run ../dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.4:12346 STORAGE_ROOT=/tmp/4 go run ../dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.5:12346 STORAGE_ROOT=/tmp/5 go run ../dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.6:12346 STORAGE_ROOT=/tmp/6 go run ../dataServer/dataServer.go &

LISTEN_ADDRESS=10.29.2.1:12346 go run ../apiServer/apiServer.go &
LISTEN_ADDRESS=10.29.2.2:12346 go run ../apiServer/apiServer.go &