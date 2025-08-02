from pymongo import MongoClient

# Connect to the mongos instance
client = MongoClient('mongodb://localhost:27020')

# List databases
print("Databases:")
print(client.list_database_names())

# Example: Insert data into the `test` database and `example` collection
db = client.test
collection = db.example

# Insert a document
collection.insert_one({"name": "MongoDB", "type": "database", "count": 1})

# Find the inserted documents
documents = collection.find()
for document in documents:
    print(document)
