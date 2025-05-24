from pyspark import SparkConf, SparkContext

# Set up configuration
conf = SparkConf().setAppName("SumListApp")
sc = SparkContext(conf=conf)

# Create an RDD with numbers from 1 to 1,000,000 split into 4 partitions
rdd = sc.parallelize(range(1, 1_000_001), numSlices=4)

# Print the number of partitions
print("Number of partitions:", rdd.getNumPartitions())

# Compute the sum
total = rdd.reduce(lambda a, b: a + b)

# Show the result
print("Total sum is:", total)

# Wait a bit so you can check the Spark UI
input("Press Enter to stop SparkContext...")

sc.stop()
