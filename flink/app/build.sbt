name := "flink-scala-kafka"

version := "0.1"

scalaVersion := "2.12.17"

libraryDependencies ++= Seq(
  "org.apache.flink" %% "flink-scala" % "1.17.0" % "provided",
  "org.apache.flink" %% "flink-streaming-scala" % "1.17.0" % "provided",
  "org.apache.flink" % "flink-connector-kafka" % "1.17.0"
)
