docker exec -it namenode bash

hdfs dfsadmin -report
hdfs dfs -ls /

hdfs dfs -mkdir /user
hdfs dfs -mkdir /user/test
echo "Hello Hadoop!" >hello.txt
cd /hadoop/dfs/name
hdfs dfs -put hello.txt /user/test/
hdfs dfs -cat /user/test/hello.txt

# Try to put a very big parquet file inside the cluster. See it get split
#

docker exec -it hive-metastore bash
schematool -initSchema -dbType postgres

docker exec -it hiveserver bash
beeline -u jdbc:hive2://hiveserver:10000/default -n hiveuser -p hivepassword
or
hive

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

LOAD DATA LOCAL INPATH '/opt/hive/data/kv1.txt' OVERWRITE INTO TABLE users;

SELECT * FROM users;

# PARQUET
```SQL
CREATE EXTERNAL TABLE IF NOT EXISTS external_products_parquet (
  id INT,
  value INT,          -- Assuming 'value' is your second column, not 'name' or 'category'
  category STRING     -- Assuming 'category' is your third column
)
STORED AS PARQUET
LOCATION 'hdfs://namenode:9000/files/';



SELECT COUNT(1) FROM external_products_parquet;
```
