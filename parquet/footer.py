import struct

PARQUET_MAGIC = b"PAR1"


def read_parquet_footer(file_path):
    with open(file_path, "rb") as f:
        f.seek(-8, 2)  # Go to last 8 bytes
        footer_bytes = f.read(8)

        footer_len = struct.unpack("<I", footer_bytes[:4])[0]
        magic = footer_bytes[4:]

        if magic != PARQUET_MAGIC:
            raise ValueError("Not a valid Parquet file (missing magic footer)")

        footer_start = f.tell() - 8 - footer_len
        f.seek(footer_start)
        footer_data = f.read(footer_len)

        print(f"Footer length: {footer_len} bytes")
        print(f"Footer starts at byte: {footer_start}")
        print(f"Raw footer bytes (hex): {footer_data[:64].hex()} ...")  # preview

        return footer_data


# Example usage
footer = read_parquet_footer("data.parquet")
