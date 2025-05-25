# Deltalake
* Normally, you would run this in your spark cluster and use delta spark extension.
* Since I already have setup spark cluster myself before, I'm skipping that whole setup and just gonna work with standalone spark node
* Install `pyspark` and `delta-spark`
* Export Java Home to 11 - `export JAVA_HOME="/opt/homebrew/opt/openjdk@11/"`
* Using delta spark version 2.4.0 cuz the jar extension version should also be the same
