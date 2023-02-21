## Volume

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

## ConfigMap

通过环境变量获取ConfigMap中的内容。通过Volume挂载的方式将ConfigMap中的内容挂载为容器内部的文件或目录。

环境变量的名称受POSIX命名规范（[a-zA-Z_\]\[a-zA-Z0-9_]*）约束，不能以数字开头

```yaml
# env-config-map.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-appVars  # 定义名称为cm-appVars的configMap获取值
data:
  appLogLevel: info     # key=appLogLevel, value=info
  appDataDir: /var/data # key=appDataDir, value=/var/data
```

使用

```yaml
# use-env-config-map.yaml
apiVersion: v1
kind: Pod
metadata:
  name: cm-test-pod
spec:
  containers:
    - name: cm-test
      image: busybox
      command: ["/bin/sh","-c","env | grep APP"]
      env:
        - name: APP_LOG_LEVEL       # 定义环境变量
          valueFrom:
            configMapKeyRef:        # 通过引用configMap获取值
              key: appLogLevel      # 获取cm-appVars中key=appDataDir的值
              name: cm-appVars
        - name: APP_DATA_DIR
          valueFrom:
            configMapKeyRef:
              key: appDataDir       # 获取cm-appVars中key=appDataDir的值
              name: cm-appVars
```

查看启动时打印的环境变量

```sh
kubectl create -f base-config-map.yaml
kubectl create -f use-base-config-map.yaml

kubectl logs cm-test-pod
```

file-config-map

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
      <Listener className="org.apache.catalina.startup.VersionLoggerListener" />
      <Listener className="org.apache.catalina.core.AprLifecycleListener" SSLEngine="on" />
      <Listener className="org.apache.catalina.core.JreMemoryLeakPreventionListener" />
      <Listener className="org.apache.catalina.mbeans.GlobalResourcesLifecycleListener" />
      <Listener className="org.apache.catalina.core.ThreadLocalLeakPreventionListener" />
      <GlobalNamingResources>
        <Resource name="UserDatabase" auth="Container"
                  type="org.apache.catalina.UserDatabase"
                  description="User database that can be updated and saved"
                  factory="org.apache.catalina.users.MemoryUserDatabaseFactory"
                  pathname="conf/tomcat-users.xml" />
      </GlobalNamingResources>

      <Service name="Catalina">
        <Connector port="8080" protocol="HTTP/1.1"
                   connectionTimeout="20000"
                   redirectPort="8443" />
        <Connector port="8009" protocol="AJP/1.3" redirectPort="8443" />
        <Engine name="Catalina" defaultHost="localhost">
          <Realm className="org.apache.catalina.realm.LockOutRealm">
            <Realm className="org.apache.catalina.realm.UserDatabaseRealm"
                   resourceName="UserDatabase"/>
          </Realm>
          <Host name="localhost"  appBase="webapps"
                unpackWARs="true" autoDeploy="true">
            <Valve className="org.apache.catalina.valves.AccessLogValve" directory="logs"
                   prefix="localhost_access_log" suffix=".txt"
                   pattern="%h %l %u %t &quot;%r&quot; %s %b" />

          </Host>
        </Engine>
      </Service>
    </Server>
  key-logging-properties: "handlers
    = 1catalina.org.apache.juli.FileHandler, 2localhost.org.apache.juli.FileHandler,
    3manager.org.apache.juli.FileHandler, 4host-manager.org.apache.juli.FileHandler,
    java.util.logging.ConsoleHandler\r\n\r\n.handlers = 1catalina.org.apache.juli.FileHandler,
    java.util.logging.ConsoleHandler\r\n\r\n1catalina.org.apache.juli.FileHandler.level
    = FINE\r\n1catalina.org.apache.juli.FileHandler.directory = ${catalina.base}/logs\r\n1catalina.org.apache.juli.FileHandler.prefix
    = catalina.\r\n\r\n2localhost.org.apache.juli.FileHandler.level = FINE\r\n2localhost.org.apache.juli.FileHandler.directory
    = ${catalina.base}/logs\r\n2localhost.org.apache.juli.FileHandler.prefix = localhost.\r\n\r\n3manager.org.apache.juli.FileHandler.level
    = FINE\r\n3manager.org.apache.juli.FileHandler.directory = ${catalina.base}/logs\r\n3manager.org.apache.juli.FileHandler.prefix
    = manager.\r\n\r\n4host-manager.org.apache.juli.FileHandler.level = FINE\r\n4host-manager.org.apache.juli.FileHandler.directory
    = ${catalina.base}/logs\r\n4host-manager.org.apache.juli.FileHandler.prefix =
    host-manager.\r\n\r\njava.util.logging.ConsoleHandler.level = FINE\r\njava.util.logging.ConsoleHandler.formatter
    = java.util.logging.SimpleFormatter\r\n\r\n\r\norg.apache.catalina.core.ContainerBase.[Catalina].[localhost].level
    = INFO\r\norg.apache.catalina.core.ContainerBase.[Catalina].[localhost].handlers
    = 2localhost.org.apache.juli.FileHandler\r\n\r\norg.apache.catalina.core.ContainerBase.[Catalina].[localhost].[/manager].level
    = INFO\r\norg.apache.catalina.core.ContainerBase.[Catalina].[localhost].[/manager].handlers
    = 3manager.org.apache.juli.FileHandler\r\n\r\norg.apache.catalina.core.ContainerBase.[Catalina].[localhost].[/host-manager].level
    = INFO\r\norg.apache.catalina.core.ContainerBase.[Catalina].[localhost].[/host-manager].handlers
    = 4host-manager.org.apache.juli.FileHandler\r\n\r\n"
```

Pod使用file-config-map

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: cm-test-app
spec:
  containers:
    - name: cm-test-app
      image: kubeguide/tomcat-app:v1
      ports:
        - containerPort: 8080
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

创建

```bash
kubectl create -f file-config-map.yaml
kubectl create -f use-file-config-map.yaml

kubectl exec -it cm-test-app -- bash
```

如果在引用ConfigMap时不指定items，则使用volumeMount方式在容器内的目录下为每个item都生成一个文件名为key的文件。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: cm-test-app
spec:
  containers:
    - name: cm-test-app
      image: kubeguide/tomcat-app:v1
      imagePullPolicy: Never
      ports:
        - containerPort: 8080
      volumeMounts:
        - name: server-xml
          mountPath: /config-files
  volumes:
    - name: server-xml
      configMap:
        name: cm-app-config-files
```

configMap的限制

- ConfigMap必须在Pod之前创建，Pod才能引用它。
- 如果Pod使用envFrom基于ConfigMap定义环境变量，则无效的环境变量名称（例如名称以数字开头）将被忽略，并在事件中被记录为InvalidVariableNames。
- ConfigMap受命名空间限制，只有处于相同命名空间中的Pod才可以引用它。
- ConfigMap无法用于静态Pod。

## Downward API

- 环境变量：将Pod或Container信息设置为容器内的环境变量。
- Volume挂载：将Pod或Container信息以文件的形式挂载到容器内部。

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

