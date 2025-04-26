package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// (Optional) Set Go metrics collection explicitly
	prometheus.Register(prometheus.NewGoCollector())

	http.Handle("/metrics", promhttp.Handler())

	// Your normal server code...
	go http.ListenAndServe(":8080", nil) // metrics endpoint on :8080/metrics

	fmt.Println("server started")

	select {} // Block forever
}
