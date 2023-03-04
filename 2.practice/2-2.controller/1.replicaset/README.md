# ReplicaSet

ReplicaSet 的目的是维护一组在任何时候都处于运行状态的 Pod 副本的稳定集合。确保任何时间都有指定数量的 Pod 副本在运行。建议使用 Deployment 而不是直接使用 ReplicaSet， 除非你需要自定义更新业务流程或根本不需要更新。



## 参考

https://medium.com/@suyashmohan/setting-up-postgresql-database-on-kubernetes-24a2a192e962

https://devopscube.com/deploy-postgresql-statefulset/

https://adamtheautomator.com/postgres-to-kubernetes/

https://www.containiq.com/post/deploy-postgres-on-kubernetes