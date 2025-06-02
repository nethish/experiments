# Helm
```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm search repo bitnami
helm repo update
helm install bitnami/mysql --generate-name
helm list
helm uninstall mysql-1612624192
```


```bash
helm create chaat
helm install --dry-run --debug my-release chaat
```
