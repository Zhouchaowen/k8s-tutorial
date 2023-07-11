# Pod

Pod 里的所有容器，共享的是同一个 Network Namespace，并且可以声明共享同一个 Volume。

Pod可以由一个或多个容器组合而成, 属于同一个Pod的多个容器应用之间相互访问时仅需通过localhost就可以通信，使得这一组容器被“绑定”在一个环境中。

## Pod-基础实例

```yaml
# pod.yaml 
apiVersion: v1                  #k8sApi版本 kubectl api-resource 可以获取
kind: Pod                       #资源类型
# 以上为type信息
metadata:                       #元数据:名称，命名空间，标签等
  name: ngx-pod                 #service的名称
  labels:                       #自定义标签属性列表
    env: demo
    owner: zcw
# 以上为元数据部分
spec:                           #期望状态：详细描述pod容器，卷等
  containers:
    - image: nginx:alpine       #镜像地址
      name: ngx                 #镜像名称
      ports:
        - containerPort: 80     #端口
# 以上为资源规格描述部分
```

创建Pod

```bash
kubectl apply -f pod.yaml 
```

查看Pod

```bash
kubectl get pods

NAME      READY   STATUS    RESTARTS   AGE
ngx-pod   1/1     Running   0          70s
```

查看Pod详情

```bash
kubectl describe pod ngx-pod

Name:         ngx-pod
Namespace:    default
Priority:     0
Node:         k8s-node01/10.2.0.102
Start Time:   Tue, 11 Jul 2023 14:57:47 +0800
Labels:       env=demo
              owner=zcw
Annotations:  <none>
Status:       Running
IP:           10.244.1.55
IPs:
  IP:  10.244.1.55
Containers:
  ngx:
    Container ID:   docker://2195b2c550b73a7ac8f1c5ac4e
    Image:          nginx:alpine
    Image ID:       docker-pullable://nginx@sha256:6f94b7f4208b5d5391246c83a
    Port:           80/TCP
    Host Port:      0/TCP
    State:          Running
      Started:      Tue, 11 Jul 2023 14:57:49 +0800
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-8x9sd (ro)
Conditions:
  Type              Status
  Initialized       True 
  Ready             True 
  ContainersReady   True 
  PodScheduled      True 
Volumes:
  kube-api-access-8x9sd:
    Type:                    Projected
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type    Reason     Age   From               Message
  ----    ------     ----  ----               -------
  Normal  Scheduled  97s   default-scheduler  Successfully assigned default/ngx-pod to k8s-node01
  Normal  Pulled     95s   kubelet            Container image "nginx:alpine" already present on machine
  Normal  Created    95s   kubelet            Created container ngx
  Normal  Started    95s   kubelet            Started container ngx
```

查看Pod日志

```bash
kubectl logs ngx-pod [-c containerName]
```

端口转发，访问Pod

```bash
kubectl port-forward ngx-pod 8080:80
```

通过标签选择Pod

```bash
kubectl get pod -l env
```

## Pod-命令行使用

```yaml
# pod-command.yaml 
apiVersion: v1                  #k8sApi版本 kubectl api-resource 可以获取
kind: Pod                       #资源类型
metadata:                       #元数据:名称，命名空间，标签等
  name: pod-command             #service的名称
  labels:                       #自定义标签属性列表
    env: demo
    owner: zcw
spec:                           #期望状态：详细描述pod容器，卷等
  containers:
    - image: nginx:alpine       #镜像地址
      name: ngx                 #镜像名称
      command:                  #执行命令,定义后会覆盖 Dockerfile 中的 CMD,ENTRYPOINT 命令
        - /bin/sh
        - -c
        - "echo hello ngx-container;sleep 3600;"
      ports:
        - containerPort: 80     #端口
```

```bash
kubectl apply -f pod-command.yaml 
```

```bash
# kubectl logs pod-command
hello ngx-container
```

## Pod-环境变量使用

- value

```yaml
# pod-env.yaml
apiVersion: v1                  #k8sApi版本 kubectl api-resource 可以获取
kind: Pod                       #资源类型
metadata:                       #元数据:名称，命名空间，标签等
  name: pod-env                 #service的名称
  labels:                       #自定义标签属性列表
    env: demo
    owner: zcw
spec:                           #期望状态：详细描述pod容器，卷等
  containers:
    - image: nginx:alpine       #镜像地址
      name: ngx                 #镜像名称
      command:                  #执行命令,定义后会覆盖 Dockerfile 中的 CMD,ENTRYPOINT 命令
        - /bin/sh
        - -c
        - "echo NGINX_ENV:$(NGINX_ENV);sleep 3600;"
      ports:
        - containerPort: 80     #端口
      env:                      #定义容器的环境变量
        - name: NGINX_ENV
          value: "nginx_env"
        - name: NGINX_TEST
          value: "nginx_text"
```

```bash
kubectl apply -f pod-env.yaml
```

```bash
# kubectl logs -f ngx-pod
NGINX_ENV:nginx_env
```

进入容器查看环境变量

```bash
# kubectl exec -it ngx-pod -- /bin/sh
/ # env
KUBERNETES_SERVICE_PORT=443
KUBERNETES_PORT=tcp://10.96.0.1:443
HOSTNAME=ngx-pod
SHLVL=1
HOME=/root
PKG_RELEASE=1
NGINX_ENV=nginx_env						# yaml中设置的环境变量NGINX_ENV
TERM=xterm
KUBERNETES_PORT_443_TCP_ADDR=10.96.0.1
NGINX_VERSION=1.23.3
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
KUBERNETES_PORT_443_TCP_PORT=443
NJS_VERSION=0.7.9
NGINX_TEST=nginx_text					# yaml中设置的环境变量NGINX_TEST
KUBERNETES_PORT_443_TCP_PROTO=tcp
KUBERNETES_SERVICE_PORT_HTTPS=443
KUBERNETES_PORT_443_TCP=tcp://10.96.0.1:443
KUBERNETES_SERVICE_HOST=10.96.0.1
PWD=/
```

- valueFrom
  - fieldRef
  - resourceFieldRef
  - configMapKeyRef
  - secretKeyRef

## Pod-多容器

```yaml
# pod-multicontainer.yaml
apiVersion: v1
kind: Pod
metadata:
  name: mc-pod
  labels:
    app: mc-pod
spec:
  volumes:                    # 声明挂载卷 nginx-vol
    - name: nginx-vol
      emptyDir: {}
  containers:
    - name: nginx-container   # 容器一
      image: nginx:alpine
      volumeMounts:           # 声明卷挂载
        - name: nginx-vol
          mountPath: /usr/share/nginx/html
    - name: content-container # 容器二
      image: alpine
      command: ["/bin/sh","-c","while true;do sleep 1; date > /app/index.html;done;"]
      volumeMounts:           # 挂载 nginx-vol 卷
        - name: nginx-vol
          mountPath: /app
```

```bash
kubectl apply -f pod-multicontainer.yaml
```

```bash
kubectl logs mc-pod <nginx-container|content-container>
```

```bash
# kubectl get pod -o wide
NAME     READY   STATUS    RESTARTS   AGE   IP            NODE
mc-pod   2/2     Running   0          42s   10.244.1.62   k8s-node01

# curl 10.244.1.62
Tue Jul 11 11:13:38 UTC 2023
```



## Reference

