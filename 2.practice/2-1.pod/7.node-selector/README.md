# Pod

Pod 里的所有容器，共享的是同一个 Network Namespace，并且可以声明共享同一个 Volume。

Pod可以由一个或多个容器组合而成, 属于同一个Pod的多个容器应用之间相互访问时仅需通过localhost就可以通信，使得这一组容器被“绑定”在一个环境中。

```yaml
# 例1: nginx-pod.yaml
apiVersion: v1                  #k8sApi版本 kubectl api-resource 可以获取
kind: Pod                       #资源类型
# 以上为type信息
metadata:                       #元数据:名称，命名空间，标签等
  name: ngx-pod                 #service的名称
  labels:                       #自定义标签属性列表
    env: demo
    owner: zcw
# 以上为元数据部分
spec:                           #期望状态：详细描述pod容器，卷等
  containers:
    - image: nginx:alpine       #镜像地址
      name: ngx                 #镜像名称
      ports:
        - containerPort: 80     #端口
# 以上为资源规格描述部分
```

创建Pod

```bash
kubectl create -f nginx-pod.yaml
```

查看Pod

```bash
kubectl get pods
```

查看Pod详情

```bash
kubectl describe pod ngx-pod
```

查看Pod日志

```bash
kubectl logs ngx-pod [-c containerName]
```

端口转发，访问Pod

```bash
kubectl port-forward ngx-pod 8080:80
```

通过标签选择Pod

```bash
kubectl get pod -l env
```