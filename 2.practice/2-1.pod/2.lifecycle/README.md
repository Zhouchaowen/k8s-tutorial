# Lifecycle

生命周期事件包括 postStart 和 preStop函数

## postStart

这个回调在容器被创建之后立即被执行。 但是，不能保证回调会在容器入口点（ENTRYPOINT）之前执行。 没有参数传递给处理程序。

- exec

```yaml
# post-start-exec.yaml
apiVersion: v1
kind: Pod
metadata:
  name: ps-exec
spec:
  containers:
    - name: post-start-exec
      image: nginx:alpine
      lifecycle:
        postStart:    # 容器创建后立即执行
          exec:
            command:
              - "/bin/sh"
              - "-c"
              - "echo Hello from the postStart handler > /usr/share/message"
```

```bash
kubectl apply -f post-start-exec.yaml
```

```bash
# kubectl exec -it ps-exec -- /bin/sh

/ # cat /usr/share/message 
Hello from the postStart handler
```

- httpGet

```yaml
postStart:
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
postStart:
  tcpSocket:
    host: 10.0.0.2
    port: 8080
```

##preStop

在容器因 API 请求或者管理事件（诸如存活态探针、启动探针失败、资源抢占、资源竞争等） 而被终止之前，此回调会被调用。 如果容器已经处于已终止或者已完成状态，则对 preStop 回调的调用将失败。 在用来停止容器的 TERM 信号被发出之前，回调必须执行结束。 Pod 的终止宽限周期在 `PreStop` 回调被执行之前即开始计数， 所以无论回调函数的执行结果如何，容器最终都会在 Pod 的终止宽限期内被终止。

- exec

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pre-stop
spec:
  volumes:
    - name: message
      hostPath:
        path: /tmp
  containers:
    - name: pre-stop
      image: nginx:alpine
      ports:
        - containerPort: 80
      volumeMounts:
        - name: message
          mountPath: /usr/share/
      lifecycle:
        preStop:
          exec:
            command:
              - "/bin/sh"
              - "-c"
              - "echo Hello from the postStart handler > /usr/share/message"
#          httpGet:
#            host: 10.0.0.2
#            httpHeaders:
#              - name: xx
#                value: xx
#            path: /
#            port: 8080
#          tcpSocket:
#            host: 10.0.0.2
#            port: 8080
```





## Reference

