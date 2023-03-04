## ConfigMap

configMap 卷提供了向 Pod 注入配置数据的方法。 ConfigMap 对象中存储的数据可以被 `configMap` 类型的卷引用，然后被 Pod 中运行的容器化应用使用。

1. **通过环境变量获取ConfigMap中的内容:**

- 编写yaml

```yaml
# env-config-map.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-app-vars  # 定义名称为cm-app-vars的configMap获取值
data:
  appLogLevel: info     # key=appLogLevel, value=info
  appDataDir: /var/data # key=appDataDir, value=/var/data
---
apiVersion: v1
kind: Pod
metadata:
  name: cm-test-pod
spec:
  containers:
    - name: cm-test-env
      image: busybox:1.28
      command: ["/bin/sh","-c","env | grep APP"]
      env:
        - name: APP_LOG_LEVEL       # 定义环境变量
          valueFrom:
            configMapKeyRef:        # 通过引用configMap获取值
              key: appLogLevel      # 获取cm-appVars中key=appDataDir的值
              name: cm-app-vars
        - name: APP_DATA_DIR
          valueFrom:
            configMapKeyRef:
              key: appDataDir       # 获取cm-appVars中key=appDataDir的值
              name: cm-app-vars
```

- kubectl执行

```bash
root@k8s-master:~/k8s/config-map# kubectl apply -f env-config-map.yaml 
configmap/cm-app-vars created
pod/cm-test-pod created
```

- 查看创建的ConfigMap卷和详细内容

```bash
root@k8s-master:~/k8s/config-map# kubectl get cm
NAME               DATA   AGE
cm-app-vars        2      11s

root@k8s-master:~/k8s/config-map# kubectl describe cm cm-app-vars
Name:         cm-app-vars
Namespace:    default
Labels:       <none>
Annotations:  <none>

Data
====
appDataDir:
----
/var/data
appLogLevel:
----
info

BinaryData
====

Events:  <none>
```

- 查看创建的Pod和打印的Log

```sh
root@k8s-master:~/k8s/config-map# kubectl get pod
NAME          READY   STATUS             RESTARTS     AGE
cm-test-pod   0/1     CrashLoopBackOff   1 (3s ago)   5s

root@k8s-master:~/k8s/config-map# kubectl logs cm-test-pod
APP_DATA_DIR=/var/data
APP_LOG_LEVEL=info
```

- envfrom的形式，执行流程同上

```yaml
# envfrom-config-map.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-app-vars  # 定义名称为cm-app-vars的configMap获取值
data:
  appLogLevel: info     # key=appLogLevel, value=info
  appDataDir: /var/data # key=appDataDir, value=/var/data
---
apiVersion: v1
kind: Pod
metadata:
  name: cm-test-pod
spec:
  containers:
    - name: cm-test
      image: busybox:1.28
      command: [ "/bin/sh", "-c", "env | grep app" ]
      envFrom:
      - configMapRef:
          name: cm-app-vars
  restartPolicy: Never
```

2. **通过Volume挂载的方式将ConfigMap中的内容挂载为容器内部的文件或目:**

- 编写yaml

```yaml
# file-config-map.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-app-config-files
data:
  key-server-xml: |
    <?xml version='1.0' encoding='utf-8'?>
    <Server port="8005" shutdown="SHUTDOWN">
      ......
    </Server>
  key-logging-properties: "logging-properties content"
---
apiVersion: v1
kind: Pod
metadata:
  name: cm-test-app
spec:
  containers:
    - name: cm-test-app
      image: busybox:1.28
      command: [ "/bin/sh","-c","ls /config-files && cat /config-files/server.xml && cat /config-files/logging.properties" ]
      volumeMounts:
        - name: server-xml              # 引用名称为server-xml的volume
          mountPath: /config-files      # 挂载到容器的目录
  volumes:
    - name: server-xml                  # 定义volume的名称
      configMap:                        # volume的值引用configMap
        name: cm-app-config-files       # 引用名称为cm-app-config-files的configMap
        items:
          - key: key-server-xml         # 引用的key
            path: server.xml            # 映射key为key-server-xml的值到文件server.xml
          - key: key-logging-properties # 引用的key
            path: logging.properties    # 映射key为key-logging-properties的值到文件logging.properties
```

- kubectl执行

```bash
root@k8s-master:~/k8s/config-map# kubectl apply -f file-config-map.yaml 
configmap/cm-app-config-files created
pod/cm-test-app created
```

- 查看创建的ConfigMap卷和详细内容

