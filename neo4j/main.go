package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type User struct {
	Name   string
	Phone  string
	CardNo string
	IP     string
}

func main() {
	// Neo4j connection configuration
	// Replace with your own credentials if different
	uri := "neo4j://localhost:7687"
	username := "neo4j"
	password := "password"

	// Initialize the driver
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatal("Error creating Neo4j driver:", err)
	}
	defer driver.Close()

	// Verify connectivity
	err = driver.VerifyConnectivity()
	if err != nil {
		log.Fatal("Error connecting to Neo4j:", err)
	}
	fmt.Println("Successfully connected to Neo4j!")

	// Initialize session
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close()

	// Clear existing data for clean runs
	clearDatabase(session)

	// Generate random users
	users := generateRandomUsers(20)

	// Create user nodes
	userIds := createUserNodes(session, users)
	fmt.Printf("Created %d user nodes\n", len(userIds))

	// Create random relationships between users
	createRandomRelationships(session, userIds)
	fmt.Println("Created random relationships between users")

	// Query first-degree relationships
	fmt.Println("\n--- First-Degree Relationships ---")
	queryFirstDegreeRelationships(session, userIds[0])

	// Query second-degree relationships
	fmt.Println("\n--- Second-Degree Relationships ---")
	querySecondDegreeRelationships(session, userIds[0])
}

func clearDatabase(session neo4j.Session) {
	_, err := session.Run("MATCH (n) DETACH DELETE n", map[string]interface{}{})
	if err != nil {
		log.Printf("Warning: Error clearing database: %v", err)
	}
}

func generateRandomUsers(count int) []User {
	rand.Seed(time.Now().UnixNano())
	users := make([]User, count)

	for i := 0; i < count; i++ {
		users[i] = User{
			Name:   fmt.Sprintf("User%d", i+1),
			Phone:  generateRandomPhone(),
			CardNo: generateRandomCardNo(),
			IP:     generateRandomIP(),
		}
	}
	return users
}

func generateRandomPhone() string {
	return fmt.Sprintf("+1-%d%d%d-%d%d%d-%d%d%d%d",
		rand.Intn(10), rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10))
}

func generateRandomCardNo() string {
	return fmt.Sprintf("%d%d%d%d-%d%d%d%d-%d%d%d%d-%d%d%d%d",
		rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10))
}

func generateRandomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func createUserNodes(session neo4j.Session, users []User) []int64 {
	var userIds []int64

	for _, user := range users {
		result, err := session.Run(
			"CREATE (u:User {name: $name, phone: $phone, cardNo: $cardNo, ip: $ip}) RETURN id(u)",
			map[string]interface{}{
				"name":   user.Name,
				"phone":  user.Phone,
				"cardNo": user.CardNo,
				"ip":     user.IP,
			})
		if err != nil {
			log.Printf("Error creating user node: %v", err)
			continue
		}

		if result.Next() {
			// jsonBytes, _ := json.MarshalIndent(result.Record(), "", "  ")
			// fmt.Println(string(jsonBytes))
			userIds = append(userIds, result.Record().Values[0].(int64))
		}
	}
	return userIds
}

func createRandomRelationships(session neo4j.Session, userIds []int64) {
	relationshipTypes := []string{"KNOWS", "WORKS_WITH", "FRIENDS_WITH"}

	// Create about 40 random relationships
	for i := 0; i < 40; i++ {
		// Choose random source and target users
		sourceIdx := rand.Intn(len(userIds))
		targetIdx := rand.Intn(len(userIds))

		// Avoid self-relationships
		if sourceIdx == targetIdx {
			continue
		}

		// Choose a random relationship type
		relType := relationshipTypes[rand.Intn(len(relationshipTypes))]

		// Create the relationship with a random weight property
		_, err := session.Run(
			fmt.Sprintf("MATCH (a:User), (b:User) WHERE id(a) = $sourceId AND id(b) = $targetId "+
				"CREATE (a)-[r:%s {weight: $weight}]->(b)", relType),
			map[string]interface{}{
				"sourceId": userIds[sourceIdx],
				"targetId": userIds[targetIdx],
				"weight":   rand.Intn(10) + 1,
			})
		if err != nil {
			log.Printf("Error creating relationship: %v", err)
		}
	}
}

func queryFirstDegreeRelationships(session neo4j.Session, userId int64) {
	result, err := session.Run(
		"MATCH (u:User)-[r]->(friend:User) WHERE id(u) = $userId "+
			"RETURN type(r) as relationship, friend.name as friendName, "+
			"friend.phone as friendPhone, friend.cardNo as friendCardNo, "+
			"friend.ip as friendIP",
		map[string]interface{}{
			"userId": userId,
		})
	if err != nil {
		log.Printf("Error querying first-degree relationships: %v", err)
		return
	}

	fmt.Printf("First-degree connections for user with ID %d:\n", userId)
	count := 0
	for result.Next() {
		record := result.Record()
		relType := record.Values[0].(string)
		friendName := record.Values[1].(string)
		friendPhone := record.Values[2].(string)
		friendCardNo := record.Values[3].(string)
		friendIP := record.Values[4].(string)

		fmt.Printf("  Relationship: %s to %s (Phone: %s, Card: %s, IP: %s)\n",
			relType, friendName, friendPhone, friendCardNo, friendIP)
		count++
	}

	if count == 0 {
		fmt.Println("  No first-degree connections found")
	}
}

func querySecondDegreeRelationships(session neo4j.Session, userId int64) {
	result, err := session.Run(
		"MATCH (u:User)-[r1]->(friend:User)-[r2]->(friendOfFriend:User) "+
			"WHERE id(u) = $userId AND id(friendOfFriend) <> $userId "+
			"RETURN DISTINCT type(r1) as firstRel, friend.name as friendName, "+
			"type(r2) as secondRel, friendOfFriend.name as fofName, "+
			"friendOfFriend.phone as fofPhone, friendOfFriend.cardNo as fofCardNo, "+
			"friendOfFriend.ip as fofIP",
		map[string]interface{}{
			"userId": userId,
		})
	if err != nil {
		log.Printf("Error querying second-degree relationships: %v", err)
		return
	}

	fmt.Printf("Second-degree connections for user with ID %d:\n", userId)
	count := 0
	for result.Next() {
		record := result.Record()
		firstRel := record.Values[0].(string)
		friendName := record.Values[1].(string)
		secondRel := record.Values[2].(string)
		fofName := record.Values[3].(string)
		fofPhone := record.Values[4].(string)
		fofCardNo := record.Values[5].(string)
		fofIP := record.Values[6].(string)

		fmt.Printf("  Path: -%s-> %s -%s-> %s (Phone: %s, Card: %s, IP: %s)\n",
			firstRel, friendName, secondRel, fofName, fofPhone, fofCardNo, fofIP)
		count++
	}

	if count == 0 {
		fmt.Println("  No second-degree connections found")
	}
}
