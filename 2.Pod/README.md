## Pod

Pod可以由1个或多个容器组合而成, 属于同一个Pod的多个容器应用之间相互访问时仅需通过localhost就可以通信，使得这一组容器被“绑定”在一个环境中。

```yaml
# 例1: nginx-pod.yaml
apiVersion: v1                  #版本
kind: Pod                       #资源类型
metadata:                       #元数据
  name: ngx-pod                 #service的名称
  labels:                       #自定义标签属性列表
    env: demo
    owner: zcw
spec:                           #详细描述
  containers:
    - image: nginx:alpine       #镜像地址
      name: ngx                 #镜像名称
      ports:
        - containerPort: 80     #端口
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
kubectl describe pod podNameXXX
```

## init-container

init container与应用容器在本质上是一样的，但它们是仅运行一次就结束的任务，并且必须在成功运行完成后，系统才能继续执行下一个容器。当设置了多个init container时，将按顺序逐个运行，并且只有前一个init container运行成功后才能运行后一个init container。

- 等待其他关联组件正确运行（例如数据库或某个后台服务）。
- 基于环境变量或配置模板生成配置文件。
- 从远程数据库获取本地所需配置，或者将自身注册到某个中央数据库中。
- 下载相关依赖包，或者对系统进行一些预配置操作。

```yaml
# init-container-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: init-demo
spec:
  volumes:
    - name: workdir
      emptyDir: {}
  initContainers:       # 提前初始化容器命令
    - name: install
      image: busybox
      command:          # 这些命令
        - wget
        - "-O"
        - "/work-dir/index.html"
        - http://www.baidu.com
      volumeMounts:
        - name: workdir
          mountPath: "/work-dir"
  containers:
    - name: nginx
      image: nginx:alpine
      ports:
        - containerPort: 80
      volumeMounts:
        - name: workdir
          mountPath: /usr/share/nginx/html
```

## 启动

```bash
kubectl create -f nginx-pod.yaml

kubectl get pods

kubectl describe pod podNameXXX
```

