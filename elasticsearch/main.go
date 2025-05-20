package main

import (
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
	doc := `{"title":"Go in Action","author":"William Kennedy","year":2016}`
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

	// Done
	fmt.Println("\nAll done at", time.Now())
}
