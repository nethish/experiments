from fastavro import writer, reader, parse_schema

schema = {
    "type": "record",
    "name": "User",
    "fields": [
        {"name": "name", "type": "string"},
        {"name": "age", "type": "int"},
        {"name": "email", "type": ["null", "string"], "default": None}
    ]
}

records = [
    {"name": "Alice", "age": 30, "email": "alice@example.com"},
    {"name": "Bob", "age": 25, "email": None}
]

with open("users.avro", "wb") as out:
    writer(out, parse_schema(schema), records)
