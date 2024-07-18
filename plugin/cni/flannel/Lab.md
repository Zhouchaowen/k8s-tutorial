# 实验一

```bash
# 创建一对 veth-pair 
ip link add veth1a type veth peer name veth1b

# 分别设置 ip
ip addr add 10.244.1.2 dev veth1a
ip addr add 10.244.1.3 dev veth1b

# 启动 veth 和 pair
ip link set veth1a up
ip link set veth1b up
```

```bash
# 从veth1a ping 10.244.1.3
ping -I veth1a 10.244.1.3

# 在veth1a抓包
tcpdump -s0 -e -v -n -i veth1a
```

由于 veth1a 和 veth1b 处于同一个网段，且是第一次连接，所以会事先发 ARP 包，但 veth1b 并没有响应 ARP 包。这是由于我使用的 Ubuntu 系统内核中一些 ARP 相关的默认配置限制所导致的，需要修改一下配置项：

```bash
# 设置是否允许接收从本机IP地址上发送给本机的数据包
echo 1 > /proc/sys/net/ipv4/conf/veth1/accept_local
echo 1 > /proc/sys/net/ipv4/conf/veth2/accept_local
# 关闭反向路由监测
echo 0 > /proc/sys/net/ipv4/conf/all/rp_filter
echo 0 > /proc/sys/net/ipv4/conf/veth1/rp_filter
echo 0 > /proc/sys/net/ipv4/conf/veth2/rp_filter
```

https://ctimbai.github.io/2019/03/03/tech/net/vnet/veth-pair%E8%AF%A6%E8%A7%A3/

https://zhuanlan.zhihu.com/p/586504492

# 实验二

veth-pair实现一个命名空间和本机host通信

```bash
# 创建veth pair
ip link add veth1 type veth peer name veth2

# 创建ns1网络命名空间
ip netns add ns1

# 将veth2移动到ns1空间下
ip link set veth2 netns ns1

# 分别给veth1, veth2设置IP
ip addr add 10.244.1.2/24 dev veth1
ip netns exec ns1 ip addr add 10.244.1.3/24 dev veth2

ip netns exec ns1 ip a

# 分别设置veth1,veth2为启动状态
ip link set veth1 up
ip netns exec ns1 ip link set veth2 up
```

```bash
# 在主网络命名空间下，去ping ns1空间里的veth2网卡
ping 10.244.1.3 -I veth1
```

在ns1网络命名空间里，测试是否能ping通主网络命名空间里的eth0网卡?

```bash
# 配置路由
ip netns exec ns1 route add -net 10.211.55.0/24 dev veth2
```

```bash
# 在ns1网络命名空间里去ping 宿主机网卡
ip netns exec ns1 ping 10.211.55.122
```

https://zhuanlan.zhihu.com/p/587261667

# 实验三

veth-pair实现一个命名空间和其它机器host通信

```bash
ip netns add ns1

ip link add veth1 type veth peer name veth2

ip link set veth2 netns ns1

ip addr add 10.244.1.2/24 dev veth1
ip link set veth1 up

ip netns exec ns1 ip addr add 10.244.1.3/24 dev veth2
ip netns exec ns1 ip link set veth2 up

ip netns exec ns1 route add default gw 10.244.1.2
```

```bash
ip netns exec ns1 ping 10.211.55.123

# 设置MASQUERADE转换
iptables -t nat -A POSTROUTING -s 10.244.1.0/24 -o eth0 -j MASQUERADE
```

https://zhuanlan.zhihu.com/p/589345014

# 实验四

veth-pair实现两个命名空间通信

```bash
# 查看网络命名空间列表
ip netns list

# 添加网络命名空间
ip netns a ns1a
ip netns a ns1b

# 创建一对 veth-pair veth1a veth1b
ip l a veth1a type veth peer name veth1b

# 将 veth1a 和 veth1b 分别加入两个命名空间
ip l s veth1a netns ns1a
ip l s veth1b netns ns1b

# 将veth0和veth1分别配上ip
ip netns exec ns1a ip a a 10.244.1.2/24 dev veth1a
ip netns exec ns1a ip l s veth1a up

ip netns exec ns1b ip a a 10.244.1.3/24 dev veth1b
ip netns exec ns1b ip l s veth1b up

# 从veth0 ping veth1
ip netns exec ns1a ping 10.244.1.3

# 删除网络命名空间
ip netns del ns1a
ip netns del ns1b
```

# 实验五

不同网段命名空间通信（同一主机，把主机当路由器）

方式一：

```bash
ip link add veth1a type veth peer name veth1b
ip link add veth2a type veth peer name veth2b

ip netns add ns1
ip netns add ns2

ip link set veth1a netns ns1
ip link set veth2a netns ns2

# 注意这里配置的是/32，代表是一个ip地址，所以不会创建默认路由，需要手动添加
ip netns exec ns1 ip addr add 10.244.1.3/32 dev veth1a
ip addr add 10.244.1.2 dev veth1b

ip netns exec ns2 ip addr add 10.244.2.3/32 dev veth2a
ip addr add 10.244.2.2 dev veth2b

ip netns exec ns1 ip link set veth1a up
ip link set veth1b up

ip netns exec ns2 ip link set veth2a up
ip link set veth2b up

# 主机路由：1.目的地址为 10.244.1.3 走 veth1b 2.目的地址为 10.244.2.3 走 veth2b
ip route add 10.244.1.3 dev veth1b
ip route add 10.244.2.3 dev veth2b

# ns路由：10.244.2.0/24 ping出去的时候需要目的地ip的 mac 地址 
# 10.244.1.0/24 回复 ICMP 时需要使用
ip netns exec ns1 route add -net 10.244.1.0/24 dev veth1a
ip netns exec ns1 route add -net 10.244.2.0/24 dev veth1a

ip netns exec ns2 route add -net 10.244.1.0/24 dev veth2a
ip netns exec ns2 route add -net 10.244.2.0/24 dev veth2a
```

