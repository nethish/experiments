from pyspark.sql import SparkSession
from pyspark.sql.functions import from_json, col, window, to_timestamp
from pyspark.sql.types import (
    StructType,
    StringType,
    IntegerType,
    FloatType,
)

# Create SparkSession
spark = SparkSession.builder.appName("KafkaSparkStructuredStreaming").getOrCreate()
spark.sparkContext.setLogLevel("WARN")

# Kafka config
kafka_bootstrap_servers = "kafka:9092"
kafka_topic = "micro-batch"  # change this to your Kafka topic

# Define schema for JSON data
schema = (
    StructType()
    .add("time", StringType())
    .add("code", StringType())
    .add("int", IntegerType())
    .add("float", FloatType())
)

# Read from Kafka
kafka_df = (
    spark.readStream.format("kafka")
    .option("kafka.bootstrap.servers", kafka_bootstrap_servers)
    .option("subscribe", kafka_topic)
    .option("startingOffsets", "latest")
    .load()
)

# Cast the binary "value" column to string
json_df = kafka_df.selectExpr("CAST(value AS STRING) as json_str")

# Parse JSON using schema
parsed_df = json_df.select(from_json(col("json_str"), schema).alias("data")).select(
    "data.*"
)

# Convert string "time" to TimestampType
final_df = parsed_df.withColumn(
    "time", to_timestamp(col("time"), "yyyy-MM-dd'T'HH:mm:ss")
)

# Optional: Print schema to confirm datatypes (prints once at startup)
final_df.printSchema()

# Group by 1-minute windows based on the "time" field
# windowed_df = final_df.groupBy(window(col("time"), "1 minute")).agg(
# For demo, just collect all rows as list for each window (optional)
# or just count or show first rows, here we just pass through
# )

# Instead of grouping (if you want raw data every 1 min), just trigger every 1 min:
query = (
    final_df.writeStream.outputMode("append")
    .format("console")
    .option("truncate", "false")
    .trigger(processingTime="1 minute")
    .start()
)

query.awaitTermination()
