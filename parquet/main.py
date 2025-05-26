import pyarrow.parquet as pq
import sys


def show_footer(file_path):
    pf = pq.ParquetFile(file_path)

    print(f"Number of Row Groups: {pf.num_row_groups}")
    print(f"Schema:\n{pf.schema}\n")

    metadata = pf.metadata
    print(f"Number of Columns: {metadata.num_columns}")
    print(f"Number of Rows: {metadata.num_rows}")
    print(f"Created by: {metadata.created_by}")
    print()

    for i in range(pf.num_row_groups):
        row_group = metadata.row_group(i)
        print(f"Row Group {i}:")
        print(f"  Number of rows: {row_group.num_rows}")
        print(f"  Total byte size: {row_group.total_byte_size}")
        for j in range(row_group.num_columns):
            column = row_group.column(j)
            print(f"    Column {j}: {column.path_in_schema}")
            print(f"      Type: {column.physical_type}")
            print(f"      Encodings: {column.encodings}")
            print(f"      Compression: {column.compression}")
            print(f"      Data page offset: {column.data_page_offset}")
            print(f"      Dictionary page offset: {column.dictionary_page_offset}")
            print(f"      Total compressed size: {column.total_compressed_size}")

            stats = column.statistics
            if stats is not None:
                print(f"      Stats:")
                print(f"        Null count: {stats.null_count}")
                print(f"        Distinct count: {stats.distinct_count}")
                print(f"        Min: {stats.min}")
                print(f"        Max: {stats.max}")
            else:
                print(f"      Stats: None")

        print()
        break


if __name__ == "__main__":
    show_footer("./data.parquet")
