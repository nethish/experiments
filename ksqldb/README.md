# ksqldb
* An SQL based stream processor developed by confluent
* It uses underlying Kafka as the storage.
* It streams, transforms, aggregates and pushes output to another kafka topic
* Why should you use it? You don't have a Java team to use Flink and you have SQL teams

## Features
* Stream Join with time window
* Materialized views (pushes it into kafka topic)
* UDFs and Custom  UDFs
* https://github.com/confluentinc/ksql


## Setup

```bash
docker-compose exec ksqldb-cli ksql http://ksqldb-server:8088

CREATE STREAM SALES_STREAM (
  region STRING,
  amount DOUBLE
) WITH (
  KAFKA_TOPIC='sales',
  VALUE_FORMAT='JSON'
);

SHOW STREAMS;

CREATE TABLE sales_aggregated AS
  SELECT
    region,
    WINDOWSTART AS window_start,
    WINDOWEND AS window_end,
    SUM(amount) AS total_sales
  FROM SALES_STREAM
  WINDOW TUMBLING (SIZE 30 SECONDS)
  GROUP BY region
  EMIT CHANGES;


kafka-console-producer --bootstrap-server localhost:9092 --topic sales
{"region": "north", "amount": 100}

kafka-console-consumer --bootstrap-server localhost:9092 --topic SALES_AGGREGATED --from-beginning



```

```
```
