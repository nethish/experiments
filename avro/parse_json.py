import json
import time

def parse_json(path):
    count = 0
    with open(path, "r") as f:
        for line in f:
            json.loads(line)
            count += 1
    return count

if __name__ == "__main__":
    start = time.time()
    n = parse_json("data/data.json")
    duration = time.time() - start
    print(f"Parsed {n} JSON records in {duration:.2f} seconds")
