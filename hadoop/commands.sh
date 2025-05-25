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
