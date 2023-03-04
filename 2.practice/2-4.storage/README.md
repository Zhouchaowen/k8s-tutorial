# PV/PVC

`PV` 的全称是：`PersistentVolume`（持久化卷），是对底层共享存储的一种抽象

`PVC` 的全称是：`PersistentVolumeClaim`（持久化卷声明），PVC 是用户存储的一种声明，PVC 和 Pod 比较类似，Pod 消耗的是节点，PVC 消耗的是 PV 资源，Pod 可以请求 CPU 和内存，而 PVC 可以请求特定的存储空间和访问模式

`StorageClass`: 通过 `StorageClass` 的定义，管理员可以将存储资源定义为某种类型的资源，比如快速存储、慢速存储等，用户根据 StorageClass 的描述就可以非常直观的知道各种存储资源的具体特性了，这样就可以根据应用的特性去申请合适的存储资源了

```yaml
apiVersion: v1
kind: PersistentVolume        # 类型为PV
metadata:
  name: host-path-pv           # pv的名称
  labels:
    type: local
spec:
  storageClassName: manual
  capacity:
    storage: 1Gi             # pv可用的大小
  accessModes:
    - ReadWriteOnce          # PV以read-write挂载到
  hostPath:
    path: "/data/k8s/test/hostPath"
```

**storageClassName**就是刚才说过的，对存储类型的抽象 StorageClass。这个 PV 是我们手 

动管理的，名字可以任意起。

**accessModes**定义了存储设备的访问模式，简单来说就是虚拟盘的读写权限，和 Linux 的文 

件访问模式差不多，目前 Kubernetes 里有 3 种：

ReadWriteOnce：存储卷可读可写，但只能被一个节点上的 Pod 挂载。 

ReadOnlyMany：存储卷只读不可写，可以被任意节点上的 Pod 多次挂载。 

ReadWriteMany：存储卷可读可写，也可以被任意节点上的 Pod 多次挂载。

PersistentVolumeClaim

```yaml
# 定义PVC，用于消费PV
apiVersion: v1
kind: PersistentVolumeClaim   # 类型
metadata:
  name: host-path-pvc          # PVC的名称
spec:
  storageClassName: manual
  accessModes:                # 访问模式
    - ReadWriteOnce           # PVC以read-write挂载到一个节点
  resources:
    requests:
      storage: 1Gi            # PVC允许申请的大小
```

注意：

- PersistentVolume 简称为 PV，是 Kubernetes 对存储设备的抽象，由系统管理员维护，需 要描述清楚存储设备的类型、访问模式、容量等信息。 

- PersistentVolumeClaim 简称为 PVC，代表 Pod 向系统申请存储资源，它声明对存储的要 求，Kubernetes 会查找最合适的 PV 然后绑定。 

- StorageClass 抽象特定类型的存储系统，归类分组 PV 对象，用来简化 PV/PVC 的绑定过 程。 

- HostPath 是最简单的一种 PV，数据存储在节点本地，速度快但不能跟随 Pod 迁移。









## EmptyDir

当 Pod 分派到某个节点上时，`emptyDir` 卷会被创建，卷最初是空的,并且在 Pod 在该节点上运行期间，卷一直存在。 尽管 Pod 中的容器挂载 `emptyDir` 卷的路径可能相同也可能不同，这些容器都可以读写 `emptyDir` 卷中相同的文件。 当 Pod 因为某些原因被从节点上删除时，`emptyDir` 卷中的数据也会被永久删除。

`emptyDir` 的一些用途：

- 缓存空间，例如基于磁盘的归并排序。
- 为耗时较长的计算任务提供检查点，以便任务能方便地从崩溃前状态恢复执行。
- 在 Web 服务器容器服务数据时，保存内容管理器容器获取的文件。





## HostPath

`hostPath` 卷能将主机节点文件系统上的文件或目录挂载到你的 Pod 中。 虽然这不是大多数 Pod 需要的，但是它为一些应用程序提供了强大的逃生舱。

例如，`hostPath` 的一些用法有：

- 运行一个需要访问 Docker 内部机制的容器；可使用 `hostPath` 挂载 `/var/lib/docker` 路径。
- 在容器中运行 cAdvisor 时，以 `hostPath` 方式挂载 `/sys`。
- 允许 Pod 指定给定的 `hostPath` 在运行 Pod 之前是否应该存在，是否应该创建以及应该以什么方式存在。

除了必需的 `path` 属性之外，你可以选择性地为 `hostPath` 卷指定 `type`

同一个Pod中的多个容器能够共享Pod级别的存储卷Volume

```yaml
# 例2: share-volume-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: counter
spec:
  volumes:            # 定义名称为varlog的volume
    - name: varlog
      hostPath:
        path: /var/log/counter	# 主机目录
  containers:
    - name: count     # 容器 1
      image: busybox
      args:
        - /bin/sh
        - -c
        - >
          i=0;
          while true;
          do
            echo "$i: $(date)" >> /var/log/1.log;
            i=$((i+1));
            sleep 1;
          done
      volumeMounts:
        - name: varlog
          mountPath: /var/log
    - name: count-log   # 容器 2
      image: busybox
      args: [/bin/sh, -c, 'tail -n+1 -f /var/log/1.log']
      volumeMounts:
        - name: varlog
          mountPath: /var/log

# 判断是否需要在 Pod 中使用多个容器的时候，我们可以按照如下的几个方式来判断：
#
# 1.这些容器是否一定需要一起运行，是否可以运行在不同的节点上
# 2.这些容器是一个整体还是独立的组件
# 3.这些容器一起进行扩缩容会影响应用吗
```

