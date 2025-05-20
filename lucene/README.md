# Lucene
* From what I see it's just a library that parses text, tokenizes it and stores data to disk.
* Elastic Search is built on top of Lucene which adds distributed compute, aggregations and backup features
* It primarily uses Inverted Index to quickly retrieve documents matching a term.
* But it uses bunch of other algos to retrieve the documents as well (TODO)

# TODO
* Create a standalone Java file that uses lucene to store and retrieve data
* There is Lucene Luke to analyze the stored indexes in lucene
* These Java apps are not easy to just write and test :sad:
* See what other algos does lucene engine uses to retrieve docs
