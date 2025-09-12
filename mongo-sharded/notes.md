# Notes

```bash
docker exec -it mongos mongosh --port 27020

sh.status()  // Optional: check sharding status

use test
sh.enableSharding("test")  // Enable sharding on the `test` database

sh.shardCollection("test.users", { _id: "hashed" })  // Shard the collection
```
