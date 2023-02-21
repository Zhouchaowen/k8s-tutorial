# Pod调度

```bash
# 给Node打标签
kubectl label nodes <node-name> <key=value>
```





## Deployment



### Pod Scale

```yaml
# nginx-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name:  nginx-deploy
  namespace: default
spec:
  replicas: 3     # 期望的 Pod 副本数量，默认值为1
  selector:       # Label Selector，必须匹配 Pod 模板中的标签
    matchLabels:
      app: nginx
  template:  # Pod 模板
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:alpine
          ports:
            - containerPort: 80
```

验证

```bash
# 创建
kubectl apply -f tmp.yaml 

# 查看pod信息
kubectl  get pods -o wide

# 查看deploy详细信息
kubectl describe deploy

# 查看rs详细信息
kubectl describe rs

# 查看pod详细信息
kubectl describe pod [podName]
```

扩容

```bash
kubectl scale deployment nginx-deploy --replicas 5

kubectl get pods --watch -o wide
#NAME                           READY   STATUS    AGE   IP           NODE
#nginx-deploy-f7ccf9478-9twq9   1/1     Running   28m   172.17.2.6   k8s-node02
#nginx-deploy-f7ccf9478-dhmn9   1/1     Running   28m   172.17.2.5   k8s-node02
#nginx-deploy-f7ccf9478-r52f8   1/1     Running   28m   172.17.1.5   k8s-node01

#nginx-deploy-f7ccf9478-gr2wm   0/1     Pending   0s    <none>       <none>
#nginx-deploy-f7ccf9478-dd5cf   0/1     Pending   0s    <none>       <none>
#nginx-deploy-f7ccf9478-gr2wm   0/1     Pending   0s    <none>       k8s-node01
#nginx-deploy-f7ccf9478-dd5cf   0/1     Pending   0s    <none>       k8s-node01
#nginx-deploy-f7ccf9478-gr2wm   0/1     ContainerCreating   0s    <none>       k8s-node01
#nginx-deploy-f7ccf9478-dd5cf   0/1     ContainerCreating   0s    <none>       k8s-node01
#nginx-deploy-f7ccf9478-gr2wm   1/1     Running             3s    172.17.1.6   k8s-node01
#nginx-deploy-f7ccf9478-dd5cf   1/1     Running             3s    172.17.1.7   k8s-node01
```

缩容

```bash
kubectl scale deployment nginx-deploy --replicas=1

#kubectl get pods --watch -o wide
#NAME                           READY   STATUS        AGE     IP           NODE
#nginx-deploy-f7ccf9478-9twq9   1/1     Running       33m     172.17.2.6   k8s-node02
#nginx-deploy-f7ccf9478-dd5cf   1/1     Running       4m52s   172.17.1.7   k8s-node01
#nginx-deploy-f7ccf9478-dhmn9   1/1     Running       33m     172.17.2.5   k8s-node02
#nginx-deploy-f7ccf9478-gr2wm   1/1     Running       4m52s   172.17.1.6   k8s-node01
#nginx-deploy-f7ccf9478-r52f8   1/1     Running       33m     172.17.1.5   k8s-node01
#nginx-deploy-f7ccf9478-r52f8   1/1     Terminating   33m     172.17.1.5   k8s-node01
#nginx-deploy-f7ccf9478-dhmn9   1/1     Terminating   33m     172.17.2.5   k8s-node02
#nginx-deploy-f7ccf9478-gr2wm   1/1     Terminating   4m58s   172.17.1.6   k8s-node01
#nginx-deploy-f7ccf9478-dd5cf   1/1     Terminating   4m58s   172.17.1.7   k8s-node01
#nginx-deploy-f7ccf9478-r52f8   0/1     Terminating   33m     172.17.1.5   k8s-node01
#nginx-deploy-f7ccf9478-r52f8   0/1     Terminating   33m     172.17.1.5   k8s-node01
#nginx-deploy-f7ccf9478-r52f8   0/1     Terminating   33m     172.17.1.5   k8s-node01
#nginx-deploy-f7ccf9478-gr2wm   0/1     Terminating   4m59s   172.17.1.6   k8s-node01
#nginx-deploy-f7ccf9478-gr2wm   0/1     Terminating   4m59s   172.17.1.6   k8s-node01
#nginx-deploy-f7ccf9478-gr2wm   0/1     Terminating   4m59s   172.17.1.6   k8s-node01
#nginx-deploy-f7ccf9478-dd5cf   0/1     Terminating   4m59s   172.17.1.7   k8s-node01
#nginx-deploy-f7ccf9478-dd5cf   0/1     Terminating   4m59s   172.17.1.7   k8s-node01
#nginx-deploy-f7ccf9478-dd5cf   0/1     Terminating   4m59s   172.17.1.7   k8s-node01
#nginx-deploy-f7ccf9478-dhmn9   0/1     Terminating   33m     172.17.2.5   k8s-node02
#nginx-deploy-f7ccf9478-dhmn9   0/1     Terminating   33m     172.17.2.5   k8s-node02
#nginx-deploy-f7ccf9478-dhmn9   0/1     Terminating   33m     172.17.2.5   k8s-node02
```









