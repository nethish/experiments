# Envoy

```bash
curl -i http://localhost:10000/get

for i in {1..10}; do curl -s -o /dev/null -w "%{http_code}\n" http://localhost:10000/get; done

```
