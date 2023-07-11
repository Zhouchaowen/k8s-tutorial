# Init-Container

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
# init-containers.yaml
apiVersion: v1
kind: Pod
metadata:
  name: init-containers
spec:
  volumes:              # volumes卷稍后介绍，用于临时存储数据
    - name: workdir
      emptyDir: {}
  initContainers:       # Pod在启动containers之前，先要【运行完】initContainers的所有容器，所以这些容器必须有终结，不能一直运行
    - name: install-index-html
      image: busybox:1.28
      command:          # 执行命令,定义后会覆盖 Dockerfile 中的 CMD,ENTRYPOINT 命令
        - /bin/sh
        - "-c"
        - "date > index.html"
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

```bash
# kubectl get pod -o wide
NAME              READY   STATUS            RESTARTS   AGE   IP            NODE         
init-containers   0/1     PodInitializing   0          8s    10.244.1.70   k8s-node01

# curl 10.244.1.70
Tue Jul 11 13:51:26 UTC 2023
```

## Reference

