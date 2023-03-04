

# Deployment

> 常规业务：扩缩容,滚动更新

Deployment 使得 Pod 和 ReplicaSet 能够进行声明式更新。

------

- **apiVersion**: apps/v1
- **kind**: Deployment

- **metadata** ([ObjectMeta](https://kubernetes.io/zh-cn/docs/reference/kubernetes-api/common-definitions/object-meta/#ObjectMeta))
- **spec** ([DeploymentSpec](https://kubernetes.io/zh-cn/docs/reference/kubernetes-api/workload-resources/deployment-v1/#DeploymentSpec)) Deployment 预期行为的规约。
- **status** ([DeploymentStatus](https://kubernetes.io/zh-cn/docs/reference/kubernetes-api/workload-resources/deployment-v1/#DeploymentStatus)) 最近观测到的 Deployment 状态。

DeploymentSpec 定义 Deployment 预期行为的规约。

------

- **selector** ([LabelSelector](https://kubernetes.io/zh-cn/docs/reference/kubernetes-api/common-definitions/label-selector/#LabelSelector))，必需

  供 Pod 所用的标签选择算符。通过此字段选择现有 ReplicaSet 的 Pod 集合， 被选中的 ReplicaSet 将受到这个 Deployment 的影响。此字段必须与 Pod 模板的标签匹配。

- **template** ([PodTemplateSpec](https://kubernetes.io/zh-cn/docs/reference/kubernetes-api/workload-resources/pod-template-v1/#PodTemplateSpec))，必需

  template 描述将要创建的 Pod。

- **replicas** (int32)

  预期 Pod 的数量。这是一个指针，用于辨别显式零和未指定的值。默认为 1。

- **minReadySeconds** (int32)

  新建的 Pod 在没有任何容器崩溃的情况下就绪并被系统视为可用的最短秒数。 默认为 0（Pod 就绪后即被视为可用）。

- **strategy** (DeploymentStrategy)

  **补丁策略：retainKeys**

  将现有 Pod 替换为新 Pod 时所用的部署策略。

  **DeploymentStrategy 描述如何将现有 Pod 替换为新 Pod。**

  - **strategy.type** (string)

    部署的类型。取值可以是 “Recreate” 或 “RollingUpdate”。默认为 RollingUpdate。

  - **strategy.rollingUpdate** (RollingUpdateDeployment)

    滚动更新这些配置参数。仅当 type = RollingUpdate 时才出现。

    **控制滚动更新预期行为的规约。**

    - **strategy.rollingUpdate.maxSurge** (IntOrString)

      超出预期的 Pod 数量之后可以调度的最大 Pod 数量。该值可以是一个绝对数（例如： 5）或一个预期 Pod 的百分比（例如：10%）。如果 MaxUnavailable 为 0，则此字段不能为 0。 通过向上取整计算得出一个百分比绝对数。默认为 25%。例如：当此值设为 30% 时， 如果滚动更新启动，则可以立即对 ReplicaSet 扩容，从而使得新旧 Pod 总数不超过预期 Pod 数量的 130%。 一旦旧 Pod 被杀死，则可以再次对新的 ReplicaSet 扩容， 确保更新期间任何时间运行的 Pod 总数最多为预期 Pod 数量的 130%。

      **IntOrString 是可以保存 int32 或字符串的一个类型。 当用于 JSON 或 YAML 编组和取消编组时，它会产生或消费内部类型。 例如，这允许你拥有一个可以接受名称或数值的 JSON 字段。**

    - **strategy.rollingUpdate.maxUnavailable** (IntOrString)

      更新期间可能不可用的最大 Pod 数量。该值可以是一个绝对数（例如： 5）或一个预期 Pod 的百分比（例如：10%）。通过向下取整计算得出一个百分比绝对数。 如果 MaxSurge 为 0，则此字段不能为 0。默认为 25%。 例如：当此字段设为 30%，则在滚动更新启动时 ReplicaSet 可以立即缩容为预期 Pod 数量的 70%。 一旦新的 Pod 就绪，ReplicaSet 可以再次缩容，接下来对新的 ReplicaSet 扩容， 确保更新期间任何时间可用的 Pod 总数至少是预期 Pod 数量的 70%。

      **IntOrString 是可以保存 int32 或字符串的一个类型。 当用于 JSON 或 YAML 编组和取消编组时，它会产生或消费内部类型。 例如，这允许你拥有一个可以接受名称或数值的 JSON 字段。**

- **revisionHistoryLimit** (int32)

  保留允许回滚的旧 ReplicaSet 的数量。这是一个指针，用于辨别显式零和未指定的值。默认为 10。

- **progressDeadlineSeconds** (int32)

  Deployment 在被视为失败之前取得进展的最大秒数。Deployment 控制器将继续处理失败的部署， 原因为 ProgressDeadlineExceeded 的状况将被显示在 Deployment 状态中。 请注意，在 Deployment 暂停期间将不会估算进度。默认为 600s。

- **paused** (boolean)

  指示部署被暂停。

实战 https://kubernetes.io/zh-cn/docs/tasks/run-application/run-stateless-application-deployment/

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

# 查看deploy
kubectl get deploy

# 查看deploy详细信息
kubectl describe deploy

# 查看ReplicaSet
kubectl get rs

# 查看rs详细信息
kubectl describe rs


# 查看pod信息
kubectl  get pods -o wide

# 查看pod详细信息
kubectl describe pod [podName]
```

Pod Scale 扩容

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




