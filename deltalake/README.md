# Deltalake
* It's a collection of immutable parquet files
* Each operation will create a new parquet file, and adds transaction logs to `_delta_logs` directory
* The union of all transactions is the final table
* When something gets updated, the whole parquet file where the data belongs is copied, modified and put into a new file. The transaction log will not say `Remove old file, add newFile`
* The `OPTIMIZE` command merges small parquets into larger ones
* The transaction log is retained only for a set period of time (configurable)


## Spark
* Normally, you would run this in your spark cluster and use delta spark extension.
* Since I already have setup spark cluster myself before, I'm skipping that whole setup and just gonna work with standalone spark node
* Install `pyspark` and `delta-spark`
* Export Java Home to 11 - `export JAVA_HOME="/opt/homebrew/opt/openjdk@11/"`
* Using delta spark version 2.4.0 cuz the jar extension version should also be the same

