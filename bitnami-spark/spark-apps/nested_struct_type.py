from pyspark.sql import SparkSession
from pyspark.sql.functions import col

spark = (
    SparkSession.builder.appName("Flatten Nested Dict").master("local[*]").getOrCreate()
)

# Data as list of dicts with nested dicts
data = [
    {
        "id": 1,
        "person": {"name": "Alice", "address": {"city": "New York", "state": "NY"}},
    },
    {
        "id": 2,
        "person": {"name": "Bob", "address": {"city": "Los Angeles", "state": "CA"}},
    },
]

# Define schema (optional but recommended for nested structures)
from pyspark.sql.types import StructType, StructField, IntegerType, StringType

schema = StructType(
    [
        StructField("id", IntegerType(), False),
        StructField(
            "person",
            StructType(
                [
                    StructField("name", StringType(), True),
                    StructField(
                        "address",
                        StructType(
                            [
                                StructField("city", StringType(), True),
                                StructField("state", StringType(), True),
                            ]
                        ),
                        True,
                    ),
                ]
            ),
            True,
        ),
    ]
)

# Create DataFrame with schema
df = spark.createDataFrame(data, schema)

print("Original nested DataFrame:")
df.show(truncate=False)
df.printSchema()

# Flatten nested fields
flat_df = df.select(
    col("id"),
    col("person.name").alias("name"),
    col("person.address.city").alias("city"),
    col("person.address.state").alias("state"),
)

print("Flattened DataFrame:")
flat_df.show(truncate=False)
flat_df.printSchema()

spark.stop()
