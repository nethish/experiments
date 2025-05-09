curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" -d '{
  "name": "mongodb-source-connector",
  "config": {
    "connector.class": "com.mongodb.kafka.connect.MongoSourceConnector",
    "tasks.max": "1",
    "connection.uri": "mongodb://mongo1:27017,mongo2:27017,mongo3:27017/?replicaSet=myReplicaSet",
    "database": "testdb",
    "collection": "testcoll",
    "topic.prefix": "mongo",
    "output.format.value": "json"
  }
}'
