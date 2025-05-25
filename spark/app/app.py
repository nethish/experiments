from pyspark.sql import SparkSession

spark = SparkSession.builder.appName("DistributedCSVReader").getOrCreate()

# Read the CSV (note: assumes the file is local to the container)
df = spark.read.option("header", "true").csv("/app/data.csv")

# Repartition explicitly across 2 partitions
df = df.repartition(2)


# Function to inspect partitions
def print_partition(index, iterator):
    print(f"=======================Partition {index}:")
    for row in iterator:
        print(row)
    yield


# Trigger partition-wise execution
df.rdd.mapPartitionsWithIndex(print_partition).collect()


spark.stop()
