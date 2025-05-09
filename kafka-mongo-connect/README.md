# Mongo CDC into Kafka Topic with Kafka Connect
* Create Mongo Cluster with replica set.
* Network `mongoCluster` (use the script from mongocdc folder)
* Register the connector with `register-mongo-connector.sh`
* Run the `consume.sh` file to see consumed events.
  * The topic name somehow is mongo.testdb.testcoll
