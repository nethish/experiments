from pyspark.sql import SparkSession
from pyspark.sql.functions import expr
import os


warehouse_path = f"file://{os.path.abspath('./tmp')}"
print("Warehouse path:", warehouse_path)

# Step 1: Initialize Spark with Iceberg support
spark = (
    SparkSession.builder.appName("IcebergDemo")
    .config(
        "spark.jars.packages", "org.apache.iceberg:iceberg-spark-runtime-3.4_2.12:1.5.0"
    )
    .config(
        "spark.sql.extensions",
        "org.apache.iceberg.spark.extensions.IcebergSparkSessionExtensions",
    )
    .config("spark.sql.catalog.local", "org.apache.iceberg.spark.SparkCatalog")
    .config("spark.sql.catalog.local.type", "hadoop")
    .config("spark.sql.catalog.local.warehouse", warehouse_path)
    .getOrCreate()
)

# Step 2: Create a new Iceberg table
spark.sql("DROP TABLE IF EXISTS local.db.people")
spark.sql("""
    CREATE TABLE local.db.people (
        id INT,
        name STRING,
        age INT
    )
    USING iceberg
""")

# Step 3: Insert data
spark.sql("""
    INSERT INTO local.db.people VALUES
    (1, 'Alice', 30),
    (2, 'Bob', 25),
    (3, 'Charlie', 35)
""")

# Step 4: Read the table
print("✅ Table contents:")
spark.sql("SELECT * FROM local.db.people").show()

# Step 5: Update rows
spark.sql("""
    UPDATE local.db.people
    SET age = age + 1
    WHERE name = 'Alice'
""")

# Step 6: Delete rows
spark.sql("""
    DELETE FROM local.db.people
    WHERE name = 'Bob'
""")

# Step 7: Final data
print("✅ After update and delete:")
spark.sql("SELECT * FROM local.db.people").show()

# Step 8: View table history
print("✅ Table history (snapshots):")
spark.sql("SELECT * FROM local.db.people.snapshots").show(truncate=False)

spark.stop()
