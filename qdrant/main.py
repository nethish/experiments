from qdrant_client import QdrantClient
from qdrant_client.models import Distance, VectorParams, PointStruct
from sentence_transformers import SentenceTransformer

# 1. Start client
client = QdrantClient(host="localhost", port=6333)

# 2. Create collection with vector size 384 (MiniLM model)
DIM = 384
COLLECTION_NAME = "my_collection"

client.recreate_collection(
    collection_name=COLLECTION_NAME,
    vectors_config=VectorParams(size=DIM, distance=Distance.COSINE),
)

# 3. Create embeddings
model = SentenceTransformer("all-MiniLM-L6-v2")
sentences = [
    "I like learning",
    "Machine learning is the best course",
    "Deep learning is the best course",
    "I love deep learning",
    "Neural networks are powerful",
    "Bananas are yellow",
]

embeddings = model.encode(sentences).tolist()

# 4. Upload vectors
points = [
    PointStruct(id=i, vector=embeddings[i], payload={"text": sentences[i]})
    for i in range(len(sentences))
]
client.upsert(collection_name=COLLECTION_NAME, points=points)

# 5. Query a similar sentence
query = "I like machine learning"
query_vector = model.encode(query).tolist()

hits = client.search(
    collection_name=COLLECTION_NAME,
    query_vector=query_vector,
    limit=2,
)

# 6. Show results
print("🔍 Query:", query)
for hit in hits:
    print(f"Score: {hit.score:.3f} → {hit.payload['text']}")

