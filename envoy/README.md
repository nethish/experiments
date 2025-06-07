# Envoy
Envoy is a cloud native reverse proxy for microservices.

The configuration looks something like below

bootstrap -> listeners -> filters -> routes -> clusters -> endpoints

* Listeners are entry points -- ports envoy listens on for incoming traffic
* Filters processes the requests before it's forwarded
* Clusters define where Envoy proxies requests to (i.e., backend services).
* Route config is where requests are matched and proxied


```bash
kubectl delete configmap envoy-config
kubectl create configmap envoy-config --from-file=envoy.yaml=envoy-config.yaml

kubectl delete pod envoy
kubectl apply -f envoy.yaml

# Creates a service
kubectl expose pod envoy --port=10000 --type=NodePort

minikube service envoy --url

# After 5 retries it gets rate limited
curl http://127.0.0.1:54839

kubectl port-forward pod/envoy 9901:9901

```

