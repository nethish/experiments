import pandas as pd
import numpy as np


def generate_large_data(
    num_rows=100_000_000, csv_path="data.csv", parquet_path="data.parquet"
):
    # Generate synthetic data
    df = pd.DataFrame(
        {
            "id": np.arange(num_rows),
            "value": np.random.randn(num_rows),
            "category": np.random.choice(["A", "B", "C", "D"], size=num_rows),
        }
    )

    # Write to CSV
    print(f"Writing {num_rows} rows to CSV: {csv_path}")
    df.to_csv(csv_path, index=False)

    # Write to Parquet
    print(f"Writing {num_rows} rows to Parquet: {parquet_path}")
    df.to_parquet(parquet_path, index=False)


if __name__ == "__main__":
    generate_large_data()
