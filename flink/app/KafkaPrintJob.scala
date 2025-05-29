// src/main/scala/KafkaPrintJob.scala

import org.apache.flink.api.common.serialization.SimpleStringSchema
import org.apache.flink.streaming.api.scala._
import org.apache.flink.streaming.connectors.kafka.FlinkKafkaConsumer

import java.util.Properties

object KafkaPrintJob {
  def main(args: Array[String]): Unit = {
    val env = StreamExecutionEnvironment.getExecutionEnvironment

    val props = new Properties()
    props.setProperty("bootstrap.servers", "kafka:9092")
    props.setProperty("group.id", "flink-scala-consumer")

    val kafkaStream = new FlinkKafkaConsumer[String](
      "input-topic",
      new SimpleStringSchema(),
      props
    )

    val stream = env.addSource(kafkaStream)

    stream.print()

    env.execute("Scala Kafka Consumer Job")
  }
}
