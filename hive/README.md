# Hive
* In addition to the hadoop cluster, postgres (for metastore), metastore server itself, and hiveserver are the three other components of Hive

## Hiveserver2
* First thing, you can connect to this server using DBeaver
* Standard interface with JDBC, ODBC, and thrift
* Authentication, Authorization
* Query Compilation, Execution and Orchestration
  * Consults Hive Metastore and verifies table name, column etc
  * Optimization - Partition pruning, predicate pushdown, aggregation
  * Physical Plan Generation - This plan is a Directed Acyclic Graph (DAG) consisting of multiple stages (e.g., MapReduce jobs, Tez tasks, Spark jobs) that will run on the Hadoop cluster.
* Query Execution
  * Takes the DAG of Jobs/ Tasks
  * Interact with Hadoop YARN (Yet another resource negotiator)
  * YARN allocates resources and launches individual tasks on Data nodes
  * The tasks read data from HDFS
* Results are gathered and shown

## Beeline
* CLI Utility to interact with Hive, built on top of JDBC
* You use Beeline to type and execute the HiveQL, so it's just a tool to connect to the server

## Metastore
* It stores the names and properties of the databases created.
* Tables: Names, column definitions (schema), data types, partitioning keys, bucketing information, storage formats (e.g., TextFile, ORC, Parquet), and SerDe (Serializer/Deserializer) properties for all tables.
* Metadata about partitions and locations in HDFS
* Hive Metastore exposes its services via Thrift Interface
* Thrift uses Interface Definition Language much like our Protocol Buffers and have built in mechanism for generating rpc stubs
* Some important table described by Gemini -- Connect to postgres with username password from dockerfile
```SQL
SELECT *
FROM public."DBS";

select *
from public."TBLS";

select *
from public."SDS";

select *
from public."COLUMNS_V2";

select *
from public."PARTITIONS" p;

```




## Create a database
* Initialize metastore first
```bash
docker exec -it hive-metastore bash
schematool -initSchema -dbType postgres
```

* Exec into the hiveserver and connect using beeline or hive
```bash
docker exec -it hiveserver bash
beeline -u jdbc:hive2://hiveserver:10000/default -n hiveuser -p hivepassword
# You can optionally use `hive` command
```

* Database commands
```SQL
SHOW DATABASES;
CREATE DATABASE first_database;
USE first_database;

CREATE TABLE users (
  id INT,
  fruit STRING
)
ROW FORMAT DELIMITED
FIELDS TERMINATED BY '\t'
STORED AS TEXTFILE;

-- Don't forget to copy the kv1.txt available locally
LOAD DATA LOCAL INPATH '/opt/hive/data/kv1.txt' OVERWRITE INTO TABLE users;

SELECT * FROM users;


```
