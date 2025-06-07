# Istio
Istio is an open-source service mesh platform that helps you connect, secure, control, and observe microservices in a distributed system — especially in Kubernetes environments.

* Service discovery and routing
* Traffic management 
* Security and Observability

So, you install istio containers inside your cluster. Label your namepsace with injection enabled.
Start deploying your cluster

## How does istio intercepts traffic? 
When the pod starts:
* Istio runs a script inside the sidecar container (istio-init or istio-iptables.sh)
* That script sets up iptables rules inside the pod’s network namespace
* The iptables rules do:
* Redirect outbound traffic to port 15001 (Envoy outbound listener)
* Redirect inbound traffic to port 15006 (Envoy inbound listener)


```bash
curl -L https://istio.io/downloadIstio | sh -
cd istio-*/bin
export PATH=$PWD:$PATH

istioctl install --set profile=demo -y

kubectl get pods -n istio-system

# Enable istio injection. This will inject the sidecar when you create a deployment in the default namespace
kubectl label namespace default istio-injection=enabled

# This creates the pod
kubectl apply -f app.yaml

# You should 2/2 - both httpbin and istio running
kubectl get pods

# Start another container
kubectl apply -f curl.yaml

# Try to connect to httpbin from curl
kubectl exec -it curl -- curl http://httpbin.default.svc.cluster.local:8000/get

# Istio is able to detect services running in the cluster and just be able to route
kubectl create namespace nethish
kubectl apply -f app.yaml -n nethish

kubectl exec -it curl -- curl http://httpbin.nethish.svc.cluster.local:8000/get
```

## TODO
How does istio sets up iptables for redirection
