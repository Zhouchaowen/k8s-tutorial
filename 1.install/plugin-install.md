## dashboard

安装

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml
```

获取Token

```bash
kubectl get secret -n kubernetes-dashboard | grep token | grep admin

kubectl describe secret admin-user-token-sz4vc -n kubernetes-dashboard
```

需要权限

```bash

```

