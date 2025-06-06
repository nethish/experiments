# Envoy

```bash
kubectl create configmap envoy-config --from-file=envoy.yaml=envoy-config.yaml

kubectl expose pod envoy --port=10000 --type=NodePort

minikube service envoy --url

curl http://127.0.0.1:54839

kubectl port-forward pod/envoy 9901:9901

```


## TODO
Understand reverse proxy first
How do you solve the problem without envoy - Alternatives and Kube native tool

Envoy internals and it's features
How to make it run as side car
