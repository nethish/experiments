curl -X POST http://localhost:8083/connectors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "mongo-sink-connector",
    "config": {
      "connector.class": "com.mongodb.kafka.connect.MongoSinkConnector",
      "tasks.max": "1",
      "topics": "mongo.testdb.testcoll",

      "connection.uri": "mongodb://mongo1:27017,mongo2:27017,mongo3:27017/?replicaSet=myReplicaSet",

      "database": "testdb",
      "collection": "history_sink",

      "key.converter": "org.apache.kafka.connect.storage.StringConverter",
      "value.converter": "org.apache.kafka.connect.json.JsonConverter",
      "value.converter.schemas.enable": "false",

      "document.id.strategy": "com.mongodb.kafka.connect.sink.processor.id.strategy.ProvidedInKeyStrategy"
    }
  }'
