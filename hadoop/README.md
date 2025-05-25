# Hadoop
* Man this is complicated like hell
* Namenode is like the master node that receives commands. Knows where the file is and stores metadata
* Datanodes handle the actual storage of data
* I copied the `data.parquet` file generated as part of `spark/app` experiment into `./hadoop_namenode/` and copied it into the cluster. The file was split into 134217728 bytes parts and spread across two nodes
* `hdfs dfs` is a wrapper around hadoops file system API. `ls, cat, rm, mkdir, get, put`
