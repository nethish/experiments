# Cassandra
* Datacenter, Rack, and Num Tokens, Virtuals Tokens -- How does it work together?



## CQL
```SQL
docker exec -it cassandra1 cqlsh
DESCRIBE KEYSPACES;
SELECT cluster_name, listen_address FROM system.local;
SELECT peer, data_center, rack FROM system.peers;


CREATE KEYSPACE keyspace1
WITH REPLICATION = {
  'class': 'NetworkTopologyStrategy',
  'datacenter1': 3
};

use keyspace1;


CREATE TABLE users (
  id UUID PRIMARY KEY,
  name TEXT,
  email TEXT,
  age INT
);

-- Insert data
INSERT INTO users (id, name, email, age) 
VALUES (uuid(), 'John Doe', 'john@example.com', 30);

-- Query data
SELECT * FROM users;
```

```bash
docker exec cassandra1 nodetool status
docker exec cassandra1 nodetool ring
docker exec cassandra1 nodetool info

# See data distribution
docker exec cassandra1 nodetool cfstats

# Check node health
docker exec cassandra1 nodetool tpstats

# Repair data
docker exec cassandra1 nodetool repair

# Flush memtables to disk
docker exec cassandra1 nodetool flush
```




```bash
# Install via pip
pip install cqlsh

# Connect to any node
cqlsh localhost 9042  # cassandra1
cqlsh localhost 9043  # cassandra2
cqlsh localhost 9044  # cassandra3
cqlsh localhost 9045  # cassandra4
```
