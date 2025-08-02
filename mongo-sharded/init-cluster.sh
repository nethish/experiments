#!/bin/bash

sleep 10

docker exec -it configsvr mongosh --port 27019 --eval 'rs.initiate({_id: "configReplSet", configsvr: true, members: [{_id: 0, host: "configsvr:27019"}]})'
docker exec -it shard1 mongosh --port 27018 --eval 'rs.initiate({_id: "shard1ReplSet", members: [{_id: 0, host: "shard1:27018"}]})'
docker exec -it shard2 mongosh --port 27017 --eval 'rs.initiate({_id: "shard2ReplSet", members: [{_id: 0, host: "shard2:27017"}]})'

sleep 10

docker exec -it mongos mongosh --port 27020 --eval 'sh.addShard("shard1ReplSet/shard1:27018")'
docker exec -it mongos mongosh --port 27020 --eval 'sh.addShard("shard2ReplSet/shard2:27017")'
