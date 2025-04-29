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
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// Mention the primary node
	uri := "mongodb://localhost:27017/?replicaSet=myReplicaSet&directConnection=true"
	// uri := "mongodb://localhost:27017,localhost:27018,localhost:27019/?replicaSet=myReplicaSet"
	clientOpts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(ctx)

	// Check connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Could not connect to Replica Set:", err)
	}
	fmt.Println("Connected to Replica Set successfully.")

	userColl := client.Database("mydb").Collection("users")
	historyColl := client.Database("mydb").Collection("user_history")

	session, err := client.StartSession()
	if err != nil {
		log.Fatal("Error starting session:", err)
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// Insert into users collection
		_, err := userColl.InsertOne(sc, bson.M{
			"name":  "John Doe",
			"email": "john@example.com",
		})
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed inserting into users: %w", err)
		}

		fmt.Println("Inserted user collection")
		time.Sleep(10 * time.Second)
		// Insert into user_history collection
		_, err = historyColl.InsertOne(sc, bson.M{
			"user_id":   "123",
			"change":    "Updated profile",
			"timestamp": time.Now(),
		})
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed inserting into user_history: %w", err)
		}

		fmt.Println("Inserted history collection")
		if err := session.CommitTransaction(sc); err != nil {
			return fmt.Errorf("commit failed: %w", err)
		}

		fmt.Println("Committed transaction")
		return nil
	})

	if err != nil {
		log.Println("Transaction failed and rolled back:", err)
	} else {
		log.Println("Transaction committed successfully!")
	}
}
