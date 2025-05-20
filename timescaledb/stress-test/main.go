package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const (
	dbConnStr     = "host=localhost port=5432 user=postgres password=secret dbname=metrics sslmode=disable"
	concurrency   = 20    // number of concurrent workers
	insertCount   = 10000 // total inserts per worker
	queryInterval = 5 * time.Second
)

func main() {
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal("DB open error:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("DB ping error:", err)
	}

	fmt.Println("Starting stress test with", concurrency, "workers")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Launch workers to insert data concurrently
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go insertWorker(ctx, db, i, &wg)
	}

	// Periodic query goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(queryInterval):
				queryRecent(db)
			}
		}
	}()

	wg.Wait()
	fmt.Println("Stress test completed")
}

func insertWorker(ctx context.Context, db *sql.DB, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	location := fmt.Sprintf("worker-%d", id)
	stmt, err := db.PrepareContext(ctx, "INSERT INTO temperature (time, location, value) VALUES ($1, $2, $3)")
	if err != nil {
		log.Printf("[worker %d] Prepare error: %v\n", id, err)
		return
	}
	defer stmt.Close()

	start := time.Now()
	for i := 0; i < insertCount; i++ {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, err := stmt.ExecContext(ctx, time.Now(), location, 20+10*rand.Float64())
		if err != nil {
			log.Printf("[worker %d] Insert error at iteration %d: %v\n", id, i, err)
		}
	}
	elapsed := time.Since(start)
	log.Printf("[worker %d] Inserted %d rows in %v (%.2f inserts/sec)\n", id, insertCount, elapsed, float64(insertCount)/elapsed.Seconds())
}

func queryRecent(db *sql.DB) {
	start := time.Now()
	rows, err := db.Query("SELECT time, location, value FROM temperature ORDER BY time DESC LIMIT 5")
	if err != nil {
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	log.Println("Recent temperature entries:")
	for rows.Next() {
		var ts time.Time
		var loc string
		var val float64
		if err := rows.Scan(&ts, &loc, &val); err != nil {
			log.Println("Scan error:", err)
			return
		}
		log.Printf("%s | %-12s | %.2f°C\n", ts.Format(time.RFC3339), loc, val)
	}
	log.Printf("Query took %v\n", time.Since(start))
}
