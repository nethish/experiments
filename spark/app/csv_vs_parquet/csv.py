from pyspark.sql import SparkSession
import time

spark = SparkSession.builder.appName("CSVJob").getOrCreate()

start = time.time()

df = spark.read.option("header", "true").csv("/app/csv_vs_parquet/data.csv")
result = df.groupBy("category").count()
result.show()

print(f"CSV job took {time.time() - start:.2f} seconds")

spark.stop()
