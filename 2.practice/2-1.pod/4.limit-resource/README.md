# resources

```bash
kubectl explain pod.spec.containers.resources
```

## requests

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: resource-limits-demo
spec:
  containers:
    - name: resource-limits-demo
      image: nginx:alpine
      ports:
        - containerPort: 80
      resources:
        requests:         # requests 集群调度使用的资源
          memory: 50Mi    # 定义容器或Pod对内存资源的需求量
          cpu: 50m        # 定义容器或Pod对CPU资源的需求量
```

## limits

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: resource-limits-demo
spec:
  containers:
    - name: resource-limits-demo
      image: nginx:alpine
      ports:
        - containerPort: 80
      resources:
        requests:         # requests 集群调度使用的资源
          memory: 50Mi    # 定义容器或Pod对内存资源的需求量
          cpu: 50m        # 定义容器或Pod对CPU资源的需求量
        limits:           # limits 容器资源限制的配置
          memory: 100Mi   # 限制容器或Pod可以使用的内存资源的最大数量
          cpu: 100m       # 限制容器或Pod可以使用的CPU资源的最大数量
```





## Reference

