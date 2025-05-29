# Flink 
* Flink is a stream processor framework


## Components
* Job manager - Master Node (Control Plane)
  * Accepts and schedules flink jobs
  * Divides jobs into tasks and distributes them to task managers
* Task Manager - Worker node
  * Manages task slots and executors

## How to run
* Go to app and run `sbt clean package` to build the `/app/target/scala-<version>/<project-name>.jar`
* The Flink libs only work with Java 11, so export both `JAVA_HOME` and the PATH to Java11 bin
* Inside the job manager run `flink run /app/target/scala-2.12/flink-scala-kafka_2.12-0.1.jar`
* The logs will be printed inside `taskmanager`

## Kafka 
* For the KafkaPrintJob example
```bash
kafka-console-producer --broker-list kafka:9092 --topic input-topic
```

```json
{"user_id":123,"event":"login","timestamp":"2024-05-29T10:00:00Z"}
```

* For the Kafka Stateful Processing example
```bash

flink run -c StatefulUserEventCounter /app/target/scala-2.12/flink-scala-kafka_2.12-0.1.jar

kafka-console-producer --broker-list kafka:9092 --topic user-events

{"user_id":123,"event":"login","timestamp":"2024-05-29T10:00:00Z"}
{"user_id":123,"event":"click","timestamp":"2024-05-29T10:01:00Z"}
{"user_id":456,"event":"logout","timestamp":"2024-05-29T10:02:00Z"}
```

## Dependencies
Download deps from
```bash
curl -LO https://repo1.maven.org/maven2/org/apache/flink/flink-connector-kafka/1.17.0/

curl -LO https://repo1.maven.org/maven2/org/apache/kafka/kafka-clients/2.8.0/
```

