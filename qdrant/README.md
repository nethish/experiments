# qdrant - Vector database
It is a vector database


## Embeddings
Let's say the dimension of each token is 5. Then for a "hi there", there are two tokens and each will give you a vector. 
Then pooling is applied to get the final result. One example is mean pooling. 

```
"hi there" = [[1, 2, 3, 4, 5], [6, 7, 8, 9, 10]]
Then pooling = [[mean(1, 6), mean(2, 7) ...]] -- Reduced to one vector
```
