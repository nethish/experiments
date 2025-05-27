from hdfs import InsecureClient
from concurrent.futures import ThreadPoolExecutor

# Configuration
HDFS_URL = "http://namenode:9870"  # Replace with your actual namenode URL
HDFS_FILE_PATH = "/user/demo/data.parquet"
BLOCK_SIZE = 128 * 1024 * 1024  # 128 MB block size
NUM_THREADS = 4

# Connect to HDFS
client = InsecureClient(HDFS_URL, user="hdfs")  # replace with actual user if needed

# Get file size
status = client.status(HDFS_FILE_PATH)
print("=========== Status ===========\n", status)
file_size = status["length"]


# Function to read a block
def read_block(start, length):
    with client.read(HDFS_FILE_PATH, offset=start, length=length) as reader:
        data = reader.read()
        print(f"Read {len(data)} bytes from offset {start}")
        return data


# Create tasks
blocks = []
for i in range(0, file_size, BLOCK_SIZE):
    block_len = min(BLOCK_SIZE, file_size - i)
    blocks.append((i, block_len))

# Read in parallel
results = []
with ThreadPoolExecutor(max_workers=NUM_THREADS) as executor:
    futures = [executor.submit(read_block, start, length) for start, length in blocks]
    for future in futures:
        results.append(future.result())

# `results` now holds the data from each block
combined_data = "".join(results)
print("Combined data length:", len(combined_data))