### downwardAPI

`downwardAPI` 卷用于为应用提供 [downward API](https://kubernetes.io/docs/concepts/workloads/pods/downward-api/) 数据。 在这类卷中，所公开的数据以纯文本格式的只读文件形式存在。



## cephfs

`cephfs` 卷允许你将现存的 CephFS 卷挂载到 Pod 中。 不像 `emptyDir` 那样会在 Pod 被删除的同时也会被删除，`cephfs` 卷的内容在 Pod 被删除时会被保留，只是卷被卸载了。 这意味着 `cephfs` 卷可以被预先填充数据，且这些数据可以在 Pod 之间共享。同一 `cephfs` 卷可同时被多个写者挂载



## PV

`PV` 的全称是：`PersistentVolume`（持久化卷），是对底层共享存储的一种抽象

`PVC` 的全称是：`PersistentVolumeClaim`（持久化卷声明），PVC 是用户存储的一种声明，PVC 和 Pod 比较类似，Pod 消耗的是节点，PVC 消耗的是 PV 资源，Pod 可以请求 CPU 和内存，而 PVC 可以请求特定的存储空间和访问模式

`StorageClass`: 通过 `StorageClass` 的定义，管理员可以将存储资源定义为某种类型的资源，比如快速存储、慢速存储等，用户根据 StorageClass 的描述就可以非常直观的知道各种存储资源的具体特性了，这样就可以根据应用的特性去申请合适的存储资源了

### nfs

`nfs` 卷能将 NFS (网络文件系统) 挂载到你的 Pod 中。 不像 `emptyDir` 那样会在删除 Pod 的同时也会被删除，`nfs` 卷的内容在删除 Pod 时会被保存，卷只是被卸载。 这意味着 `nfs` 卷可以被预先填充数据，并且这些数据可以在 Pod 之间共享。

安装NFS

https://blog.csdn.net/bingju328/article/details/124489385

```bash
# 目标机器10.2.0.106

# 安装nfs
yum -y install nfs-utils rpcbind

# 创建nfs目录
mkdir -p /nfs/data/
mkdir -p /nfs/data/k8s

# 授予权限
chmod -R 777 /nfs/data

# 编辑export文件
# vim /etc/exports
/nfs/data *(rw,no_root_squash,sync)  # 这里给的是root权限---生产环境不推荐
# 或者/nfs/data   0.0.0.0/0(rw,sync,all_squash)  # 所有用户权限被映射成服务端上的普通用户nobody,权限被压缩

# 使得配置生效
exportfs -r

# 查看生效
exportfs

# 启动rpcbind、nfs服务
systemctl start rpcbind && systemctl enable rpcbind   #端口是111
systemctl start nfs && systemctl enable nfs           # 端口是 2049 

# 查看rpc服务的注册情况
rpcinfo -p localhost

# showmount测试
showmount -e ip(ip地址)
```

创建 namespace/pv/pvc/deploy

https://zhuanlan.zhihu.com/p/434209418

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

---
# 定义PVC，用于消费PV
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

---
# 定义Pod，指定需要使用的PVC
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

测试

```bash
# 在nginx-test-pvc pod的/usr/share/nginx/html目录写入index.html
echo "add another file to /nfs/data path " >  /usr/share/nginx/html/index.html

# 在nfs服务器/nfs/data目录下查看该文件

# 删除 nginx-test-pvc pod后再查看nfs服务器/nfs/data目录下index.html是否存在

# nginx-test-pvc pod被deployment查询拉取后查看是否能访问index.html
```



















## hostPath

创建pv

```yaml
# hostpath-pv.yaml
apiVersion: v1
kind: PersistentVolume        # 类型为PV
metadata:
  name: hostPath-pv           # pv的名称
  labels:
    type: local
spec:
  storageClassName: manual
  capacity:
    storage: 1Gi             # pv可用的大小
  accessModes:
    - ReadWriteOnce          # PV以read-write挂载到
  hostPath:
    path: "/data/k8s/test/hostPath"
```

配置文件中指定了该卷位于集群节点上的 `/data/k8s/test/hostPath` 目录，指定 1G 大小的空间和 `ReadWriteOnce` 的访问模式，这意味着该卷可以在单个节点上以读写方式挂载，另外定义了名称为 `manual` 的 `StorageClass`，该名称用来将 `PersistentVolumeClaim` 请求绑定到该 `PersistentVolum`。

- Capacity（存储能力）：一般来说，一个 PV 对象都要指定一个存储能力，通过 PV 的 `capacity` 属性来设置的，目前只支持存储空间的设置，就是我们这里的 `storage=10Gi`，不过未来可能会加入 `IOPS`、吞吐量等指标的配置。
- AccessModes（访问模式）：用来对 PV 进行访问模式的设置，用于描述用户应用对存储资源的访问权限，访问权限包括下面几种方式：
  - ReadWriteOnce（RWO）：读写权限，但是只能被单个节点挂载
  - ReadOnlyMany（ROX）：只读权限，可以被多个节点挂载
  - ReadWriteMany（RWX）：读写权限，可以被多个节点挂载

```bash
echo 'Hello from Kubernetes hostpath storage' > /data/k8s/test/hostpath/index.html
```



```bash
kubectl apply -f hostpath-pv.yaml

kubectl get pv hostPath-pv   
```

PV 的状态，实际上描述的是 PV 的生命周期的某个阶段，一个 PV 的生命周期中，可能会处于4种不同的阶段：

- Available（可用）：表示可用状态，还未被任何 PVC 绑定
- Bound（已绑定）：表示 PVC 已经被 PVC 绑定
- Released（已释放）：PVC 被删除，但是资源还未被集群重新声明
- Failed（失败）： 表示该 PV 的自动回收失败

创建pvc

```yaml
# 定义PVC，用于消费PV
apiVersion: v1
kind: PersistentVolumeClaim   # 类型
metadata:
  name: pvc-hostPath          # PVC的名称
spec:
  storageClassName: manual
  accessModes:                # 访问模式
    - ReadWriteOnce           # PVC以read-write挂载到一个节点
  resources:
    requests:
      storage: 1Gi            # PVC允许申请的大小
```

创建Pod

```yaml
# nginx-hostpath-pvc.yaml
# 定义PV
---
apiVersion: v1
kind: PersistentVolume        # 类型为PV
metadata:
  name: host-path-pv           # PV的名称
  labels:
    type: local
spec:
  storageClassName: manual
  capacity:
    storage: 1Gi             # PV可用的大小
  accessModes:
    - ReadWriteOnce          # PV以read-write挂载,只能被单个节点挂载
  hostPath:
    path: "/data/k8s/test/hostPath"

# 定义PVC，用于消费PV
---
apiVersion: v1
kind: PersistentVolumeClaim   # 类型
metadata:
  name: host-path-pvc          # PVC的名称
spec:
  storageClassName: manual
  accessModes:                # 访问模式
    - ReadWriteOnce           # PVC以read-write挂载,只能被单个节点挂载
  resources:
    requests:
      storage: 1Gi            # PVC允许申请的大小

# 定义Pod，指定需要使用的PVC
---
apiVersion: v1
kind: Pod
metadata:
  name: host-path-pvc-pod
spec:
  volumes:
    - name: host-path-pvc
      persistentVolumeClaim:
        claimName: host-path-pvc
  nodeSelector:                         # Pod只能调度到k8s-node01上
    kubernetes.io/hostname: k8s-node01
  containers:
    - name: nginx-host-path-pvc
      image: nginx:1.7.9
      ports:
        - containerPort: 80
      volumeMounts:
        - mountPath: "/usr/share/nginx/html"
          name: host-path-pvc
```

测试

```bash
kubectl apply -f nginx-hostpath-pvc.yaml


# 创建index.html文件
kubectl exec host-path-pvc-pod -- bash -c 'echo "this k8s-node01 hahah" > /usr/share/nginx/html/index.html'

# 访问Pod ip
curl

# 删除pod后,重建


# 访问Pod ip index.html内容是否相同
```









## NFS

安装NFS

https://blog.csdn.net/bingju328/article/details/124489385

```bash
# 目标机器10.2.0.106

# 安装nfs
yum -y install nfs-utils rpcbind

# 创建nfs目录
mkdir -p /nfs/data/
mkdir -p /nfs/data/k8s

# 授予权限
chmod -R 777 /nfs/data

# 编辑export文件
# vim /etc/exports
/nfs/data *(rw,no_root_squash,sync)  # 这里给的是root权限---生产环境不推荐
# 或者/nfs/data   0.0.0.0/0(rw,sync,all_squash)  # 所有用户权限被映射成服务端上的普通用户nobody,权限被压缩

# 使得配置生效
exportfs -r

# 查看生效
exportfs

# 启动rpcbind、nfs服务
systemctl start rpcbind && systemctl enable rpcbind   #端口是111
systemctl start nfs && systemctl enable nfs           # 端口是 2049 

# 查看rpc服务的注册情况
rpcinfo -p localhost

# showmount测试
showmount -e ip(ip地址)
```

创建 namespace/pv/pvc/deploy

https://zhuanlan.zhihu.com/p/434209418

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

---
# 定义PVC，用于消费PV
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

---
# 定义Pod，指定需要使用的PVC
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

测试

```bash
# 在nginx-test-pvc pod的/usr/share/nginx/html目录写入index.html
echo "add another file to /nfs/data path " >  /usr/share/nginx/html/index.html

# 在nfs服务器/nfs/data目录下查看该文件

# 删除 nginx-test-pvc pod后再查看nfs服务器/nfs/data目录下index.html是否存在

# nginx-test-pvc pod被deployment查询拉取后查看是否能访问index.html
```



















