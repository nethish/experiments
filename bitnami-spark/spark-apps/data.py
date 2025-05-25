from pyspark.sql import SparkSession
from pyspark.sql.functions import spark_partition_id, sum

# Create Spark session
spark = (
    SparkSession.builder.appName("CSV Partition Example")
    .master("local[*]")
    .config("spark.hadoop.hadoop.security.authentication", "simple")  # Disable Kerberos
    .config(
        "spark.hadoop.fs.file.impl", "org.apache.hadoop.fs.LocalFileSystem"
    )  # Force local FS
    .getOrCreate()
)

# Create a test CSV (you can replace this with a real path)
csv_path = "file:///opt/bitnami/spark/data/data.csv"

# Load CSV with multiple partitions
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
