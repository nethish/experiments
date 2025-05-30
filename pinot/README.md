# Pinot
* It's a realtime OLAP database that has first class integration to consume real time events
* It differs from other OLAP databases in it's indexing approach and integrations with real time data
* It's near instant actually, and I'm impressed
* Uses zookeeper for coordination

## Components
Many components to support high throughput
### Controller
* Cluster manager
* Table creation, segment management
* Web UI
### Broker
* Query Router
* Plans the query
* Aggregates responses
### Worker
* Data processor/ Executor
* Stores data segments
### Minion
* Runs asynchronous background jobs
* Segment merging and data cleanup

## How to define table schema
### You upload a `schema.json` doc
* dimensionFieldSpecs
  * Filtering
  * Grouping
  * These are no aggregatable
* metricFieldSpecs
  * Aggregations
* timeFieldSpec
  * Time series based

## Segments
* A segment is the fundamental unit of data storage and query in Pinot.
* It’s essentially a self-contained columnar data file that stores a chunk of your ingested data — along with its indexes — in a format optimized for fast analytical queries.

