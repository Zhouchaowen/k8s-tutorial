# kube-scheduler

主要作用就是根据特定的调度算法和调度策略将 Pod 调度到合适的 Node 节点上去，是一个独立的二进制程序，启动之后会一直监听 API Server，获取到 `PodSpec.NodeName` 为空的 Pod，对每个 Pod 都会创建一个 binding

```bash
kubectl apply -> pod -> etcd 

kube-scheduler->watch apiserver -> pod -> schedule -> filter(nodes) -> score(优先级) -> node(bind) -> update pod -> save etcd

Kubelet -> watch apiserver -> pod(pull info) -> start container ->pod status -> update etcd
```

详细的流程是这样的：

- 首先，客户端通过 API Server 的 REST API 或者 kubectl 工具创建 Pod 资源
- API Server 收到用户请求后，存储相关数据到 etcd 数据库中
- 调度器监听 API Server 查看到还未被调度(bind)的 Pod 列表，循环遍历地为每个 Pod 尝试分配节点，这个分配过程就是我们上面提到的两个阶段：
  - 预选阶段(Predicates)，过滤节点，调度器用一组规则过滤掉不符合要求的 Node 节点，比如 Pod 设置了资源的 request，那么可用资源比 Pod 需要的资源少的主机显然就会被过滤掉
  - 优选阶段(Priorities)，为节点的优先级打分，将上一阶段过滤出来的 Node 列表进行打分，调度器会考虑一些整体的优化策略，比如把 Deployment 控制的多个 Pod 副本尽量分布到不同的主机上，使用最低负载的主机等等策略
- 经过上面的阶段过滤后选择打分最高的 Node 节点和 Pod 进行 `binding` 操作，然后将结果存储到 etcd 中 最后被选择出来的 Node 节点对应的 kubelet 去执行创建 Pod 的相关操作（当然也是 watch APIServer 发现的）。









## Reference

https://www.qikqiak.com/k8strain2/scheduler/overview/