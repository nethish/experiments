# Citus postgres
* Sharded postgres extension. Has coordinator and workers
* Workers query for the data, and coordinator combines and aggregates data if needed.
* By default, Citus creates 32 shards when you distribute a table using create_distributed_table() without specifying a shard count.

## Setup
Once the 3 containers come up, 
```bash
docker exec -it citus-worker1 psql -U postgres -d citus -c "CREATE EXTENSION citus;"
docker exec -it citus-coordinator psql -U postgres -d citus -c "CREATE EXTENSION citus;"
docker exec -it citus-worker2 psql -U postgres -d citus -c "CREATE EXTENSION citus;"


docker exec -it citus-coordinator psql -U postgres -d citus -c "SELECT * from master_add_node('citus-worker1', 5432);"
docker exec -it citus-coordinator psql -U postgres -d citus -c "SELECT * from master_add_node('citus-worker2', 5432);"

docker exec -it citus-coordinator psql -U postgres -d citus -c "SELECT * FROM master_get_active_worker_nodes();"

docker exec -it citus-coordinator psql -U postgres -d citus -c "
CREATE TABLE events (
    id bigserial,
    user_id int,
    created_at timestamp default now()
);
"

docker exec -it citus-coordinator psql -U postgres -d citus -c "SELECT create_distributed_table('events', 'user_id');"

docker exec -it citus-coordinator psql -U postgres -d citus -c "
INSERT INTO events (user_id) 
SELECT generate_series(1, 1000);
"

docker exec -it citus-coordinator psql -U postgres -d citus -c "SELECT count(*) FROM events;"
"

```

To query the shard details
```SQL
SELECT 
    shardid,
    table_name,
    shard_name 
FROM pg_dist_placement p
JOIN pg_dist_shard s ON p.shardid = s.shardid
JOIN pg_class c ON s.logicalrelid = c.oid
WHERE nodename = 'citus-worker1'
ORDER BY shardid;
```

## Useful commands
```SQL
--- Peek into num shards
SELECT * FROM pg_dist_shard WHERE logicalrelid = 'events'::regclass;
```
