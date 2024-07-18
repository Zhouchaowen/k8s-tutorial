## 知识点

- 命名空间
- 路由表
- arp表
- veth-pair
- 网桥
- nat转换
- ip-ip
- Overlay Network
  - UDP
  - VXLAN
  - host-gw



## 实战

1. 命名空间与外界通信
2. 同主机命名空间与命名空间通信
3. 不同主机命名空间与命名空间通信



## 路由

在一台linux机器上要访问一个目标ip时四步口诀：

1. 如果本机有目标ip，则会直接访问本地; 如果本地没有目标ip，则看第2步
2. 用route -n查看路由，如果路由条目里包含了目标ip的网段，则数据包就会从对应路由条目后面的网卡出去
3. 如果没有对应网段的路由条目，则全部都走网关
4. 如果网关也没有，则报错：网络不可达

**注意:** 当不能直接到达目标ip, 那么每到达一个机器都会重复上面四步，直到找到目标

网关只能加路由条目里已有的路由网段里的一个IP (ping不通此IP都可以） 加网关不需要指定子网掩码

```bash
# 临时配置与删除(立即生效,重启网络服务就没了)
route add default gw x.x.x.x # route del default gw x.x.x.x
```

**ip_forward**: linux内核里的一个参数.当两边机器不同网段IP通过中间双网卡机器进行路由交互时,需要将此参数值改为1,也就是打开ip_forward。



路由策略

```bash
ip rule list
```

查看路由表

```bash
# 查看路由表 or route -n
ip route list table main
```

添加路由表

```bash
ip route add default via 192.168.2.1 dev eth0
```



```
分布式
rdeis
mongo
docker
网络协议
psql
mysql
grpc
gin

```



