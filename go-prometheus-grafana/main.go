package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var gauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "random_number",
		Help: "A random number between 0 and 100",
		ConstLabels: map[string]string{
			"tagA": "1",
			"tagB": "2",
		},
	},
)

var counter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Namespace:   "",
		Subsystem:   "",
		Name:        "counter",
		Help:        "A counter",
		ConstLabels: map[string]string{},
	},
)

func main() {
	// (Optional) Set Go metrics collection explicitly
	prometheus.Register(prometheus.NewGoCollector())
	prometheus.Register(gauge)
	prometheus.Register(counter)

	go func() {
		for {
			gauge.Add(1)
			counter.Add(10)
			time.Sleep(2 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	// Your normal server code...
	go http.ListenAndServe(":8080", nil) // metrics endpoint on :8080/metrics

	fmt.Println("server started")

	select {} // Block forever
}
