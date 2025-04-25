package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Add a delay to ensure replica set is ready
	time.Sleep(2 * time.Second)

	// MongoDB connection string with replica set
	uri := "mongodb://localhost:27017/?replicaSet=myReplicaSet&directConnection=true"

	// Create a new client and connect to the server
	clientOpts := options.Client().ApplyURI(uri).SetServerSelectionTimeout(10 * time.Second)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	// Ping the server to confirm connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	fmt.Println("Successfully connected to MongoDB")

	// Get the database and collection
	database := client.Database("test")
	collection := database.Collection("users")

	// Create a change stream
	opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	changeStream, err := collection.Watch(context.TODO(), mongo.Pipeline{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer changeStream.Close(context.TODO())

	fmt.Println("Listening for MongoDB CDC events...")
	fmt.Println("Try inserting/updating/deleting documents in the 'users' collection")

	// Process change events
	for changeStream.Next(context.TODO()) {
		var changeDoc bson.M
		if err := changeStream.Decode(&changeDoc); err != nil {
			log.Fatal(err)
		}

		// Print the change event
		fmt.Printf("Change Event: %v\n", changeDoc)
	}

	if err := changeStream.Err(); err != nil {
		log.Fatal(err)
	}
}
