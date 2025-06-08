# CRD

This is a simple mysql CRD. When you create, update or delete the object, the callbacks gets called.
The controller runs locally, and listens to minikube. For production grade, you'll simply deploy this in your cluster, setup service account with necessary roles (permissions) so that your controller can do kubectl operations.

```bash
kubectl apply -f mysqlinstance-crd.yaml

kubectl apply -f mysqlinstance-crd.yaml
go run main.go

kubectl apply -f mysqlinstance.yaml
```
