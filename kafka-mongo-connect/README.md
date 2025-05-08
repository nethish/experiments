# Mongo CDC into Kafka Topic with Kafka Connect
* Create Mongo Cluster with replica set.
* Network `mongoCluster` (use the script from mongocdc folder)
* Register the connector with `create-mongo-connector.sh`
* Run the `consume.sh` file to see consumed events.
  * The topic name somehow is mongo..testdb.testcoll (where did two dots come from bruh. Spent like an hour)
