# Sharded Mongo

1. Docker Compose Up
2. ./init-cluster.sh
3. Run the below commands in mongos
```bash
docker exec -it mongos mongosh --port 27020

sh.status()  // Optional: check sharding status

use test
sh.enableSharding("test")  // Enable sharding on the `test` database

sh.shardCollection("test.users", { _id: "hashed" })  // Shard the collection
```