### Pod Update

如果Pod是通过Deployment创建的，则用户可以在运行时修改Deployment的Pod定义（spec.template）或镜像名称，并应用到Deployment对象上(apply)，系统即可完成Deployment的rollout动作，rollout可被视为Deployment的自动更新或者自动部署动作。

如果在更新过程中发生了错误，则还可以通过回滚操作恢复Pod的版本。

只有Pod模板定义部分（Deployment的.spec.template）的属性发生改变时才会触发Deployment的rollout行为，对于其他的比如修改Pod的副本数量（spec.replicas）的值，则不会触发rollout行为。

```yaml
# nginx-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name:  nginx-deploy
  namespace: default
spec:
  replicas: 3     # 期望的 Pod 副本数量，默认值为1
  selector:       # Label Selector，必须匹配 Pod 模板中的标签
    matchLabels:
      app: nginx
  template:  # Pod 模板
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:alpine
          ports:
            - containerPort: 80
```

验证

```bash
# 创建
kubectl apply -f tmp.yaml 

# 查看pod信息
kubectl  get pods -o wide

# 查看deploy详细信息
kubectl describe deploy

# 查看rs详细信息
kubectl describe rs

# 查看pod详细信息
kubectl describe pod [podName]
```

升级

```bash
# update
kubectl set image deployment/nginx-deploy nginx=nginx:1.9.1
# or
kubectl edit deployment/nginx-deploy

# 窗口2查看 pod被替换的过程
kubectl get pods --watch
#NAME                           READY   STATUS    RESTARTS   AGE
#nginx-deploy-f7ccf9478-czblt   1/1     Running   0          8m2s
#nginx-deploy-f7ccf9478-r77n6   1/1     Running   0          8m2s
#nginx-deploy-f7ccf9478-rdhwv   1/1     Running   0          8m2s

#nginx-deploy-5bfdf46dc6-zz4b8   0/1     Pending   0          0s
#nginx-deploy-5bfdf46dc6-zz4b8   0/1     Pending   0          0s
#nginx-deploy-5bfdf46dc6-zz4b8   0/1     ContainerCreating   0          0s
#nginx-deploy-5bfdf46dc6-zz4b8   1/1     Running             0          2s
#nginx-deploy-f7ccf9478-czblt    1/1     Terminating         0          8m47s
#nginx-deploy-5bfdf46dc6-2h7sk   0/1     Pending             0          0s
#nginx-deploy-5bfdf46dc6-2h7sk   0/1     Pending             0          0s
#nginx-deploy-5bfdf46dc6-2h7sk   0/1     ContainerCreating   0          0s
#nginx-deploy-f7ccf9478-czblt    0/1     Terminating         0          8m49s
#nginx-deploy-f7ccf9478-czblt    0/1     Terminating         0          8m49s
#nginx-deploy-f7ccf9478-czblt    0/1     Terminating         0          8m49s
#nginx-deploy-5bfdf46dc6-2h7sk   1/1     Running             0          3s
#nginx-deploy-f7ccf9478-r77n6    1/1     Terminating         0          8m50s
#nginx-deploy-5bfdf46dc6-9gb5z   0/1     Pending             0          0s
#nginx-deploy-5bfdf46dc6-9gb5z   0/1     Pending             0          0s
#nginx-deploy-5bfdf46dc6-9gb5z   0/1     ContainerCreating   0          0s
#nginx-deploy-f7ccf9478-r77n6    0/1     Terminating         0          8m51s
#nginx-deploy-f7ccf9478-r77n6    0/1     Terminating         0          8m51s
#nginx-deploy-f7ccf9478-r77n6    0/1     Terminating         0          8m51s
#nginx-deploy-5bfdf46dc6-9gb5z   1/1     Running             0          2s
#nginx-deploy-f7ccf9478-rdhwv    1/1     Terminating         0          8m52s
#nginx-deploy-f7ccf9478-rdhwv    0/1     Terminating         0          8m53s
#nginx-deploy-f7ccf9478-rdhwv    0/1     Terminating         0          8m53s
#nginx-deploy-f7ccf9478-rdhwv    0/1     Terminating         0          8m53s

kubectl get rs
#NAME                      DESIRED   CURRENT   READY   AGE
#nginx-deploy-5bfdf46dc6   3         3         3       26s
#nginx-deploy-f7ccf9478    0         0         0       9m11s

# deploy操作ScalingReplicaSet替换pod的过程
kubectl describe deploy
#From                   Message
#----                   -------
#deployment-controller  Scaled up replica set nginx-deploy-f7ccf9478 to 3
#deployment-controller  Scaled up replica set nginx-deploy-5bfdf46dc6 to 1
#deployment-controller  Scaled down replica set nginx-deploy-f7ccf9478 to 2
#deployment-controller  Scaled up replica set nginx-deploy-5bfdf46dc6 to 2
#deployment-controller  Scaled down replica set nginx-deploy-f7ccf9478 to 1
#deployment-controller  Scaled up replica set nginx-deploy-5bfdf46dc6 to 3
#deployment-controller  Scaled down replica set nginx-deploy-f7ccf9478 to 0
```

