import json
import random
from fastavro import writer, parse_schema
from pathlib import Path
import time

NUM_RECORDS = 10_000_000  # 1 million records
DATA_DIR = Path("data")
DATA_DIR.mkdir(exist_ok=True)

schema = {
    "type": "record",
    "name": "User",
    "fields": [
        {"name": "id", "type": "int"},
        {"name": "name", "type": "string"},
        {"name": "age", "type": "int"},
        {"name": "email", "type": ["null", "string"], "default": None},
    ],
}


def generate_record(i):
    return {
        "id": i,
        "name": f"user_{i}",
        "age": random.randint(18, 1000),
        "email": f"user_{i}@example.com" if i % 2 == 0 else None,
    }


def generate_json():
    with open(DATA_DIR / "data.json", "w") as f:
        for i in range(NUM_RECORDS):
            json.dump(generate_record(i), f)
            f.write("\n")  # One JSON per line (NDJSON)


def generate_avro():
    records = (generate_record(i) for i in range(NUM_RECORDS))
    with open(DATA_DIR / "data.avro", "wb") as out:
        writer(out, parse_schema(schema), records)


if __name__ == "__main__":
    print("Generating JSON and Avro files...")
    generate_json()
    generate_avro()
    print("Done.")
