from fastavro import reader
import time

def parse_avro(path):
    count = 0
    with open(path, "rb") as f:
        for record in reader(f):
            count += 1
    return count

if __name__ == "__main__":
    start = time.time()
    n = parse_avro("data/data.avro")
    duration = time.time() - start
    print(f"Parsed {n} Avro records in {duration:.2f} seconds")
