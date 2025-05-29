import org.apache.flink.api.common.state.{ValueState, ValueStateDescriptor}
import org.apache.flink.configuration.Configuration
import org.apache.flink.streaming.api.functions.KeyedProcessFunction
import org.apache.flink.streaming.api.scala._
import org.apache.flink.streaming.connectors.kafka.FlinkKafkaConsumer
import org.apache.flink.api.common.serialization.SimpleStringSchema
import org.apache.flink.util.Collector
import java.util.Properties

case class UserEvent(user_id: Int, event: String, timestamp: String)

object StatefulUserEventCounter {
  def main(args: Array[String]): Unit = {
    val env = StreamExecutionEnvironment.getExecutionEnvironment

    // Kafka config
    val props = new Properties()
    props.setProperty("bootstrap.servers", "kafka:9092")
    props.setProperty("group.id", "flink-user-counter")

    val consumer = new FlinkKafkaConsumer[String]("user-events", new SimpleStringSchema(), props)
    val stream = env.addSource(consumer)

    val parsed: DataStream[UserEvent] = stream
      .filter(_.nonEmpty)
      .map { json =>
        val pattern = """\{"user_id":(\d+),"event":"(.*?)","timestamp":"(.*?)"\}""".r
        val pattern(userId, event, timestamp) = json
        UserEvent(userId.toInt, event, timestamp)
      }

    // Key by user_id and apply stateful processing
    val result = parsed
      .keyBy(_.user_id)
      .process(new EventCountFunction)

    result.print()

    env.execute("Stateful Event Count Per User")
  }
}

// Stateful counter logic
class EventCountFunction extends KeyedProcessFunction[Int, UserEvent, String] {

  private var countState: ValueState[Int] = _

  override def open(parameters: Configuration): Unit = {
    val descriptor = new ValueStateDescriptor[Int]("eventCount", classOf[Int])
    countState = getRuntimeContext.getState(descriptor)
  }

  override def processElement(
      value: UserEvent,
      ctx: KeyedProcessFunction[Int, UserEvent, String]#Context,
      out: Collector[String]
  ): Unit = {
    val currentCount = Option(countState.value()).getOrElse(0)
    val newCount = currentCount + 1
    countState.update(newCount)

    out.collect(s"User ${value.user_id} has ${newCount} events so far. Latest: ${value.event}")
  }
}
