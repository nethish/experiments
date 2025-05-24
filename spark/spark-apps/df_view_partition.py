from pyspark.sql import SparkSession
from pyspark.sql.functions import spark_partition_id

# Start SparkSession
spark = (
    SparkSession.builder.appName("DataFramePartitionView")
    .master("local[*]")
    .getOrCreate()
)

# Create a DataFrame with 4 partitions
df = spark.range(1, 21).repartition(4)

# Add partition ID column
df_with_pid = df.withColumn("partition_id", spark_partition_id())

# View contents grouped by partition
df_with_pid.orderBy("partition_id", "id").show(100, truncate=False)

# Optional: group to view count per partition
df_with_pid.groupBy("partition_id").count().show()

spark.stop()
