from pyspark.sql import SparkSession
from pyspark.sql.functions import from_json, col, window
from pyspark.sql.types import (
    StructType,
    StructField,
    StringType,
    IntegerType,
    TimestampType,
)

spark = SparkSession.builder.appName("KafkaWatermarkExample").getOrCreate()

# Define schema including eventTime as TimestampType
schema = StructType(
    [
        StructField("userId", StringType()),
        StructField("action", StringType()),
        StructField("value", IntegerType()),
        StructField("eventTime", StringType()),  # initially string
    ]
)

# Read stream from Kafka
kafka_df = (
    spark.readStream.format("kafka")
    .option("kafka.bootstrap.servers", "kafka:9092")
    .option("subscribe", "events-topic")
    .option("startingOffsets", "earliest")
    .load()
)

# Convert value bytes to string, parse JSON
json_df = kafka_df.selectExpr("CAST(value AS STRING) as json_str")

parsed_df = json_df.select(from_json(col("json_str"), schema).alias("data")).select(
    "data.*"
)

# Convert eventTime string to timestamp
from pyspark.sql.functions import to_timestamp

events_df = parsed_df.withColumn(
    "eventTime", to_timestamp(col("eventTime"), "yyyy-MM-dd'T'HH:mm:ssX")
)

# Add watermark with 1 minute delay to handle late data
watermarked_df = events_df.withWatermark("eventTime", "1 minute")

# Windowed aggregation: count events per 1 minute window
windowed_counts = watermarked_df.groupBy(
    window(col("eventTime"), "1 minute"), col("action")
).count()

# Output to console
# ==================
query = (
    windowed_counts.writeStream.outputMode("append")
    .format("console")
    .option("truncate", False)
    .start()
)


from pyspark.sql.functions import to_json, struct

output_df = windowed_counts.select(
    to_json(
        struct(
            col("window.start").alias("window_start"),
            col("window.end").alias("window_end"),
            col("action"),
            col("count"),
        )
    ).alias("value")
)

# Kafka expects key and value as bytes. You can omit key or add one if needed.
# Here we omit key by default (null key).

query = (
    output_df.writeStream.format("kafka")
    .option("kafka.bootstrap.servers", "kafka:9092")
    .option("topic", "output-topic")
    .option(
        "checkpointLocation", "/tmp/spark-checkpoints/output-topic"
    )  # Required for Kafka sink
    .outputMode("update")  # Or "append" if your aggregation supports it
    .start()
)

query.awaitTermination()