```bash
root@k8s-master:~/k8s/config-map# kubectl get cm
NAME                  DATA   AGE
cm-app-config-files   2      7s

root@k8s-master:~/k8s/config-map# kubectl describe cm cm-app-config-files
Name:         cm-app-config-files
Namespace:    default
Labels:       <none>
Annotations:  <none>

Data
====
key-server-xml:
----
<?xml version='1.0' encoding='utf-8'?>
<Server port="8005" shutdown="SHUTDOWN">
  ......
</Server>

key-logging-properties:
----
logging-properties content

BinaryData
====

Events:  <none>
```

- 查看创建的Pod和打印的Log

```bash
root@k8s-master:~/k8s/config-map# kubectl get pod
NAME          READY   STATUS              RESTARTS   AGE
cm-test-app   0/1     ContainerCreating   0          4s

root@k8s-master:~/k8s/config-map# kubectl logs cm-test-app
logging.properties
server.xml
<?xml version='1.0' encoding='utf-8'?>
<Server port="8005" shutdown="SHUTDOWN">
  ......
</Server>
logging-properties content
```

- 如果在引用ConfigMap时不指定items，则使用volumeMount方式在容器内的目录下为每个item都生成一个文件名为key的文件。

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-app-config-files
data:
  key-server-xml: |
    <?xml version='1.0' encoding='utf-8'?>
    <Server port="8005" shutdown="SHUTDOWN">
      ......
    </Server>
  key-logging-properties: "logging-properties content"
---
apiVersion: v1
kind: Pod
metadata:
  name: cm-test-app
spec:
  containers:
    - name: cm-test-app
      image: busybox:1.28
      command: [ "/bin/sh","-c","ls /config-files && cat /config-files/key-server-xml && cat /config-files/key-logging-properties" ]
      volumeMounts:
        - name: server-xml
          mountPath: /config-files
  volumes:
    - name: server-xml
      configMap:
        name: cm-app-config-files
```

```bash
root@k8s-master:~/k8s/config-map# kubectl logs cm-test-app
key-logging-properties
key-server-xml
<?xml version='1.0' encoding='utf-8'?>
<Server port="8005" shutdown="SHUTDOWN">
  ......
</Server>
logging-properties content
```

3. **configMap的限制**

- ConfigMap必须在Pod之前创建，Pod才能引用它。
- 如果Pod使用envFrom基于ConfigMap定义环境变量，则无效的环境变量名称（例如名称以数字开头）将被忽略，并在事件中被记录为InvalidVariableNames。
- ConfigMap受命名空间限制，只有处于相同命名空间中的Pod才可以引用它。
- ConfigMap无法用于静态Pod。



## Secret

Secret 是一种包含少量敏感信息例如密码、令牌或密钥的对象。Secret 类似于 [ConfigMap](https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/configure-pod-configmap/) 但专门用于保存机密数据。

默认情况下，Kubernetes Secret 未加密地存储在 API 服务器的底层数据存储（etcd）中。 任何拥有 API 访问权限的人都可以检索或修改 Secret，任何有权访问 etcd 的人也可以。 此外，任何有权限在命名空间中创建 Pod 的人都可以使用该访问权限读取该命名空间中的任何 Secret； 这包括间接访问，例如创建 Deployment 的能力。

为了安全地使用 Secret，请至少执行以下步骤：

1. 为 Secret [启用静态加密](https://kubernetes.io/zh-cn/docs/tasks/administer-cluster/encrypt-data/)。
2. 以最小特权访问 Secret 并[启用或配置 RBAC 规则](https://kubernetes.io/zh-cn/docs/reference/access-authn-authz/authorization/)。
3. 限制 Secret 对特定容器的访问。
4. [考虑使用外部 Secret 存储驱动](https://secrets-store-csi-driver.sigs.k8s.io/concepts.html#provider-for-the-secrets-store-csi-driver)。

Pod 可以用三种方式之一来使用 Secret：

- 作为挂载到一个或多个容器上的[卷](https://kubernetes.io/zh-cn/docs/concepts/storage/volumes/) 中的[文件](https://kubernetes.io/zh-cn/docs/concepts/configuration/secret/#using-secrets-as-files-from-a-pod)。
- 作为[容器的环境变量](https://kubernetes.io/zh-cn/docs/concepts/configuration/secret/#using-secrets-as-environment-variables)。
- 由 [kubelet 在为 Pod 拉取镜像时使用](https://kubernetes.io/zh-cn/docs/concepts/configuration/secret/#using-imagepullsecrets)。





