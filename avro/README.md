# Avro Format
* Avro files store the schema along with the data, so you don't have to transfer the schema
* Avro is binary so the resulting file is smaller than JSON
* Mainly used in Hadoop ecosystem where I have no knowledge about.
* Surprisingly, avro deser is slower than json when writing this README. 
  * 10M records in JSON took 11 seconds while avro took 14 seconds
* Parallel processing by horizontal sharding of the avro file
