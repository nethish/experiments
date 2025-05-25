from pyspark.sql import SparkSession
from pyspark.sql.types import StructType, StructField, IntegerType, StringType

# Create Spark session
spark = (
    SparkSession.builder.appName("StructType Example").master("local[*]").getOrCreate()
)

# Define schema explicitly
schema = StructType(
    [
        StructField("id", IntegerType(), nullable=False),
        StructField("name", StringType(), nullable=True),
        StructField("age", IntegerType(), nullable=True),
    ]
)

# Sample data as list of tuples
data = [
    (1, "Alice", 30),
    (2, "Bob", 25),
    (3, "Charlie", None),
]

# Create DataFrame with schema
df = spark.createDataFrame(data, schema)

# Show the DataFrame
df.show()

# Print schema to verify
df.printSchema()

spark.stop()
