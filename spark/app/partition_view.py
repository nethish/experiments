from pyspark.sql import SparkSession
from pyspark.sql.functions import spark_partition_id, sum

spark = SparkSession.builder.appName("DistributedCSVReader").getOrCreate()

csv_path = "/app/data.csv"
df = spark.read.csv(csv_path, header=True, inferSchema=True).repartition(3)

# Add partition ID for inspection
df_with_pid = df.withColumn("partition_id", spark_partition_id())

# Show how data is distributed
print("📦 Data spread across partitions:")
df_with_pid.orderBy("partition_id").show(truncate=False)

# Group by category and sum values
agg_df = df.groupBy("category").agg(sum("value").alias("total_value"))

print("📊 Aggregated results:")
agg_df.show()

# Optional: show number of rows per partition
print("📊 Rows per partition:")
df_with_pid.groupBy("partition_id").count().show()

spark.stop()
