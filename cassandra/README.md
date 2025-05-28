# Cassandra
* Datacenter, Rack, and Num Tokens, Virtuals Tokens -- How does it work together?
* Masterless - Every node is a peer.
* Any node can act as a replica, coordinator and gossip participant
* SEED nodes are only used to discover other nodes during startup and bootstrapping.
* This promotes high availability, and add a node at any time
* To know which replicas store a particular key `nodetool getendpoints <keyspace> <table> <partition-key>`
  * This returns list of IPs
* Replicas RF=3
  * The data goes to a primary replica, and the next two nodes in clockwise direction
  * `SimpleStrategy` - Pick N-1 nodes in clockwise direction
  * `NetworkTopologyStrategy` - Picks replicas across datacenters and racks
* Adding a new node redistributes the data
```bash
# Before state
root@cassandra4:/# nodetool getendpoints keyspace1 users 1;
172.19.0.4
172.19.0.3
172.19.0.2

# After adding the new node.
root@cassandra4:/# nodetool getendpoints keyspace1 users 1;
172.19.0.4
172.19.0.6 # Data moved from 2 to 6
172.19.0.3
```
* If the node is dead or needs removal, manually remove the node using nodetool. This rebalances the data
  * `nodetool decommission`
  * `nodetool removenode <host-id>`
  * `nodetool cleanup`

## Why?
* Massive write throughput
* Horizontal Scale
* HA and Fault Tolerance
* Append only SST Tables

## Use cases
* Logging, monitoring, metrics
* Time series

## Recipes
* Timestamped sensor data
```sql

CREATE KEYSPACE IF NOT EXISTS demo WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE demo;

CREATE TABLE sensor_readings (
    sensor_id text,
    reading_date date,          -- partition key (time bucket)
    reading_time timestamp,     -- clustering column (ordered by time)
    temperature float,
    humidity float,
    PRIMARY KEY ((sensor_id, reading_date), reading_time)
) WITH CLUSTERING ORDER BY (reading_time ASC);
```

* Chat
```sql

CREATE KEYSPACE IF NOT EXISTS chat_app WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE chat_app;

CREATE TABLE messages (
    chat_id text,            -- chat room or conversation id
    message_id timeuuid,     -- unique id for message, timeuuid orders by time
    sender_id text,
    message text,
    sent_at timestamp,
    PRIMARY KEY (chat_id, message_id)
) WITH CLUSTERING ORDER BY (message_id DESC);
```


## Why is it Column Family?
* Each row can have different columns. 
* The partition key is mandatory
* The columns can be ommitted during inserts -- Flexible columns so column family. Inherited from Google BigTable where it's row -> column1:value1 column2:value2 groups
* Sparse data storage. Save storage cost
* Can have millions of rows (theoretically lol)
```sql
-- Step 1: Create a keyspace
CREATE KEYSPACE demo
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

-- Step 2: Use the keyspace
USE demo;

-- Step 3: Create a table (column family)
CREATE TABLE users (
    user_id text PRIMARY KEY,
    name text,
    email text,
    age int
);

-- Insert user with full data
INSERT INTO users (user_id, name, email, age)
VALUES ('user123', 'Alice', 'alice@example.com', 30);

-- Insert user with only name and age (no email)
INSERT INTO users (user_id, name, age)
VALUES ('user456', 'Bob', 25);

UPDATE users
SET email = 'bob@example.com'
WHERE user_id = 'user456';
```


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
