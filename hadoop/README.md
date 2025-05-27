# Hadoop
* Man this is complicated like hell
* Namenode is like the master node that receives commands. Knows where the file is and stores metadata
* Datanodes handle the actual storage of data
* I copied the `data.parquet` file generated as part of `spark/app` experiment into `./hadoop_namenode/` and copied it into the cluster. The file was split into 134217728 bytes parts and spread across two nodes
* `hdfs dfs` is a wrapper around hadoops file system API. `ls, cat, rm, mkdir, get, put`
* Look at `commands.sh` to understand how to put file in the file system
* View the UI at `http://localhost:9870/`

## hdfs cli
* Download hdfs `brew install hadoop`
* Export the `env.sh`
* `hdfs --config ./config dfs -ls /`


## Connect to Spark Cluster
* Make the below changes to spark docker-compose
  * Create a network
```
    hadoop-net:
      external:
        name: hadoop_hadoop-net
```
  * Add the `hadoop-net` to all the services 
  * Start the services
  * Use the `app/hadoop/main.py` to run the job using `spark-submit`
  * Before that, copy the data.parquet to hadoop file system
  * It's slower than FS lookup. You probably need hundreds of GBs to Terabytes of data to see real performance
