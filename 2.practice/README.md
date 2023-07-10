# 常用命令

```bash
# 非常重要的命令
kubectl api-resource [--namespace=true|false] # 查看k8s所有资源描述
kubectl explain pod.metadata									# 查看关键字相关描述
kubectl run pod-name -o yaml 									# 输出当前pod的yaml文件
```

```bash
# 高频命令
kubectl run pod-name --image=image/xxx # 启动一个pod
kubectl run pod-name --image=image/xxx --dry-run=client -o yaml # client模拟执行并输出yaml

kubectl get all	# 获取所有资源
kubectl get deploy,pod,svc [-n nsName]

kubectl apply -f xxx.yaml 	# 通过yaml创建资源
kubectl diff -f xxx.yaml 		# 对比yaml是否有变化
kubectl delete -f xxx.yaml 	# 通过yaml删除创建的资源


kubectl logs resourceName 	# 获取资源[容器]日志
kubectl describe [deploy|pod|svc] resourceName # 获取资源执行详情
kubectl exec -it podNameXXX /bin/bash	# 进入容器内部
```

## Node

```bash
# 显示集群中当前所有节点的信息
kubectl get nodes

# 获取有关特定节点的详细信息，例如节点的标签、容量和运行状况
kubectl describe node <node-name>：

# 为特定节点添加自定义标签，以便在调度Pod时进行节点选择
kubectl label node <node-name> <label-key>=<label-value>

# 给节点添加污点，以阻止Pod在该节点上被调度
kubectl taint node <node-name> <key>=<value>:<effect>

# 将节点标记为不可调度，并将节点上的Pod迁移到其他可用节点
kubectl drain <node-name>

# 将节点标记为不可调度，防止新的Pod被调度到该节点上
kubectl cordon <node-name>

# 将先前使用 kubectl cordon 命令标记为不可调度的节点恢复为可调度状态
kubectl uncordon <node-name>：

# 获取节点的资源使用情况，如CPU和内存的利用率
kubectl top node


# 从集群中移除一个节点，同时将其上的所有Pod迁移到其他可用节点
kubectl delete node <node-name>
```

## Pod

```bash
kubectl get pods
kubectl get pods --watch
kubectl get pods [-o wide]
kubectl get pod -A

# 按照标签选择器筛选并列出符合条件的 Pod
kubectl get pod -l <label-selector>

# 打印指定 Pod 中特定容器的日志
kubectl logs <pod-name> -c <container-name>

# 显示指定Pod的详细信息
kubectl describe pod <pod-name> 
# 显示集群中所有 Pod 的资源使用情况
kubectl top pod

kubectl exec -it <pod-name> -- /bin/bash

kubectl delete po all 							# 删除当前命名空间的所有pod
kubectl delete pod <pod-name> 			# 删除指定名称的pod
kubectl delete po -l lableName 			# 删除指定标签的pod

# 将本地端口与Pod的端口进行转发，以便直接访问Pod
kubectl port-forward <pod-name> <local-port>:<pod-port>

# 将文件复制到或从 Pod 中的容器
kubectl cp <pod-name>:<source-path> <destination-path>

# 以YAML格式打印指定Pod的详细配置
kubectl get pod <pod-name> -o yaml
```

## ReplicationController

```bash
kubectl get rc

kubectl edit rc rcName

kubectl delete rc

kubectl scale rc rnName --replicas=10 # 通过replicas扩缩容pod数量
```

## ReplicaSet

```bash
kubectl get rs rsName
kubectl describe rs rsName
kubectl delete rs rsName
```

##  DaemonSet

```bash

```

## Deployments

```bash
# 创建一个 Deployment：
kubectl create deployment dName(部署名称) --image=image/xxx [--date|--replicas=3|--port=80]

# 列出所有 Deployments：
kubectl get deployments
kubectl get deployments -A
kubectl get deployments --all-namespace

# 删除一个 Deployment：
kubectl delete deployment <deployment-name>

# 更新一个 Deployment 的镜像：
kubectl set image deployment/<deployment-name> <container-name>=<new-image-name>

# 扩展一个 Deployment 的副本数量：
kubectl scale deployment/<deployment-name> --replicas=<replica-count>

# 滚动更新一个 Deployment 的镜像：
kubectl set image deployment/<deployment-name> <container-name>=<new-image-name> --record
# 暂停一个 Deployment 的滚动更新：
kubectl rollout pause deployment/<deployment-name>

# 恢复一个 Deployment 的滚动更新：
kubectl rollout resume deployment/<deployment-name>

# 回滚一个 Deployment 的滚动更新：
kubectl rollout undo deployment/<deployment-name>

# 查看 Deployment 的状态和历史版本：
kubectl rollout status deployment/<deployment-name>
kubectl rollout history deployment/<deployment-name>

# 从 Deployment 中获取 Pod 的日志：
kubectl logs deployment/<deployment-name>
```

## Job

```bash

```

## Service

```yaml
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
```

```bash
kubectl get service[svc]

kubectl describe svc svcName

kubectl get endpoint serviceName
```

##Namespace 

```yaml
apiVersion: v1
kind: Namespace       # 资源类型
metadata:
  name: kube-example  # 命名空间名称 kube-example
```

```bash
kubectl get ns
kubectl get pod -n namespaceName

kubect delete ns namespaceName # 删除整个命名空间（pod将会伴随命名空间自动删除）
kubectl delete all --all # 第一个all指定正在删除所有资源类型，而--all选项指定将删除所有资源实例
```



## 核心配置

```
/etc/kubernetes
/etc/sysconfig/kubelet
```

