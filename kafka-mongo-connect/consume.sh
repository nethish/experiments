# Run Kafka console consumer
docker exec -it $(docker ps -qf "ancestor=confluentinc/cp-kafka:7.5.0") \
  kafka-console-consumer --bootstrap-server localhost:9092 --topic mongo.testdb.testcoll --from-beginning

