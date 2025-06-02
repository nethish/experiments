package main

import (
	"log"
	"os"
	"text/template"
)

type Product struct {
	Name  string
	Price float64
}

type Store struct {
	Name     string
	Location string
	Products []Product
	Ratings  map[string]float64
}

func main() {
	store := Store{
		Name:     "GoMart",
		Location: "Downtown",
		Products: []Product{
			{Name: "Laptop", Price: 1200.00},
			{Name: "Mouse", Price: 25.50},
			{Name: "Keyboard", Price: 75.00},
		},
		Ratings: map[string]float64{
			"Service": 4.8,
			"Quality": 4.5,
			"Value":   4.2,
		},
	}

	// Template for ranging over a slice of products
	const productListTpl = `
Store: {{.Name}} ({{.Location}})

Our Products:
{{range .Products}} - {{.Name}} (${{printf "%.2f" .Price}})
{{end}}
`
	t1, err := template.New("products").Parse(productListTpl)
	if err != nil {
		log.Fatalf("Error parsing products template: %v", err)
	}
	log.Println("--- Product List ---")
	err = t1.Execute(os.Stdout, store)
	if err != nil {
		log.Fatalf("Error executing products template: %v", err)
	}

	// Template for ranging over a map of ratings
	const ratingsTpl = `
Store Ratings:
{{range $key, $value := .Ratings}} - {{$key}}: {{$value}} stars
{{end}}
`
	t2, err := template.New("ratings").Parse(ratingsTpl)
	if err != nil {
		log.Fatalf("Error parsing ratings template: %v", err)
	}
	log.Println("\n--- Store Ratings ---")
	err = t2.Execute(os.Stdout, store)
	if err != nil {
		log.Fatalf("Error executing ratings template: %v", err)
	}
}

