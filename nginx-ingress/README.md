# NGINX Ingress

Ingress exposes HTTP and HTTPS routes from outside the cluster to services within the cluster. Traffic routing is controlled by rules defined on the Ingress resource.

It supports host and path based routing.



Follow `https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/`


## Prerequisites
* Install minikube and `minikube start`

```bash
minikube addons enable ingress
kubectl get pods -n ingress-nginx # Can take upto a miniute

# Deploy an app
kubectl create deployment web --image=gcr.io/google-samples/hello-app:1.0
kubectl get deployment web 

# This creates a Service of type NodePort exposed at 8080
kubectl expose deployment web --type=NodePort --port=8080
kubectl get service web

minikube service web --url 
curl http://127.0.0.1:62445 


# Create an ingress - ingress.yaml
kubectl apply -f https://k8s.io/examples/service/networking/example-ingress.yaml

kubectl describe ingress example-ingress

kubectl get ingress

minikube tunnel
curl --resolve "hello-world.example:80:127.0.0.1" -i http://hello-world.example

# Optionally if you want to access from browser, add the below line to /etc/hosts
127.0.0.1 hello-world.example

# See your logs here
kubectl logs ingress-nginx-controller-56d7c84fd4-9c6jb -n ingress-nginx | tail -n 10

kubectl edit deploy ingress-nginx-controller -n ingress-nginx
```

Create a second deployment

```bash
kubectl create deployment web2 --image=gcr.io/google-samples/hello-app:2.0
kubectl expose deployment web2 --port=8080 --type=NodePort

# You will add another path v2 which is already added
kubectl apply -f ingress.yaml

minikube tunnel

# This will resolve to first web
curl --resolve "hello-world.example:80:127.0.0.1" -i http://hello-world.example

# This will resolve to web2
curl --resolve "hello-world.example:80:127.0.0.1" -i http://hello-world.example/v2

```

