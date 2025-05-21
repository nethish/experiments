package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	token  = "demo-token"
	org    = "demo-org"
	bucket = "demo-bucket"
	url    = "http://localhost:8086"
)

func main() {
	client := influxdb2.NewClient(url, token)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(org, bucket)

	start := time.Now()
	for i := range 100000 {
		// 1. Write a point
		p := influxdb2.NewPoint("temperature",
			map[string]string{"location": "office"},
			map[string]interface{}{"value": rand.Float32() * 100},
			time.Now())
		err := writeAPI.WritePoint(context.Background(), p)
		if err != nil {
			panic(err)
		}

		if i%10000 == 0 {
			fmt.Println("Wrote 10000 records.")
		}
	}
	fmt.Println("✅ Data written in", time.Since(start))

	// 2. Query data
	query := `from(bucket:"demo-bucket")
	|> range(start: -1h)
	|> filter(fn: (r) => r._measurement == "temperature")`

	queryAPI := client.QueryAPI(org)
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		panic(err)
	}
	for result.Next() {
		fmt.Printf("🔎 Got: time=%s location=%s value=%v\n",
			result.Record().Time(), result.Record().ValueByKey("location"), result.Record().Value())
	}
	if result.Err() != nil {
		panic(result.Err())
	}
}
