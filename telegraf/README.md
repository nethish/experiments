# telegraf

Publish metrics to influxdb (a time series database). telegraf is an opensource metric exporter.

```bash
# After docker compose up

docker exec -it influxdb influx

# You should initialize the client
docker exec -it influxdb influx config create \
  --config-name local \
  --host-url http://localhost:8086 \
  --org myorg \
  --token my-super-secret-token \
  --active

USE telegraf;
SHOW MEASUREMENTS;
SELECT * FROM cpu LIMIT 5;

# In grafana -> Add Connections -> Datasources -> InfluxDB
# Use http://influxdb:8086 as the URL
# Refer to influxdb configuration in docker compose while setting up the Datasources
# * Organization = myorg
# * User = admin
# * Password = admin123
# * Default Bucket = telegraf
# * Token = my-super-secret-token
```


```bash
# Example query that you can run in grafana or `influx query 'query'`
from(bucket: "telegraf")
  |> range(start: -5m)
  |> filter(fn: (r) => r._measurement == "cpu")
  |> filter(fn: (r) => r._field == "usage_user")
  |> sort(columns: ["_value"], desc: true)
  |> limit(n:5)
```
