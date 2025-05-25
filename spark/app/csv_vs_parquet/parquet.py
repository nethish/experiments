from pyspark.sql import SparkSession
import time

spark = SparkSession.builder.appName("ParquetJob").getOrCreate()

start = time.time()

df = spark.read.parquet("/app/csv_vs_parquet/data.parquet")
result = df.groupBy("category").count()
result.show()

print(f"Parquet job took {time.time() - start:.2f} seconds")

spark.stop()
