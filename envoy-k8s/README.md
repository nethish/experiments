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
## Experiment 2
Here the envoy runs as sidecar for both app-a and app-b. The envoy config knows how to reach the other service

```bash
kubectl create configmap envoy-app-a-config --from-file=envoy.yaml=envoy-app-a.yaml
kubectl create configmap envoy-app-b-config --from-file=envoy.yaml=envoy-app-b.yaml


kubectl apply -f app-a.yaml
kubectl apply -f app-b.yaml

kubectl port-forward pod/app-a 15000:15000

curl http://localhost:15000/anything -H "Host: app-b.default.svc.cluster.local"

kubectl exec -it app-a -- curl http://app-b_service/anything

# If you want logs of container running inside a pod
kubectl logs app-a -c envoy
kubectl logs app-b -c envoy



```
