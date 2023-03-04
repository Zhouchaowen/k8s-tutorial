# HPA

Kubernetes中的某个Metrics Server持续采集所有Pod副本的指标数据。HPA控制器通过Metrics Server的API获取这些数据，基于用户定义的扩缩容规则进行计算，得到目标Pod的副本数量。当目标Pod副本数量与当前副本数量不同时，HPA控制器就向Pod的副本控制器（Deployment、RC或ReplicaSet）发起scale操作，调整Pod的副本数量，完成扩缩容操作。

1. 安装metrics-server

```bash
# 国内需要替换镜像名称 并忽略证书- --kubelet-insecure-tls
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# https://zhuanlan.zhihu.com/p/572406293
```

- 注意网络问题
- 注意apiserver访问metrics的ip地址，是集群地址还是内部地址还是外部地址。这需要配置
- 注意apiserver需要开启- --enable-aggregator-routing=true

```bash
# metrics
- --kubelet-preferred-address-types=InternalIP,Hostname,InternalDNS,ExternalDNS,ExternalIP  
- --kubelet-insecure-tls

# apiserver
- --enable-aggregator-routing=true

# https://www.cnblogs.com/fat-girl-spring/p/15936467.html
# https://www.cnblogs.com/shunzi115/p/12438702.html
```

2. 设置HPA

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: php
spec:
  selector:
    matchLabels:
      run: php
  replicas: 1
  template:
    metadata:
      labels:
        run: php
    spec:
      containers:
        - name: php
          image: hpa-php:test
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 80
          resources:
            limits:
              cpu: 400m
            requests:
              cpu: 200m

---
apiVersion: v1
kind: Service
metadata:
  name: php
  labels:
    run: php
spec:
  ports:
    - port: 80
  selector:
    run: php
```

```bash
kubectl autoscale deployment php --cpu-percent=50 --min=1 --max=10

kubectl get hpa

kubectl run -i --tty load-generator --rm --image=busybox:1.28 --restart=Never -- /bin/sh -c "while sleep 0.01; do wget -q -O- http://php-apache; done"

kubectl get hpa php --watch

# 这里启动一个容器，并将无限查询循环发送php服务
kubectl run v1 -it --image=busybox:1.28 /bin/sh
# 登录到容器，执行以下操作
/ # while true; do wget -q -O- http://php.default.svc.cluster.local; done
```

## 参考

https://www.cnblogs.com/yuhaohao/p/14109787.html

https://www.php1.cn/detail/Kubernetes-_RuHe_52dcade4.html