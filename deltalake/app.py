from pyspark.sql import SparkSession
from delta.tables import DeltaTable
from pyspark.sql.functions import expr
import os

spark = (
    SparkSession.builder.appName("DeltaLakeDemo")
    .master("local[*]")
    .config("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension")
    .config(
        "spark.sql.catalog.spark_catalog",
        "org.apache.spark.sql.delta.catalog.DeltaCatalog",
    )
    .config("spark.jars.packages", "io.delta:delta-core_2.12:2.4.0")
    .getOrCreate()
)

# Clean old data if rerunning
path = "delta_table"
if os.path.exists(path):
    import shutil

    shutil.rmtree(path)

# Create and write Delta table
data = spark.createDataFrame([(1, "Alice"), (2, "Bob"), (3, "Carol")], ["id", "name"])

data.write.format("delta").save(path)

# Read the table back
df = spark.read.format("delta").load(path)
print("=== Initial data:")
df.show()

# Update the Delta table
delta_table = DeltaTable.forPath(spark, path)
delta_table.update(condition=expr("id = 2"), set={"name": expr("'Robert'")})

print("=== After update:")
delta_table.toDF().show()

# Delete a record
delta_table.delete("id = 3")

print("=== After delete:")
delta_table.toDF().show()

# Append a new row
new_data = spark.createDataFrame([(4, "Diana")], ["id", "name"])
new_data.write.format("delta").mode("append").save(path)

print("=== After append:")
spark.read.format("delta").load(path).show()

# Time travel: read older version
print("=== Time travel to version 0:")
old_df = spark.read.format("delta").option("versionAsOf", 0).load(path)
old_df.show()

spark.stop()
