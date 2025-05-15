kafka-topics --delete \
  --bootstrap-server localhost:9092 \
  --topic mongo.testdb.testcoll

kafka-topics --describe \
  --bootstrap-server localhost:9092 \
  --topic mongo.testdb.testcoll

