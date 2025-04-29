package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoURI        = "mongodb://localhost:27017"
	dbName          = "test_db"
	collName        = "test_collection"
	readerCount     = 200
	writerCount     = 1
	runDuration     = 60 * time.Second
	delayBetweenOps = 10 * time.Millisecond
)

func seedData(coll *mongo.Collection) {
	// coll.Drop(context.Background())
	var docs []interface{}

	// Insert around 50 million data into the collection with index
	for i := 0; i < 1000000; i++ {
		docs = append(docs, bson.M{"x": rand.Int()})
	}
	for range 50 {
		_, err := coll.InsertMany(context.Background(), docs)
		if err != nil {
			log.Fatalf("Failed to seed data: %v", err)
		}
	}
}

func readWorker(coll *mongo.Collection, wg *sync.WaitGroup, stopChan <-chan struct{}) {
	defer wg.Done()
	for {
		select {
		case <-stopChan:
			return
		default:
			start := time.Now()
			ctx, _ := context.WithTimeout(context.TODO(), 20*time.Second)
			cursor, err := coll.Find(ctx, bson.M{"x": bson.M{"$eq": 0}}, options.Find().SetLimit(1))
			fmt.Println(time.Since(start), err)
			if err == nil {
				_ = cursor.All(context.Background(), &[]bson.M{})
				cursor.Close(context.Background())
			}
			// time.Sleep(delayBetweenOps)
		}
	}
}

func writeWorker(coll *mongo.Collection, wg *sync.WaitGroup, stopChan <-chan struct{}) {
	defer wg.Done()
	for {
		select {
		case <-stopChan:
			return
		default:
			_, _ = coll.InsertOne(context.Background(), bson.M{"x": time.Now().UnixNano()})
			// time.Sleep(delayBetweenOps)
		}
	}
}

func main() {
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}
	defer client.Disconnect(ctx)

	coll := client.Database(dbName).Collection(collName)
	// seedData(coll)

	var wg sync.WaitGroup
	stopChan := make(chan struct{})

	fmt.Printf("Starting %d readers and %d writers for %s...\n", readerCount, writerCount, runDuration)

	// Start reader goroutines
	for i := 0; i < readerCount; i++ {
		wg.Add(1)
		go readWorker(coll, &wg, stopChan)
	}

	// Start writer goroutines
	for i := 0; i < writerCount; i++ {
		wg.Add(1)
		go writeWorker(coll, &wg, stopChan)
	}

	time.Sleep(runDuration)
	close(stopChan)
	wg.Wait()
	fmt.Println("Simulation complete.")
}
