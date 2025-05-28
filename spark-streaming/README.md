# Spark Streaming
* Spark Streaming breaks incoming data into small batches (e.g., every 2 seconds), and then processes each batch using regular Spark jobs.
* Triggers - Spark will continuously gather data as it arrives, but it will process and emit results only once every minute.
* Watermarking with event time example - I won't accept messages with time X as I have seen messages with time >= X+T already. T is the tolerance window
  * Window start is inclusive and window end is Exclusive
* Checkpointing - Checkpointing in Spark Streaming is a fault-tolerance mechanism that helps Spark recover from failures during stream processing.
  * Saves data periodically to recover in case of failure

## Recipes
* Micro batch processing. See main.py
* Watermarking - See watermark.py
  * outputMode("complete") - Will push messages only after window is complete
  * outputMode("update") - Will push messages for every update
* Checkpointing
  * Spark tracks which output messages have been committed to the output Kafka topic.
  * This ensures exactly-once delivery guarantees to the sink.
