# DownwardApi

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

## fieldRef

```yaml
# field-ref.yaml
apiVersion: v1
kind: Pod
metadata:
  name: field-ref
spec:
  containers:
    - name: field-ref
      image: busybox:1.28
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
```

```bash
# kubectl apply -f field-ref.yaml
pod/field-ref created
# kubectl get pod
NAME        READY   STATUS    RESTARTS   AGE
field-ref   1/1     Running   0          5s 
# kubectl logs field-ref

k8s-node01
field-ref
default
10.244.1.66
default
```

## resourceFieldRef

```yaml
# resource-field-ref.yaml
apiVersion: v1
kind: Pod
metadata:
  name: resource-field-ref
spec:
  containers:
    - name: resource-field-ref
      image: busybox:1.28
      command: [ "sh", "-c"]          # 执行命令行命令
      args:                           # 参数
        - while true; do
          echo -en '\n';
          printenv MY_CPU_REQUEST MY_CPU_LIMIT;
          printenv MY_MEM_REQUEST MY_MEM_LIMIT;
          sleep 10;
          done;
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
            resourceFieldRef:                   # 通过引用container本身属性当做值
              containerName: resource-field-ref # 引用的container的名称
              resource: requests.cpu            # 引用资源限制的字段
        - name: MY_CPU_LIMIT
          valueFrom:
            resourceFieldRef:
              containerName: resource-field-ref
              resource: limits.cpu
        - name: MY_MEM_REQUEST
          valueFrom:
            resourceFieldRef:
              containerName: resource-field-ref
              resource: requests.memory
        - name: MY_MEM_LIMIT
          valueFrom:
            resourceFieldRef:
              containerName: resource-field-ref
              resource: limits.memory
```

```bash
# kubectl apply -f resource-field-ref.yaml
pod/field-ref created
# kubectl get pod
NAME        READY   STATUS    RESTARTS   AGE
field-ref   1/1     Running   0          5s 
# kubectl logs resource-field-ref
1
1
33554432
67108864
```



## Reference

https://kubernetes.io/zh-cn/docs/concepts/workloads/pods/downward-api/

