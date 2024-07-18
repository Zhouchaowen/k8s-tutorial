docker完了插件

https://www.cnblogs.com/yangmingxianshen/p/10153377.html

https://www.cnblogs.com/xuxinkun/p/5707687.html

http://ninjadq.com/2015/09/29/3rd-party-net-plugin-in-docker



https://www.koenli.com/2ec51ae9.html





## daemon

1. 初始化配置并创建目录
2. 启动前检查最新配置要求
3. 接收绑定配置参数
4. 



## cilium-docker-plugin

```
1.IpamDriver.RequestAddress
2.NetworkDriver.CreateEndpoint
3.NetworkDriver.Join


sudo -E IFACE=ens160 docker-compose up -d --remove-orphans

```





## TC

**TC**全称「**Traffic Control**」，直译过来是「**流量控制**」，在这个领域，你可能更熟悉的是**Linux iptables**或者**netfilter**，它们都能做**packet mangling**，而TC更专注于**packet scheduler**，所谓的网络包调度器，调度网络包的延迟、丢失、传输顺序和速度控制。

TC有4大组件：

- **Queuing disciplines**，简称为**qdisc**，直译是「队列规则」，它的本质是一个带有算法的队列，默认的算法是**FIFO**，形成了一个最简单的流量调度器。
- **Class**，直译是「种类」，它的本质是为上面的qdisc进行分类。因为现实情况下会有很多qdisc存在，每种qdisc有它特殊的职责，根据职责的不同，可以对qdisc进行分类。
- **Filters**，直译是「过滤器」，它是用来过滤传入的网络包，使它们进入到对应class的qdisc中去。
- **Policers**，直译是「规则器」，它其实是filter的跟班，通常会紧跟着filter出现，定义命中filter后网络包的后继操作，如丢弃、延迟或限速。

加载BPF的TC程序：

```bash
# 为目标网卡创建clsact
tc qdisc add dev [network-device] clsact
# 加载bpf程序
tc filter add dev [network-device] <direction> bpf da obj [object-name] sec [section-name]
# 查看
tc filter show dev [network-device] <direction>
```

https://davidlovezoe.club/wordpress/archives/952