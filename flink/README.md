# Flink 
* Flink is a stream processor framework

## How to run
* Go to app and run `sbt clean build` to build the `/app/target/scala-<version>/<project-name>.jar`
* The Flink libs only work with Java 11, so export both `JAVA_HOME` and the PATH to Java11 bin
* Inside the job manager run `flink run /app/target/scala-2.12/flink-scala-kafka_2.12-0.1.jar`
* The logs will be printed inside `taskmanager`

## Kafka 
```bash
kafka-console-producer --broker-list kafka:9092 --topic input-topic
```

```json
{"user_id":123,"event":"login","timestamp":"2024-05-29T10:00:00Z"}
```
