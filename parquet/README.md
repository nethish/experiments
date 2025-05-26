# Parquet (Incomplete)
* Parquet stores metadata about row group in it's footer
* The last 8 bytes tells where the footer starts
* Distributed processing can leverage this metadata and request only the data needed


## Example metadata
```
Number of Row Groups: 96
Schema:
<pyarrow._parquet.ParquetSchema object at 0x1083d4b00>
required group field_id=-1 schema {
  optional int64 field_id=-1 id;
  optional double field_id=-1 value;
  optional binary field_id=-1 category (String);
}


Number of Columns: 3
Number of Rows: 100000000
Created by: parquet-cpp-arrow version 20.0.0

Row Group 0:
  Number of rows: 1048576
  Total byte size: 17600366
    Column 0: id
      Type: INT64
      Encodings: ('PLAIN', 'RLE', 'RLE_DICTIONARY')
      Compression: SNAPPY
      Data page offset: 524601
      Total compressed size: 4475929
    Column 1: value
      Type: DOUBLE
      Encodings: ('PLAIN', 'RLE', 'RLE_DICTIONARY')
      Compression: SNAPPY
      Data page offset: 5524582
      Total compressed size: 8668467
    Column 2: category
      Type: BYTE_ARRAY
      Encodings: ('PLAIN', 'RLE', 'RLE_DICTIONARY')
      Compression: SNAPPY
      Data page offset: 13144436
      Total compressed size: 264326
```
