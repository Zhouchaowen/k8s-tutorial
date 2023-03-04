# word-press 实战


## 部署
```bash
kubectl apply -f wp-namespace.yaml
kubectl apply -f wp-deployment.yaml
kubectl apply -f wp-mysql.yaml
kubectl apply -f wp-service.yaml
```

## hpa
```bash
kubectl autoscale deployment wordpress --namespace kube-example --cpu-percent=20 --min=3 --max=6
```