```bash
ip netns exec ns1 ping 10.244.2.3
ip netns exec ns2 ping 10.244.1.3
```

```bash
echo 1 > /proc/sys/net/ipv4/conf/veth1b/proxy_arp
echo 1 > /proc/sys/net/ipv4/conf/veth2b/proxy_arp
echo 1 > /proc/sys/net/ipv4/ip_forward
```

https://vinchin.com/blog/vinchin-technique-share-details.html?id=4525

方式二：

```bash
ip link add veth1a type veth peer name veth1b
ip link add veth2a type veth peer name veth2b

ip netns add ns1
ip netns add ns2

ip link set veth1a netns ns1
ip link set veth2a netns ns2

ip netns exec ns1 ip addr add 10.244.1.3/32 dev veth1a
ip addr add 10.244.1.2 dev veth1b

ip netns exec ns2 ip addr add 10.244.2.3/32 dev veth2a
ip addr add 10.244.2.2 dev veth2b

ip netns exec ns1 ip link set veth1a up
ip link set veth1b up

ip netns exec ns2 ip link set veth2a up
ip link set veth2b up

# 主机路由：1.目的地址为 10.244.1.3 走 veth1b 2.目的地址为 10.244.2.3 走 veth2b
ip route add 10.244.1.3 dev veth1b
ip route add 10.244.2.3 dev veth2b

ip netns exec ns1 route add -net 10.244.1.0/24 dev veth1a
ip netns exec ns1 route add default gw 10.244.1.2 dev veth1a

ip netns exec ns2 route add -net 10.244.2.0/24 dev veth2a
ip netns exec ns2 route add default gw 10.244.2.2 dev veth2a

echo 1 > /proc/sys/net/ipv4/ip_forward
```

将宿主机作为路由器来使用，不用开启proxy_arp代理转发

方式三：veth-pair+bridge实现多个命名空间相互通信

```bash
brctl show

# 添加网络命名空间
ip netns a ns1
ip netns a ns2

# 创建bridge br0并启动
ip l a br0 type bridge
ip l s br0 up

# 创建两对 veth-pair
ip l a veth1 type veth peer name br-veth1
ip l a veth2 type veth peer name br-veth2

# 分别将两对 veth-pair 加入两个ns和br0
ip l s veth1 netns ns1
ip l s br-veth1 master br0
ip l s br-veth1 up

ip l s veth2 netns ns2
ip l s br-veth2 master br0
ip l s br-veth2 up

# 给两个ns中的veth配置ip并启用
ip netns exec ns1 ip a a 10.1.0.1/24 dev veth1
ip netns exec ns1 ip l s veth1 up

ip netns exec ns2 ip a a 10.1.0.2/24 dev veth2
ip netns exec ns2 ip l s veth2 up


ip netns exec ns1 tcpdump -i veth1

# ifconfig br0 down
# brctl delbr br0

# 查看bridge br0上的端口：
brctl show 

# 查看对应的MAC地址和不同ns中的MAC。
brctl showmacs br
```

```bash
# 查看arp
arp -n

# 删除一跳arp
arp -d xx.xx.xx.xx
```

# 实验六

ebpf tc 实现加速

```bash
ip link add veth1a type veth peer name veth1b
ip link add veth2a type veth peer name veth2b

ip netns add ns1
ip netns add ns2

ip link set veth1a netns ns1
ip link set veth2a netns ns2

ip netns exec ns1 ip addr add 10.244.1.3/24 dev veth1a
ip addr add 10.244.1.2 dev veth1b

ip netns exec ns2 ip addr add 10.244.2.3/24 dev veth2a
ip addr add 10.244.2.2 dev veth2b

ip netns exec ns1 ip link set veth1a up
ip link set veth1b up

ip netns exec ns2 ip link set veth2a up
ip link set veth2b up

ip netns exec ns1 route add -net 10.244.2.0/24 dev veth1a
ip netns exec ns2 route add -net 10.244.1.0/24 dev veth2a

echo 1 > /proc/sys/net/ipv4/ip_forward

```

```bash
ip netns exec ns1 arp -s 10.244.2.3 a6:27:bb:7f:53:0e -i veth1a
ip netns exec ns2 arp -s 10.244.1.3 a6:27:bb:7f:53:0e -i veth2a

ip netns exec ns1 arp -s 10.244.2.3 a6:27:bb:7f:53:0e -i veth1a
ip netns exec ns2 arp -s 10.244.1.3 a6:27:bb:7f:53:0e -i veth2a
```

```bash
ip netns exec ns1 tcpdump -e -n -i veth1a

tcpdump -e -n -i veth1b

ip netns exec ns2 tcpdump -e -n -i veth2a

tcpdump -e -n -i veth2b
```