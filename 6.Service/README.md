# Service

Service 能够将一个接收端口映射到任意的 targetPort。默认情况下，targetPort 将被设置为与 port 字段相同的值。

### iptables模式

当创建 backend Service 时，Kubernetes 会给它指派一个虚拟 IP 地址，比如 10.0.0.1。假设 Service 的端口是 1234，该 Service 会被集群中所有的 kube-proxy 实例观察到。当 kube-proxy 看到一个新的 Service，它会安装一系列的 iptables 规则，从 VIP 重定向到 `per-Service` 规则。 该 `per-Service` 规则连接到 `per-Endpoint` 规则，该 `per-Endpoint` 规则会重定向（目标 `NAT`）到后端的 Pod。

### ipvs模式

在 ipvs 模式下，kube-proxy 监视 Kubernetes 服务和端点，调用 `netlink` 接口相应地创建 IPVS 规则， 并定期将 IPVS 规则与 Kubernetes 服务和端点同步。该控制循环可确保 IPVS 状态与所需状态匹配。访问服务时，IPVS　将流量定向到后端 Pod 之一。



Service服务类型如下：

- ClusterIP：通过集群的内部 IP 暴露服务，选择该值，服务只能够在集群内部可以访问，这也是默认的服务类型。
- NodePort：通过每个 Node 节点上的 IP 和静态端口（NodePort）暴露服务。NodePort 服务会路由到 ClusterIP 服务，这个 ClusterIP 服务会自动创建。通过请求 `NodeIp:NodePort`，可以从集群的外部访问一个 NodePort 服务。
- LoadBalancer：使用云提供商的负载局衡器，可以向外部暴露服务。外部的负载均衡器可以路由到 NodePort 服务和 ClusterIP 服务，这个需要结合具体的云厂商进行操作。
- ExternalName：通过返回 `CNAME` 和它的值，可以将服务映射到 `externalName` 字段的内容（例如， foo.bar.example.com）

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


# nginx-app-service.yaml
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









Headless Service

在某些应用场景中，客户端应用不需要通过Kubernetes内置Service实现的负载均衡功能，或者需要自行完成对服务后端各实例的服务发现机制，或者需要自行实现负载均衡功能，此时可以通过创建一种特殊的名为“Headless”的服务来实现。
Headless Service的概念是这种服务没有入口访问地址（无ClusterIP地址），kube-proxy不会为其创建负载转发规则，而服务名（DNS域名）的解析机制取决于该Headless Service是否设置了Label Selector。

## DNS 配置



```bash
nslookup  nginx-app.default.svc.cluster.local
```



















https://juejin.cn/post/7015109306638942221

https://k8s.easydoc.net/doc/28366845/6GiNOzyZ/C0fakgwO

https://xigang.github.io/2019/07/21/kubernetes-service/

http://www.yunweipai.com/40986.html