# Kind - Kubernetes In Docker commands

```bash
kind create cluster # Default cluster context name is `kind`.
kind create cluster --name kind-2

kind get clusters

kubectl cluster-info --context kind-kind
kubectl cluster-info --context kind-kind-2

# Docker images can be loaded into kind
kind load docker-image my-custom-image-0 my-custom-image-1

kind get nodes

# Get the list of images present in cluster. Use the output from above command
docker exec -it kind-control-plane crictl images

# Multi node cluster
kind create cluster --config multi-node-cluster.yaml






```
