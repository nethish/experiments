curl -X POST "localhost:9200/books/_doc/1" -H 'Content-Type: application/json' -d'
{
  "title": "The Go Programming Language",
  "author": "Alan A. A. Donovan",
  "year": 2015
}'

