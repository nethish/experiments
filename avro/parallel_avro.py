from fastavro import block_reader
from multiprocessing import Pool
from pathlib import Path

AVRO_FILE = Path("data/data.avro")
NUM_WORKERS = 4  # Adjust based on your CPU


def process_block(block_info):
    block_index, records = block_info
    # For example, count records with age > 60
    count = sum(1 for r in records if r["age"] > 60)
    return count


def split_avro_and_process(path, num_workers):
    block_infos = []
    with open(path, "rb") as f:
        for i, block in enumerate(block_reader(f)):
            block_infos.append((i, block))

    print(f"Number of blocks to pprocess: {len(block_infos)}")
    with Pool(num_workers) as pool:
        results = pool.map(process_block, block_infos)

    print(f"Total records with age > 60: {sum(results)}")


if __name__ == "__main__":
    split_avro_and_process(AVRO_FILE, NUM_WORKERS)
