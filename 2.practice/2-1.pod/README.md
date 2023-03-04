# Pod

Pod 里的所有容器，共享的是同一个 Network Namespace，并且可以声明共享同一个 Volume。

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

Init 容器是一种特殊容器，在Pod内的应用容器启动之前运行。与普通的容器非常像，除了如下两点：

- 它们总是运行到完成。
- 每个都必须在下一个启动之前成功完成。

init container与应用容器在本质上是一样的，但它们是仅运行一次就结束的任务，并且必须在成功运行完成后，系统才能继续执行下一个容器。当设置了多个init container时，将按顺序逐个运行，并且只有前一个init container运行成功后才能运行后一个init container。

如果 Pod 的 Init 容器失败，kubelet 会不断地重启该 Init 容器直到该容器成功为止。 然而，如果 Pod 对应的 `restartPolicy` 值为 "Never"，并且 Pod 的 Init 容器失败， 则 Kubernetes 会将整个 Pod 状态设置为失败。

因为 Init 容器可能会被重启、重试或者重新执行，所以 Init 容器的代码应该是幂等的。

应用：

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

执行

```bash
kubectl apply -f nginx-pod.yaml

kubectl get pods

kubectl describe pod podNameXXX

kubectl logs -f pod-name -c init-container-name 
```

## lifecycle

生命周期事件包括 postStart 和 preStop函数

**PostStart**

这个回调在容器被创建之后立即被执行。 但是，不能保证回调会在容器入口点（ENTRYPOINT）之前执行。 没有参数传递给处理程序。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: hook-demo1
spec:
  containers:
    - name: hook-demo1
      image: nginx:alpine
      lifecycle:
        postStart:    # 容器创建后立即执行
          exec:
            command: ["/bin/sh", "-c", "echo Hello from the postStart handler > /usr/share/message"]
```

**PreStop**

在容器因 API 请求或者管理事件（诸如存活态探针、启动探针失败、资源抢占、资源竞争等） 而被终止之前，此回调会被调用。 如果容器已经处于已终止或者已完成状态，则对 preStop 回调的调用将失败。 在用来停止容器的 TERM 信号被发出之前，回调必须执行结束。 Pod 的终止宽限周期在 `PreStop` 回调被执行之前即开始计数， 所以无论回调函数的执行结果如何，容器最终都会在 Pod 的终止宽限期内被终止。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: hook-demo2
spec:
  containers:
    - name: hook-demo2
      image: nginx:alpine
      lifecycle:
        preStop:      # 容器终止之前立即被调用
          exec:
            command: ["/usr/sbin/nginx","-s","quit"]  # 优雅退出
```

## Limit-resource

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
        requests:     # requests 集群调度使用的资源
          memory: 50Mi
          cpu: 50m
        limits:       # limits 容器资源限制的配置
          memory: 100Mi
          cpu: 100m
