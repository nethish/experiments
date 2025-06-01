# BSON Format
Simply a binary encoding rules to store a JSON file.
For example take
```JSON
{"user_id": 1, "data": {"a":  "b"}}
```

This document is encoded to something like
```BSON
4 bytes to denote length | Data type of value | Key Length | Key | Value | Datatype of Value | Key Length | Key | Value
```

You see, there is no need to encode chars like `{`, `:`, `,` so on. Simple binary info of what's needed for less memory usage.