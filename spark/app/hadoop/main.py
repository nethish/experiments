from pyspark.sql import SparkSession
import time

start = time.time()

spark = SparkSession.builder.appName("HDFS Example").getOrCreate()

df = spark.read.parquet(
    "hdfs://namenode:9000/user/demo/data.parquet",
    header=True,
    inferSchema=True,
)


# df = spark.read.parquet("/app/csv_vs_parquet/data.parquet")
result = df.groupBy("category").count()
result.show()

print(f"Parquet job took {time.time() - start:.2f} seconds")

spark.stop()
