# KEDA
```bash
helm repo add kedacore https://kedacore.github.io/charts

minikube addons enable metrics-server

kubectl apply -f deployment.yaml
kubectl get pods -l app=nginx
# kubectl apply -f service.yaml
kubectl apply -f scaledobject.yaml

kubectl get hpa

kubectl get scaledobjects
kubectl get hpa keda-hpa-nginx-scaledobject -w


# helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
# helm repo update
# helm install prometheus prometheus-community/kube-prometheus-stack --namespace monitoring --create-namespace
# kubectl get pods -n monitoring
# kubectl get svc -n monitoring prometheus-kube-prometheus-prometheus

# kubectl delete scaledobject http-app-scaledobject -n default
# kubectl apply -f scaledobject.yaml -n default # Make sure scaledobject.yaml has the corrected serverAddress

minikube service http-app-service --url
# http://127.0.0.1:49189

go install github.com/rakyll/hey@latest

hey -z 1m -c 10 http://127.0.0.1:49189 # Send requests for 1 minute and max 10 connections at a time
```

