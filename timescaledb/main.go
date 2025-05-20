package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s not set", key)
	}
	return v
}

func connect() *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost",
		"5432",
		"postgres",
		"secret",
		"metrics",
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}
	return db
}

func insertMultiple(db *sql.DB) {
	fmt.Println("1) Inserting multiple rows...")
	rows := []struct {
		location string
		value    float64
	}{
		{"office", 22.3},
		{"lab", 25.1},
		{"server-room", 27.0},
	}
	for _, r := range rows {
		_, err := db.Exec(
			"INSERT INTO temperature (time, location, value) VALUES ($1, $2, $3)",
			time.Now(), r.location, r.value,
		)
		if err != nil {
			log.Fatalf("insert failed: %v", err)
		}
	}
	fmt.Println("  → inserted batch successfully")
}

func queryRecent(db *sql.DB) {
	fmt.Println("2) Querying recent rows (last 1 minute)...")
	rows, err := db.Query(
		`SELECT time, location, value 
		 FROM temperature 
		 WHERE time > NOW() - INTERVAL '1 minute' 
		 ORDER BY time DESC`,
	)
	if err != nil {
		log.Fatalf("query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t time.Time
		var loc string
		var val float64
		if err := rows.Scan(&t, &loc, &val); err != nil {
			log.Fatalf("scan failed: %v", err)
		}
		fmt.Printf("  → %s | %-12s | %.2f\n", t.Format(time.RFC3339), loc, val)
	}
}

func queryAggregation(db *sql.DB) {
	fmt.Println("3) Aggregation: average value per location in last hour")
	rows, err := db.Query(
		`SELECT location, AVG(value) AS avg_val
		 FROM temperature
		 WHERE time > NOW() - INTERVAL '1 hour'
		 GROUP BY location`,
	)
	if err != nil {
		log.Fatalf("aggregation query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var loc string
		var avg float64
		if err := rows.Scan(&loc, &avg); err != nil {
			log.Fatalf("scan failed: %v", err)
		}
		fmt.Printf("  → %-12s | avg=%.2f\n", loc, avg)
	}
}

func deleteOld(db *sql.DB) {
	fmt.Println("4) Deleting rows older than 24 hours...")
	res, err := db.Exec(
		"DELETE FROM temperature WHERE time < NOW() - INTERVAL '24 hours'",
	)
	if err != nil {
		log.Fatalf("delete failed: %v", err)
	}
	affected, _ := res.RowsAffected()
	fmt.Printf("  → %d rows deleted\n", affected)
}

func preparedStatement(db *sql.DB) {
	fmt.Println("5) Using prepared statement for inserts...")
	stmt, err := db.Prepare(
		"INSERT INTO temperature (time, location, value) VALUES ($1, $2, $3)",
	)
	if err != nil {
		log.Fatalf("prepare failed: %v", err)
	}
	defer stmt.Close()

	for i := 0; i < 3; i++ {
		_, err := stmt.Exec(time.Now().Add(time.Duration(i)*time.Minute), "sensor-X", 20.0+float64(i))
		if err != nil {
			log.Fatalf("stmt exec failed: %v", err)
		}
	}
	fmt.Println("  → prepared-statement batch inserted")
}

func main() {
	db := connect()
	defer db.Close()

	insertMultiple(db)
	queryRecent(db)
	queryAggregation(db)
	deleteOld(db)
	preparedStatement(db)

	fmt.Println("All operations completed.")
}
