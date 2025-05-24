from pyspark.sql import SparkSession
from pyspark.sql.types import *
from decimal import Decimal
from datetime import date, datetime

spark = (
    SparkSession.builder.appName("AllSparkSQLTypesExample")
    .master("local[*]")
    .getOrCreate()
)

# Define schema with various Spark SQL data types
schema = StructType(
    [
        StructField("string_col", StringType(), True),
        StructField("integer_col", IntegerType(), True),
        StructField("long_col", LongType(), True),
        StructField("float_col", FloatType(), True),
        StructField("double_col", DoubleType(), True),
        StructField("decimal_col", DecimalType(10, 2), True),
        StructField("boolean_col", BooleanType(), True),
        StructField("date_col", DateType(), True),
        StructField("timestamp_col", TimestampType(), True),
        StructField("binary_col", BinaryType(), True),
        StructField("array_col", ArrayType(StringType()), True),
        StructField("map_col", MapType(StringType(), IntegerType()), True),
        StructField(
            "struct_col",
            StructType(
                [
                    StructField("nested_str", StringType(), True),
                    StructField("nested_int", IntegerType(), True),
                ]
            ),
            True,
        ),
    ]
)

# Sample data matching schema
data = [
    (
        "hello",
        42,
        1234567890123,
        3.14,
        2.718281828,
        Decimal("12345.67"),
        True,
        date(2024, 5, 25),
        datetime(2024, 5, 25, 12, 34, 56),
        bytearray(b"binarydata"),
        ["a", "b", "c"],
        {"key1": 1, "key2": 2},
        ("nested", 99),
    )
]

# Create DataFrame with schema
df = spark.createDataFrame(data, schema=schema)

# Show data and schema
df.show(truncate=False)
df.printSchema()

spark.stop()
