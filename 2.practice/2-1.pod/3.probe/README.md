# Probe

```bash
kubectl explain pod.spec.containers.livenessProbe
kubectl explain pod.spec.containers.readinessProbe
```

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

## livenessProbe

- exec

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: liveness-exec
spec:
  containers:
    - name: liveness-exec
      image: busybox
      args:
        - /bin/sh
        - -c
        - touch /tmp/healthy; sleep 30; rm -rf /tmp/healthy; sleep 600
      livenessProbe:
        exec:           # 存活探针 exec 类型，通过在容器内执行命令并检查返回值来确定容器的健康状态
          command:      # 指定执行的命令
            - cat
            - /tmp/healthy
        terminationGracePeriodSeconds: 30 # K8S给你程序留的最后的缓冲时间，来处理关闭之前的操作
        initialDelaySeconds: 5  # 第一次执行探针的时候要等待5秒
        periodSeconds: 5        # 每隔5秒执行一次存活探针
        timeoutSeconds: 5       # 探测超时，到了超时时间探测还没返回结果说明失败
        successThreshold: 1     # 成功阈值，连续几次成才算成功
        failureThreshold: 3     # 失败阈值，连续几次失败才算真失败
```

- httpGet

```yaml
livenessProbe:
  httpGet:
    host: 10.0.0.2
    httpHeaders:
      - name: xx
      value: xx
    path: /
    port: 8080
```

- tcpSocket

```yaml
livenessProbe:
  tcpSocket:
    host: 10.0.0.2
    port: 8080
```

- grpc

```yaml
livenessProbe:
  tcpSocket:
    host: 10.2.0.2
    port: 8080
```

## readinessProbe

- exec

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: readiness-exec
spec:
  containers:
    - name: liveness
      image: busybox
      args:
        - /bin/sh
        - -c
        - touch /tmp/healthy; sleep 30; rm -rf /tmp/healthy; sleep 600
      readinessProbe:        
        exec:								# 就绪探针 exec 类型，通过在容器内执行命令并检查返回值来确定容器的健康状态
          command:
            - cat
            - /tmp/healthy
        initialDelaySeconds: 5  # 第一次执行探针的时候要等待5秒
        periodSeconds: 5        # 每隔5秒执行一次存活探针
        timeoutSeconds: 5       # 探测超时，到了超时时间探测还没返回结果说明失败
        successThreshold: 1     # 成功阈值，连续几次成才算成功
        failureThreshold: 3     # 失败阈值，连续几次失败才算真失败
```

- 同上




## Reference



