curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" -d '{
  "name": "mongodb-source-connector",
  "config": {
    "connector.class": "com.mongodb.kafka.connect.MongoSourceConnector",
    "tasks.max": "1",
    "connection.uri": "mongodb://mongo:27017",
    "database": "testdb",
    "collection": "testcoll",
    "topic.prefix": "mongo.",
    "output.format.value": "json"
  }
}'

