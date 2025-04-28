aws dynamodb create-table \
  --table-name People \
  --attribute-definitions AttributeName=ID,AttributeType=S \
  --key-schema AttributeName=ID,KeyType=HASH \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
  --endpoint-url http://localhost:8000

aws dynamodb put-item \
  --table-name People \
  --item '{"ID": {"S": "1"}, "Name": {"S": "Alice"}, "Age": {"N": "30"}}' \
  --endpoint-url http://localhost:8000

aws dynamodb get-item \
  --table-name People \
  --key '{"ID": {"S": "1"}}' \
  --endpoint-url http://localhost:8000

aws dynamodb scan --table-name People --select "COUNT" --endpoint-url http://localhost:8000
