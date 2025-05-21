package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 1. Index a document
	doc := `{"title":"Go in Action","author":"william kennedy","year":2016}`
	res, err := es.Index("books", strings.NewReader(doc), es.Index.WithDocumentID("1"))
	if err != nil {
		log.Fatalf("Index error: %s", err)
	}
	defer res.Body.Close()
	fmt.Println("Indexed document:", res.Status())

	// 2. Search documents with match query
	query := `{
		"query": {
			"match": { "title": "Go" }
		}
	}`
	res, err = es.Search(
		es.Search.WithIndex("books"),
		es.Search.WithBody(strings.NewReader(query)),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Search error: %s", err)
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing response body: %s", err)
	}
	fmt.Println("Search results:")
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		src := hit.(map[string]interface{})["_source"]
		b, _ := json.MarshalIndent(src, "", "  ")
		fmt.Println(string(b))
	}

	// 3. Aggregation example (books per author)
	//
	// Why author.keyword?
	// In Elasticsearch, a field like "author" that is mapped as a text type is analyzed — which means it's tokenized, lowercased, etc., for full-text search. However, you can’t use text fields in aggregations or exact matches.

	// To support both full-text search and exact matching, Elasticsearch automatically creates a .keyword subfield for text fields.
	aggQuery := `{
		"size": 0,
		"aggs": {
			"books_per_author": {
				"terms": {
					"field": "author.keyword"
				}
			}
		}
	}`
	res, err = es.Search(
		es.Search.WithIndex("books"),
		es.Search.WithBody(strings.NewReader(aggQuery)),
	)
	if err != nil {
		log.Fatalf("Agg error: %s", err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Agg decode error: %s", err)
	}
	buckets := r["aggregations"].(map[string]interface{})["books_per_author"].(map[string]interface{})["buckets"].([]interface{})
	fmt.Println("\nBooks per author:")
	for _, b := range buckets {
		bucket := b.(map[string]interface{})
		fmt.Printf("%s: %v books\n", bucket["key"], bucket["doc_count"])
	}

	// 4. Bulk indexing
	bulk := `
{ "index":{ "_index":"books" }}
{ "title":"Learning Go", "author":"Jon Bodner", "year":2021 }
{ "index":{ "_index":"books" }}
{ "title":"Go Systems Programming", "author":"Mihalis Tsoukalos", "year":2017 }
`
	res, err = es.Bulk(strings.NewReader(bulk))
	if err != nil {
		log.Fatalf("Bulk index error: %s", err)
	}
	fmt.Println("\nBulk indexing done:", res.Status())

	// 5. Range query (books after 2015)
	rangeQuery := `{
		"query": {
			"range": {
				"year": { "gte": 2016 }
			}
		}
	}`
	res, err = es.Search(
		es.Search.WithIndex("books"),
		es.Search.WithBody(strings.NewReader(rangeQuery)),
	)
	if err != nil {
		log.Fatalf("Range search error: %s", err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Range decode error: %s", err)
	}
	fmt.Println("\nBooks published after 2015:")
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		src := hit.(map[string]interface{})["_source"]
		b, _ := json.MarshalIndent(src, "", "  ")
		fmt.Println(string(b))
	}

	// 6. Index some time series documents with @timestamp and response_time
	docs := []string{
		`{"@timestamp":"2025-05-19T10:00:00Z","service":"api","response_time":150}`,
		`{"@timestamp":"2025-05-19T11:00:00Z","service":"api","response_time":250}`,
		`{"@timestamp":"2025-05-20T09:30:00Z","service":"api","response_time":300}`,
		`{"@timestamp":"2025-05-20T14:00:00Z","service":"api","response_time":100}`,
	}

	for i, doc := range docs {
		res, err := es.Index(
			"metrics",
			strings.NewReader(doc),
			es.Index.WithDocumentID(fmt.Sprintf("%d", i+1)),
			es.Index.WithRefresh("true"), // Make it visible for search immediately
		)
		if err != nil {
			log.Fatalf("Index error: %s", err)
		}
		res.Body.Close()
	}

	// 7. Date histogram aggregation: Group by day, get avg response_time
	query = `{
		"size": 0,
		"aggs": {
			"daily_response_time": {
				"date_histogram": {
					"field": "@timestamp",
					"calendar_interval": "day"
				},
				"aggs": {
					"avg_response_time": {
						"avg": { "field": "response_time" }
					}
				}
			}
		}
	}`

	res, err = es.Search(
		es.Search.WithIndex("metrics"),
		es.Search.WithBody(strings.NewReader(query)),
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Fatalf("Search error: %s", err)
	}
	defer res.Body.Close()

	r = make(map[string]any)
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Decode error: %s", err)
	}

	// Parse aggregation buckets
	buckets = r["aggregations"].(map[string]interface{})["daily_response_time"].(map[string]interface{})["buckets"].([]interface{})

	fmt.Println("Daily average response times:")
	for _, b := range buckets {
		bucket := b.(map[string]interface{})
		date := bucket["key_as_string"]
		avgResp := bucket["avg_response_time"].(map[string]interface{})["value"]
		fmt.Printf("%s => %.2f ms\n", date, avgResp)
	}

	// 8 Fuzzy
	fuzzy()

	// Done
	fmt.Println("\nAll done at", time.Now())
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func fuzzy() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	index := "users"

	// 1. Create the index (optional, safe to skip if already created)
	es.Indices.Create(index)

	// 2. Index a document
	users := []User{
		{"Nethish", 30},
		{"Nathish", 28},
		{"Neethush", 31},
		{"Nitish", 27},
	}

	for i, u := range users {
		data, _ := json.Marshal(u)
		_, err := es.Index(index,
			bytes.NewReader(data),
			es.Index.WithDocumentID(fmt.Sprint(i+1)),
			es.Index.WithRefresh("true"),
		)
		if err != nil {
			log.Fatalf("Failed to index doc %d: %v", i, err)
		}
	}

	// 3. Perform a fuzzy match query
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": map[string]interface{}{
					"query":     "Nethih", // Typo intended
					"fuzziness": "AUTO",
				},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Search error: %s", err)
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response: %s", err)
	}

	// 4. Print results
	fmt.Printf("Fuzzy Search Results:\n")
	jsonBytes, _ := json.MarshalIndent(r, "", " ")
	fmt.Println(string(jsonBytes))
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		doc := hit.(map[string]interface{})["_source"]
		j, _ := json.MarshalIndent(doc, "", "  ")
		fmt.Println(string(j))
	}
}
