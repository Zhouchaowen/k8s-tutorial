# PV-PVC

创建的 PVC 要真正被容器使用起来，就必须先和某个符合条件的 PV 进行绑定。这里要检查的条件，包括两部分：

- 第一个条件，当然是 PV 和 PVC 的 spec 字段。比如，PV 的存储（storage）大小，就必须满足 PVC 的要求。
- 第二个条件，则是 PV 和 PVC 的 storageClassName 字段必须一样。

PVC 可以理解为持久化存储的“接口”，它提供了对某种持久化存储的描述，但不提供具体的实现；而这个持久化存储的实现部分则由 PV 负责完成。这样做的好处是，作为应用开发者，我们只需要跟 PVC 这个“接口”打交道，而不必关心具体的实现是 NFS 还是 Ceph。

“持久化 Volume”，指的就是这个宿主机上的目录，具备“持久性”。即：这个目录里面的内容，既不会因为容器的删除而被清理掉，也不会跟当前的宿主机绑定。这样，当容器被重启或者在其他节点上重建出来之后，它仍然能够通过挂载这个 Volume，访问到这些内容。

Static Provisioning

```yaml
# 定义Namespace
---
apiVersion: v1
kind: Namespace       # 资源类型
metadata:
  name: dev  # 命名空间名称 kube-example
# 定义PV
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nginx-pv    # pv的名称
spec:
  accessModes:      # 访问模式
    - ReadWriteMany # PV以read-write挂载到多个节点
  capacity:         # 容量
    storage: 2Gi    # pv可用的大小
  nfs:
    path: /nfs/data/     # NFS的挂载路径
    server: 10.2.0.106    # NFS服务器地址

# 定义PVC，用于消费PV
---
apiVersion: v1
kind: PersistentVolumeClaim   # 类型
metadata:
  name: nginx-pvc             # PVC 的名字
  namespace: dev              # 命名空间
spec:
  accessModes:                # 访问模式
    - ReadWriteMany           # PVC以read-write挂载到多个节点
  resources:
    requests:
      storage: 2Gi            # PVC允许申请的大小

# 定义Pod，指定需要使用的PVC
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-pvc
  namespace: dev   # 如果前面的PVC指定了命名空间这里必须指定与PVC一致的命名空间，否则PVC不可用
spec:
  selector:
    matchLabels:
      app: nginx-pvc
  template:
    metadata:
      labels:
        app: nginx-pvc
    spec:
      containers:
        - name: nginx-test-pvc
          image: nginx:1.7.9
          imagePullPolicy: IfNotPresent
          ports:
            - name: web-port
              containerPort: 80
              protocol: TCP
          volumeMounts:
            - name: nginx-persistent-storage    # 取个名字，与下面的volumes的名字要一致
              mountPath: /usr/share/nginx/html  # 容器中的路径
      volumes:
        - name: nginx-persistent-storage
          persistentVolumeClaim:
            claimName: nginx-pvc  # 引用前面声明的PVC
```

Dynamic Provisioning

StorageClass



Local Persistent Volume

相比于正常的 PV，一旦这些节点宕机且不能恢复时，Local Persistent Volume 的数据就可能丢失。这就要求使用 Local Persistent Volume 的应用必须具备数据备份和恢复的能力，允许你把这些数据定时备份在其他位置。



一个 Local Persistent Volume 对应的存储介质，一定是一块额外挂载在宿主机的磁盘或者块设备（“额外”的意思是，它不应该是宿主机根目录所使用的主硬盘）。这个原则，我们可以称为“一个 PV 一块盘”。