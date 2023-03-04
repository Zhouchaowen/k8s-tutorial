# Service

Service是集群内部的负载均衡机制，用来解决服务发现的问题。kubernetes会分配一个静态IP地址给Service，然后它再去自动管理、维护后面动态变化的Pod集合。

Service 能够将一个接收端口映射到任意的 targetPort。默认情况下，targetPort 将被设置为与 port 字段相同的值。

## Service服务类型

- `ClusterIP`：通过集群的内部 IP 暴露服务，选择该值时服务只能够在集群内部访问。 这也是你没有为服务显式指定 `type` 时使用的默认值。 你可以使用 [Ingress](https://kubernetes.io/zh-cn/docs/concepts/services-networking/ingress/) 或者 [Gateway API](https://gateway-api.sigs.k8s.io/) 向公众暴露服务。
- [`NodePort`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/#type-nodeport)：通过每个节点上的 IP 和静态端口（`NodePort`）暴露服务。 为了让节点端口可用，Kubernetes 设置了集群 IP 地址，这等同于你请求 `type: ClusterIP` 的服务。
- [`LoadBalancer`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/#loadbalancer)：使用云提供商的负载均衡器向外部暴露服务。 外部负载均衡器可以将流量路由到自动创建的 `NodePort` 服务和 `ClusterIP` 服务上。
- [`ExternalName`](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/#externalname)：通过返回 `CNAME` 记录和对应值，可以将服务映射到 `externalName` 字段的内容（例如，`foo.bar.example.com`）。 无需创建任何类型代理。

### NodePort 类型

```yaml
# nginx-service-deployment.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name:  nginx-deploy
  namespace: default
spec:
  replicas: 3     # 期望的 Pod 副本数量，默认值为1
  selector:       # Label Selector，必须匹配 Pod 模板中的标签
    matchLabels:
      app: nginx-app
  template:  # Pod 模板
    metadata:
      labels:
        app: nginx-app
    spec:
      containers:
        - name: nginx
          image: nginx:1.7.9
          ports:
            - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-app
spec:
  type: NodePort
  ports:
    - port: 8080
      targetPort: 80
      nodePort: 32380
  selector:
    app: nginx-app
```

验证

```bash
kubectl apply -f nginx-service-deployment.yaml

# 进入容器修改index.html
kubectl exec -it nginx-deploy-ContainerXXX -- /bin/bash
echo "k8s-nodeXX" > /usr/share/nginx/html/index.html

# 如何访问
# 1.通过pod ip 访问 curl podIp:containerPort
curl 10.244.1.9

# 2.通过cluster ip 访问 curl clusterIp:port
curl 10.101.85.130:8080

# 3.通过node ip 访问 curl nodeIp:port
curl 10.2.0.101:32380
```

### ExternalName

```yaml
kind: Service
apiVersion: v1
metadata:
  name: my-service
  namespace: prod
spec:
  type: ExternalName
  externalName: my.database.example.com
```

当访问地址 `my-service.prod.svc.cluster.local`（后面服务发现的时候我们会再深入讲解）时，集群的 DNS 服务将返回一个值为 my.database.example.com 的 `CNAME` 记录。

### Headless Service

有时不需要或不想要负载均衡，以及单独的 Service IP。 遇到这种情况，可以通过指定 Cluster IP（`spec.clusterIP`）的值为 `"None"` 来创建 `Headless` Service。

Headless Service的概念是这种服务没有入口访问地址（无ClusterIP地址），kube-proxy不会为其创建负载转发规则，而服务名（DNS域名）的解析机制取决于该Headless Service是否设置了Label Selector。

## DNS 配置

- 普通的 Service：会生成 `servicename.namespace.svc.cluster.local` 的域名，会解析到 Service 对应的 ClusterIP 上，在 Pod 之间的调用可以简写成 `servicename.namespace`，如果处于同一个命名空间下面，甚至可以只写成 `servicename` 即可访问
- Headless Service：无头服务，就是把 clusterIP 设置为 None 的，会被解析为指定 Pod 的 IP 列表，同样还可以通过 `podname.servicename.namespace.svc.cluster.local` 访问到具体的某一个 Pod。

```bash
nslookup  nginx-app.default.svc.cluster.local
```



## iptables模式

当创建 backend Service 时，Kubernetes 会给它指派一个虚拟 IP 地址，比如 10.0.0.1。假设 Service 的端口是 1234，该 Service 会被集群中所有的 kube-proxy 实例观察到。当 kube-proxy 看到一个新的 Service，它会安装一系列的 iptables 规则，从 VIP 重定向到 `per-Service` 规则。 该 `per-Service` 规则连接到 `per-Endpoint` 规则，该 `per-Endpoint` 规则会重定向（目标 `NAT`）到后端的 Pod。

## ipvs模式

在 ipvs 模式下，kube-proxy 监视 Kubernetes 服务和端点，调用 `netlink` 接口相应地创建 IPVS 规则， 并定期将 IPVS 规则与 Kubernetes 服务和端点同步。该控制循环可确保 IPVS 状态与所需状态匹配。访问服务时，IPVS　将流量定向到后端 Pod 之一。



## Ingress

对于小规模的应用我们使用 NodePort 或许能够满足我们的需求，但是当你的应用越来越多的时候，你就会发现对于 NodePort 的管理就非常麻烦了，这个时候使用 Ingress 就非常方便了，可以避免管理大量的端口。

Ingress 其实就是从 Kuberenets 集群外部访问集群的一个入口，将外部的请求转发到集群内不同的 Service 上，其实就相当于 nginx、haproxy 等负载均衡代理服务器。

Ingress Controller 可以理解为一个监听器，通过不断地监听 kube-apiserver，实时的感知后端 Service、Pod 的变化，当得到这些信息变化后，Ingress Controller 再结合 Ingress 的配置，更新反向代理负载均衡器，达到服务发现的作用。



安装ingress-nginx:

```bash
https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.1.1/deploy/static/provider/cloud/deploy.yaml

# 1.需要修改images为国内镜像
# 2.
```



```yaml
# nginx1-7.yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx1-7
spec:
  ports:
    - port: 80
      targetPort: 80
  selector:
    app: nginx1-7
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx1-7-deployment
spec:
  replicas: 2
  selector:       # Label Selector，必须匹配 Pod 模板中的标签
    matchLabels:
      app: nginx1-7
  template:
    metadata:
      labels:
        app: nginx1-7
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

```yaml
# nginx1-9.yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx1-9
spec:
  ports:
    - port: 80
      targetPort: 80
  selector:
    app: nginx1-9
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx1-9-deployment
spec:
  replicas: 2
  selector:       # Label Selector，必须匹配 Pod 模板中的标签
    matchLabels:
      app: nginx1-9
  template:
    metadata:
      labels:
        app: nginx1-9
    spec:
      containers:
      - name: nginx
        image: nginx:1.9.1
        ports:
        - containerPort: 80
```

```yaml
# test-ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /$2
spec:
  ingressClassName: nginx
  rules:
    - host: nginx.my.com
      http:
        paths:
          - pathType: Prefix
            path: "/v17"
            backend:
              service:
                name: nginx1-7
                port:
                  number: 80
          - pathType: Prefix
            path: "/v19"
            backend:
              service:
                name: nginx1-9
                port:
                  number: 80
```

```bash
kubectl exec -it nginx-deploy-ContainerXXX -- /bin/bash
echo "k8s-nodeXX" > /usr/share/nginx/html/index.html

# path路径问题 https://coding.imooc.com/learn/questiondetail/152694.html

curl nginx.my.com:31691/v17
curl nginx.my.com:31691/v19
```

https://www.cnblogs.com/askajohnny/p/16160721.html



https://juejin.cn/post/7015109306638942221

https://k8s.easydoc.net/doc/28366845/6GiNOzyZ/C0fakgwO

https://xigang.github.io/2019/07/21/kubernetes-service/

http://www.yunweipai.com/40986.html