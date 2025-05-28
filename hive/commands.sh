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

docker exec -it hive-client bash
beeline -u jdbc:hive2://hiveserver:10000/default -n hiveuser -p hivepassword
or
hive

SHOW DATABASES;
CREATE DATABASE my_first_hive_db;
USE my_first_hive_db;
CREATE TABLE IF NOT EXISTS users (id INT, name STRING) ROW FORMAT DELIMITED FIELDS TERMINATED BY ',';
LOAD DATA LOCAL INPATH '/opt/hive/examples/files/kv1.txt' OVERWRITE INTO TABLE users; -- This is an example, you might need to copy a file first or use HDFS path
SELECT * FROM users;
