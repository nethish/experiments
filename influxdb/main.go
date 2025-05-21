package main

import (
	"context"
	"fmt"
	"log"
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
	for i := range 10000 {
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

	queryAPI := client.QueryAPI(org)

	// // 2. Query data
	// query := `from(bucket:"demo-bucket")
	// |> range(start: -1h)
	// |> filter(fn: (r) => r._measurement == "temperature")`
	//
	// result, err := queryAPI.Query(context.Background(), query)
	// if err != nil {
	// 	panic(err)
	// }
	// for result.Next() {
	// 	fmt.Printf("🔎 Got: time=%s location=%s value=%v\n",
	// 		result.Record().Time(), result.Record().ValueByKey("location"), result.Record().Value())
	// }
	// if result.Err() != nil {
	// 	panic(result.Err())
	// }

	// Flux query combining min, max, avg over 10-minute windows
	fluxQuery := `
data = from(bucket:"` + bucket + `")
  |> range(start: -1h)
  |> filter(fn: (r) => r._measurement == "temperature")
  |> aggregateWindow(every: 10m, fn: mean, createEmpty: false)
  |> rename(columns: {_value: "avg"})

maxData = from(bucket:"` + bucket + `")
  |> range(start: -1h)
  |> filter(fn: (r) => r._measurement == "temperature")
  |> aggregateWindow(every: 10m, fn: max, createEmpty: false)
  |> rename(columns: {_value: "max"})

minData = from(bucket:"` + bucket + `")
  |> range(start: -1h)
  |> filter(fn: (r) => r._measurement == "temperature")
  |> aggregateWindow(every: 10m, fn: min, createEmpty: false)
  |> rename(columns: {_value: "min"})

join1 = join(tables: {avg: data, max: maxData}, on: ["_time", "_field"])
result = join(tables: {join1: join1, min: minData}, on: ["_time", "_field"])
  |> map(fn: (r) => ({
    _time: r._time,
    _field: r._field,
    min: r.min,
    max: r.max,
    avg: r.avg
  }))

result
`

	// Execute query
	result, err := queryAPI.Query(context.Background(), fluxQuery)
	if err != nil {
		log.Fatalf("Query error: %v", err)
	}

	// Iterate over query response
	for result.Next() {
		// Access values by column name
		record := result.Record()
		fmt.Printf("Time: %v, Field: %s, Min: %v, Max: %v, Avg: %v\n",
			record.Time(),
			record.ValueByKey("_field"),
			record.ValueByKey("min"),
			record.ValueByKey("max"),
			record.ValueByKey("avg"),
		)
	}

	// Check for errors during iteration
	if result.Err() != nil {
		log.Fatalf("Query parsing error: %v", result.Err())
	}

	client.Close()
}
