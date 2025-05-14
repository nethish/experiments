# Confluent Schema Registry
* You can register Avro, Protobuf and JSON schema in the schema registry
* You register the schema, and CSR returns an ID
* For each message [0 | Schema ID | Serialized Message] is sent to the topic
* On the consumer side, you get the schema from schema ID, and deserialize appropriately
* This folder demonstrates example for Avro
* For each new schema registered (even with same subject), the Schema ID will be incremented by 1. i.e there is always a unique id for each new schema
* You can distribute libraries that integrate with CSR and encode/ decode the message appropriately