```

## Downward API

- 环境变量：将Pod或Container信息设置为容器内的环境变量。
- Volume挂载：将Pod或Container信息以文件的形式挂载到容器内部。

可以使用 `fieldRef` 传递来自可用的 Pod 级字段的信息。

- `metadata.name`Pod 的名称
- `metadata.namespace`Pod 的命名空间
- `metadata.uid`Pod 的唯一 ID
- `metadata.annotations['<KEY>']`Pod 的注解`<KEY>` 的值（例如：`metadata.annotations['myannotation']`）
- `metadata.labels['<KEY>']`Pod 的标签`<KEY>` 的值（例如：`metadata.labels['mylabel']`）

以下信息可以通过环境变量获得，但**不能作为 downwardAPI 卷 fieldRef** 获得：

- `spec.serviceAccountName`Pod 的服务账号名称
- `spec.nodeName`Pod 运行时所处的节点
- `status.hostIP`Pod 所在节点的主 IP 地址
- `status.podIP`Pod 的主 IP 地址（通常是其 IPv4 地址）

以下信息可以通过 `downwardAPI` 卷 `fieldRef` 获得，但**不能作为环境变量**获得：

- `metadata.labels`Pod 的所有标签，格式为 `标签键名="转义后的标签值"`，每行一个标签
- `metadata.annotations`Pod 的全部注解，格式为 `注解键名="转义后的注解值"`，每行一个注解

可以使用 `resourceFieldRef` 传递来自可用的 Container 级字段的信息。

- `resource: limits.cpu`容器的 CPU 限制值
- `resource: requests.cpu`容器的 CPU 请求值
- `resource: limits.memory`容器的内存限制值
- `resource: requests.memory`容器的内存请求值
- `resource: limits.hugepages-*`容器的巨页限制值（前提是启用了 `DownwardAPIHugePages` 特性门控）
- `resource: requests.hugepages-*`容器的巨页请求值（前提是启用了 `DownwardAPIHugePages` 特性门控）
- `resource: limits.ephemeral-storage`容器的临时存储的限制值
- `resource: requests.ephemeral-storage`容器的临时存储的请求值

EVN方式

将Pod信息设置为容器内的环境变量，下面的例子通过Downward API将Pod的IP、名称和所在命名空间注入容器的环境变量中，Pod的YAML文件内容如下

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: d-api-enVars-field-ref
spec:
  containers:
    - name: test-container
      image: busybox
      command: [ "sh", "-c"]
      args:
        - while true; do
          echo -en '\n';
          printenv MY_NODE_NAME MY_POD_NAME MY_POD_NAMESPACE;
          printenv MY_POD_IP MY_POD_SERVICE_ACCOUNT;
          sleep 10;
          done;
      env:
        - name: MY_NODE_NAME
          valueFrom:              # 从pod获取spec.nodeName
            fieldRef:
              fieldPath: spec.nodeName
        - name: MY_POD_NAME
          valueFrom:              # 从pod获取metadata.name
            fieldRef:
              fieldPath: metadata.name
        - name: MY_POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: MY_POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        - name: MY_POD_SERVICE_ACCOUNT
          valueFrom:
            fieldRef:
              fieldPath: spec.serviceAccountName
  restartPolicy: Never
```

将Container信息设置为容器内的环境变量, 下面的例子通过Downward API将Container的资源请求和资源限制信息设置为容器内的环境变量，Pod的YAML文件内容如下：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: d-api-enVars-resource-field-ref
spec:
  containers:
    - name: test-container
      image: busybox
      imagePullPolicy: Never
      command: [ "sh", "-c"]          # 执行命令行命令
      args:                           # 参数
        - while true; do
          echo -en '\n';
          printenv MY_CPU_REQUEST MY_CPU_LIMIT;
          printenv MY_MEM_REQUEST MY_MEM_LIMIT;
          sleep 10;
          done;
      args:
        - while true; do
          echo -en '\n';
          printenv MY_CPU_REQUEST MY_CPU_LIMIT;
          printenv MY_MEM_REQUEST MY_MEM_LIMIT;
          sleep 3600;
          done;
      resources:                      # 设置资源限制 下线/上线
        requests:
          memory: "32Mi"
          cpu: "125m"
        limits:
          memory: "64Mi"
          cpu: "250m"
      env:
        - name: MY_CPU_REQUEST        # 定义容器环境变量
          valueFrom:
            resourceFieldRef:         # 通过引用container本身属性当做值
              containerName: test-container # 引用的container的名称
              resource: requests.cpu  # 引用资源限制的字段
        - name: MY_CPU_LIMIT
          valueFrom:
            resourceFieldRef:
              containerName: test-container
              resource: limits.cpu
        - name: MY_MEM_REQUEST
          valueFrom:
            resourceFieldRef:
              containerName: test-container
              resource: requests.memory
        - name: MY_MEM_LIMIT
          valueFrom:
            resourceFieldRef:
              containerName: test-container
              resource: limits.memory
  restartPolicy: Never
