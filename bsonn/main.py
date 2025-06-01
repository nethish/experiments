import bson
import datetime

def print_bson_document(python_dict):
    """
    Encodes a Python dictionary into BSON, then decodes it back
    to demonstrate reading, and prints the decoded dictionary.
    """
    print(f"Original Python Dictionary:\n{python_dict}\n")

    # 1. Encode the Python dictionary into BSON bytes
    # This simulates what MongoDB does when storing your data
    bson_bytes = bson.BSON.encode(python_dict)
    print(f"Encoded BSON (raw bytes) with length {len(bson_bytes)}:\n{bson_bytes}\n")

    # 2. Decode the BSON bytes back into a Python dictionary
    # This simulates what MongoDB does when retrieving your data
    decoded_dict = bson.BSON.decode(bson_bytes)
    print(f"Decoded BSON (back to Python dictionary):\n{decoded_dict}\n")

    # Verify that the decoded dictionary is the same as the original
    print(f"Is decoded dictionary equal to original? {decoded_dict == python_dict}")

# --- Example 1: Basic Document ---
print("--- Example 1: Basic Document ---")
document1 = {
    "name": "Alice",
    "age": 30,
    "isStudent": False,
    "courses": ["Math", "Science"],
    "address": {"street": "123 Main St", "city": "Anytown"}
}
print_bson_document(document1)

print("\n" + "="*50 + "\n")

# --- Example 2: Document with BSON-specific types ---
print("--- Example 2: Document with BSON-specific types ---")
document2 = {
    "_id": bson.ObjectId(),  # MongoDB's unique ID type
    "price": bson.Decimal128("19.99"), # High-precision decimal
    "last_updated": datetime.datetime.now(), # Date type
    "binary_data": bson.binary.Binary(b'\x01\x02\x03\x04'), # Binary data
    "null_value": None,
    "regex_pattern": bson.regex.Regex('^abc', 'i') # Regular expression
}
print_bson_document(document2)

print("\n" + "="*50 + "\n")

# --- Example 3: Empty Document ---
print("--- Example 3: Empty Document ---")
document3 = {}
print_bson_document(document3)