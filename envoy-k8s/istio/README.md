# Istio
Istio is an open-source service mesh platform that helps you connect, secure, control, and observe microservices in a distributed system — especially in Kubernetes environments.

* Service discovery and routing
* Traffic management 
* Security and Observability

```bash
curl -L https://istio.io/downloadIstio | sh -
cd istio-*/bin
export PATH=$PWD:$PATH

istioctl install --set profile=demo -y

kubectl get pods -n istio-system

# Enable istio injection
kubectl label namespace default istio-injection=enabled

kubectl apply -f app.yaml

kubectl get pods

kubectl apply -f curl.yaml

kubectl exec -it curl -- curl http://httpbin.default.svc.cluster.local:8000/get

# Istio is able to detect services running in the cluster and just be able to route
kubectl create namespace nethish
kubectl apply -f app.yaml -n nethish

kubectl exec -it curl -- curl http://httpbin.nethish.svc.cluster.local:8000/get
```