```

File方式

## Probe

针对运行中的容器，`kubelet` 可以选择是否执行以下三种探针，以及如何针对探测结果作出反应：

- `livenessProbe`

  指示容器是否正在运行。如果存活态探测失败，则 kubelet 会杀死容器， 并且容器将根据其[重启策略](https://kubernetes.io/zh-cn/docs/concepts/workloads/pods/pod-lifecycle/#restart-policy)决定未来。如果容器不提供存活探针， 则默认状态为 `Success`。

- `readinessProbe`

  指示容器是否准备好为请求提供服务。如果就绪态探测失败， 端点控制器将从与 Pod 匹配的所有服务的端点列表中删除该 Pod 的 IP 地址。 初始延迟之前的就绪态的状态值默认为 `Failure`。 如果容器不提供就绪态探针，则默认状态为 `Success`。

- `startupProbe`

  指示容器中的应用是否已经启动。如果提供了启动探针，则所有其他探针都会被 禁用，直到此探针成功为止。如果启动探测失败，`kubelet` 将杀死容器， 而容器依其[重启策略](https://kubernetes.io/zh-cn/docs/concepts/workloads/pods/pod-lifecycle/#restart-policy)进行重启。 如果容器没有提供启动探测，则默认状态为 `Success`。

使用探针来检查容器有四种不同的方法。 每个探针都必须准确定义为这四种机制中的一种：

- `exec`

  在容器内执行指定命令。如果命令退出时返回码为 0 则认为诊断成功。

- `grpc`

  使用 [gRPC](https://grpc.io/) 执行一个远程过程调用。 目标应该实现 [gRPC 健康检查](https://grpc.io/grpc/core/md_doc_health-checking.html)。 如果响应的状态是 "SERVING"，则认为诊断成功。 gRPC 探针是一个 Alpha 特性，只有在你启用了 "GRPCContainerProbe" [特性门控](https://kubernetes.io/zh-cn/docs/reference/command-line-tools-reference/feature-gates/)时才能使用。

- `httpGet`

  对容器的 IP 地址上指定端口和路径执行 HTTP `GET` 请求。如果响应的状态码大于等于 200 且小于 400，则诊断被认为是成功的。

- `tcpSocket`

  对容器的 IP 地址上的指定端口执行 TCP 检查。如果端口打开，则诊断被认为是成功的。 如果远程系统（容器）在打开连接后立即将其关闭，这算作是健康的。

**存活探针 livenessProbe**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: liveness-exec
spec:
  containers:
    - name: liveness
      image: busybox
      args:
        - /bin/sh
        - -c
        - touch /tmp/healthy; sleep 30; rm -rf /tmp/healthy; sleep 600
      livenessProbe:      # 存活探针
        exec:							# exec：执行一段命令
          command:
            - cat
            - /tmp/healthy
        initialDelaySeconds: 5  # 第一次执行探针的时候要等待5秒
        periodSeconds: 5        # 每隔5秒执行一次存活探针

---
apiVersion: v1
kind: Pod
metadata:
  name: liveness-http
spec:
  containers:
    - name: liveness
      image: cnych/liveness
      args:
        - /server
      livenessProbe:
        httpGet:			# http：检测某个 http 请求
          path: /healthz
          port: 8080
          httpHeaders:
            - name: X-Custom-Header
              value: Awesome
        initialDelaySeconds: 3
        periodSeconds: 3
```

**就绪探针 readinessProbe**

**慢启动探针 startupProbe**

任何给定的 Pod 从不会被“重新调度（rescheduled）”到不同的节点； 相反，这一 Pod 可以被一个新的、几乎完全相同的 Pod 替换掉。

如果某物声称其生命期与某 Pod 相同，例如存储[卷](https://kubernetes.io/zh-cn/docs/concepts/storage/volumes/)， 这就意味着该对象在此 Pod 存在期间也一直存在。 如果 Pod 因为任何原因被删除，甚至某完全相同的替代 Pod 被创建时， 这个相关的对象（例如这里的卷）也会被删除并重建。列如：Volume



