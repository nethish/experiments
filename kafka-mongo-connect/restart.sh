cd smt-encrypt/kafka-encrypt/
mvn clean package
cd ../..
cp smt-encrypt/kafka-encrypt/target/kafka-encrypt-1.0-SNAPSHOT.jar connect-plugins
docker-compose restart kafka-connect
sleep 20
bash delete-connector.sh && bash transforms.sh
