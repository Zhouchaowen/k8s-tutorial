# 应用：部署PostgreSQL

```yaml
# web-statefulset.yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  ports:
    - port: 80
      name: web
  clusterIP: None
  selector:
    app: nginx
---
apiVersion: v1
kind: PersistentVolume        # 类型为PV
metadata:
  name: host-path-pv           # pv的名称
  labels:
    type: local
spec:
  storageClassName: manual
  capacity:
    storage: 1Gi             # pv可用的大小
  accessModes:
    - ReadWriteOnce          # PV以read-write挂载到
  hostPath:
    path: "/data/k8s/test/hostPath"
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web
spec:
  serviceName: "nginx"
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.7.9
          ports:
            - containerPort: 80
              name: web
          volumeMounts:
            - name: www
              mountPath: /usr/share/nginx/html
  volumeClaimTemplates:
    - metadata:
        name: www
      spec:
        accessModes: [ "ReadWriteOnce" ]
        resources:
          requests:
            storage: 1Gi
```

创建

```bash
kubectl apply -f web-statefulset.yaml
```

顺序创建Pod

```bash
kubectl get pods -w -l app=nginx
```







## 参考

https://devopscube.com/deploy-postgresql-statefulset/
https://adamtheautomator.com/postgres-to-kubernetes/
https://medium.com/@suyashmohan/setting-up-postgresql-database-on-kubernetes-24a2a192e962