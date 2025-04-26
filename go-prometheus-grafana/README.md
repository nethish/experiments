# Go Prometheus Grafana to monitor Go runtime metrics
```
docker-compose up -d
go run ./main.go
```

* Open prometheus at `localhost:9090`
* Open Grafana at `localhost:3000`
* Open your app metrics at `localhost:8080`
* In the prometheus.yml config, the hostname is set to `host.docker.internal` which resolves to local host of the machine

## Prometheus
* Scrapes metrics every X internal configured in `prometheus.yml`
* Use PromQL to query data, aggregate or group

## Grafana
* Add prometheus as source
  * Usually you point to `localhost:9090` where prometheus is running
* Import Pre Build Dashboard from `https://grafana.com/grafana/dashboards`
  * Search for go runtime
  * Copy paste the id of the dashboard to import
