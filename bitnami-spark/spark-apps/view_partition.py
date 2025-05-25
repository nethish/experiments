from pyspark import SparkConf, SparkContext

conf = SparkConf().setAppName("PartitionViewer").setMaster("local[*]")
sc = SparkContext(conf=conf)

# Create RDD with 4 partitions
rdd = sc.parallelize(range(1, 21), numSlices=4)


# Inspect partitions
def inspect_partitions(index, items):
    return [f"Partition {index}: {list(items)}"]


partitioned_data = rdd.mapPartitionsWithIndex(inspect_partitions).collect()

for entry in partitioned_data:
    print(entry)

sc.stop()
