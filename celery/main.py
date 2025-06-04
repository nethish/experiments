from tasks import add

result = add.delay(4, 6)
print(result.get(timeout=10))  # Output: 10