回滚

```bash
# 查看历史版本
kubectl rollout history deployment/nginx-deploy

# 查看指定历史版本详情
kubectl rollout history deployment/nginx-deploy --revision=3

# 退回上一个版本
kubectl rollout undo deployment/nginx-deploy

# 退回指定版本
kubectl rollout undo deployment/nginx-deploy --to-revision=2

kubectl get pods --watch
#NAME                            READY   STATUS    RESTARTS   AGE
#nginx-deploy-5bfdf46dc6-2h7sk   1/1     Running   0          7m15s
#nginx-deploy-5bfdf46dc6-9gb5z   1/1     Running   0          7m12s
#nginx-deploy-5bfdf46dc6-zz4b8   1/1     Running   0          7m17s

#nginx-deploy-f7ccf9478-dhmn9    0/1     Pending   0          0s
#nginx-deploy-f7ccf9478-dhmn9    0/1     Pending   0          1s
#nginx-deploy-f7ccf9478-dhmn9    0/1     ContainerCreating   0          1s
#nginx-deploy-f7ccf9478-dhmn9    1/1     Running             0          3s
#nginx-deploy-5bfdf46dc6-2h7sk   1/1     Terminating         0          20m
#nginx-deploy-f7ccf9478-r52f8    0/1     Pending             0          0s
#nginx-deploy-f7ccf9478-r52f8    0/1     Pending             0          0s
#nginx-deploy-f7ccf9478-r52f8    0/1     ContainerCreating   0          0s
#nginx-deploy-5bfdf46dc6-2h7sk   0/1     Terminating         0          20m
#nginx-deploy-5bfdf46dc6-2h7sk   0/1     Terminating         0          20m
#nginx-deploy-5bfdf46dc6-2h7sk   0/1     Terminating         0          20m
#nginx-deploy-f7ccf9478-r52f8    1/1     Running             0          2s
#nginx-deploy-5bfdf46dc6-zz4b8   1/1     Terminating         0          20m
#nginx-deploy-f7ccf9478-9twq9    0/1     Pending             0          0s
#nginx-deploy-f7ccf9478-9twq9    0/1     Pending             0          0s
#nginx-deploy-f7ccf9478-9twq9    0/1     ContainerCreating   0          0s
#nginx-deploy-5bfdf46dc6-zz4b8   0/1     Terminating         0          20m
#nginx-deploy-5bfdf46dc6-zz4b8   0/1     Terminating         0          20m
#nginx-deploy-5bfdf46dc6-zz4b8   0/1     Terminating         0          20m
#nginx-deploy-f7ccf9478-9twq9    1/1     Running             0          2s
#nginx-deploy-5bfdf46dc6-9gb5z   1/1     Terminating         0          20m
#nginx-deploy-5bfdf46dc6-9gb5z   0/1     Terminating         0          20m
#nginx-deploy-5bfdf46dc6-9gb5z   0/1     Terminating         0          20m
#nginx-deploy-5bfdf46dc6-9gb5z   0/1     Terminating         0          20m

kubectl describe deploy
#From                   Message
#----                   -------
#deployment-controller  Scaled up replica set nginx-deploy-f7ccf9478 to 1
#deployment-controller  Scaled down replica set nginx-deploy-5bfdf46dc6 to 2
#deployment-controller  Scaled up replica set nginx-deploy-f7ccf9478 to 2
#deployment-controller  Scaled up replica set nginx-deploy-f7ccf9478 to 3
#deployment-controller  Scaled down replica set nginx-deploy-5bfdf46dc6 to 1
#deployment-controller  Scaled down replica set nginx-deploy-5bfdf46dc6 to 0
```

暂停更新操作

```bash
# pause and resume
# 暂停
kubectl rollout pause deployment/nginx-deploy

# 修改,修改,修改
kubectl set image deploy/nginx-deploy nginx=nginx:1.9.1
kubectl rollout history deploy/nginx-deploy
kubectl set resources deployment nginx-deploy -c=nginx --limits=cpu=200m,memory=512Mi

# 恢复
kubectl rollout resume deploy nginx-deploy#
```

更新属性设置

在Deployment的定义中，可以通过spec.strategy指定Pod更新的策略，目前支持两种策略：Recreate（重建）和RollingUpdate（滚动更新），默认值为RollingUpdate。

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name:  nginx-deploy
  namespace: default
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  minReadySeconds: 5
  strategy:
    type: RollingUpdate  # 指定更新策略：RollingUpdate和Recreate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:alpine
          ports:
            - containerPort: 80
```

- spec.strategy.type：
  - Recreate：设置spec.strategy.type=Recreate，表示Deployment在更新Pod时，会先“杀掉”所有正在运行的Pod，然后创建新的Pod。
  - RollingUpdate：设置spec.strategy.type=RollingUpdate，表示Deployment会以滚动更新的方式来逐个更新Pod。同时，可以通过设置spec.strategy.rollingUpdate下的两个参数（maxUnavailable和maxSurge）来控制滚动更新的过程。

- spec.strategy.rollingUpdate.maxUnavailable：用于指定Deployment在更新过程中不可用状态的Pod数量的上限。例如：maxUnavaible=1，则表示 Kubernetes 整个升级过程中最多会有1个 Pod 处于无法服务的状态。
- spec.strategy.rollingUpdate.maxSurge：用于指定在Deployment更新Pod的过程中Pod总数量超过Pod期望副本数量部分的最大值。



### HPA

Kubernetes中的某个Metrics Server持续采集所有Pod副本的指标数据。HPA控制器通过Metrics Server的API获取这些数据，基于用户定义的扩缩容规则进行计算，得到目标Pod的副本数量。当目标Pod副本数量与当前副本数量不同时，HPA控制器就向Pod的副本控制器（Deployment、RC或ReplicaSet）发起scale操作，调整Pod的副本数量，完成扩缩容操作。



安装metrics-server

```bash
# 国内需要替换镜像名称 并忽略证书- --kubelet-insecure-tls
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# https://zhuanlan.zhihu.com/p/572406293
```

- 注意网络问题
- 注意apiserver访问metrics的ip地址，是集群地址还是内部地址还是外部地址。这需要配置
- 注意apiserver需要开启- --enable-aggregator-routing=true

```bash
# metrics
- --kubelet-preferred-address-types=InternalIP,Hostname,InternalDNS,ExternalDNS,ExternalIP  
- --kubelet-insecure-tls

# apiserver
- --enable-aggregator-routing=true

# https://www.cnblogs.com/fat-girl-spring/p/15936467.html
# https://www.cnblogs.com/shunzi115/p/12438702.html
```

验证

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: php
spec:
  selector:
    matchLabels:
      run: php
  replicas: 1
  template:
    metadata:
      labels:
        run: php
    spec:
      containers:
        - name: php
          image: hpa-php:test
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 80
          resources:
            limits:
              cpu: 400m
            requests:
              cpu: 200m

---
apiVersion: v1
kind: Service
metadata:
  name: php
  labels:
    run: php
spec:
  ports:
    - port: 80
  selector:
    run: php
```

```bash
kubectl autoscale deployment php --cpu-percent=50 --min=1 --max=10

kubectl get hpa

kubectl run -i --tty load-generator --rm --image=busybox:1.28 --restart=Never -- /bin/sh -c "while sleep 0.01; do wget -q -O- http://php-apache; done"

kubectl get hpa php --watch

# 这里启动一个容器，并将无限查询循环发送php服务
kubectl run v1 -it --image=busybox:1.28 /bin/sh
# 登录到容器，执行以下操作
/ # while true; do wget -q -O- http://php.default.svc.cluster.local; done
```

参考：https://www.cnblogs.com/yuhaohao/p/14109787.html























