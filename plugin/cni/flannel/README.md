## Flannel

Flannel实质上是一种“覆盖网络(overlaynetwork)”，也就是将TCP数据包装在另一种网络包里面进行路由转发和通信，目前已经支持udp、vxlan、host-gw、aws-vpc、gce和alloc路由等数据转发方式，默认的节点间数据通信方式是UDP转发。

## veth-pair+bridge实现多个命名空间相互通信

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
#ip netns exec ns0 ifconfig veth0 10.1.2.1/24
#ip netns exec ns0 ifconfig veth0 up
ip netns exec ns1 ip a a 10.1.0.1/24 dev veth1
ip netns exec ns1 ip l s veth1 up

#ip netns exec ns1 ifconfig veth1 10.1.2.2/24
#ip netns exec ns1 ifconfig veth1 up
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

## 容器夸主机通信



iptables

NAT（网络地址转换）

veth、VLAN、VXLAN、Macvlan

tun/tap

网络接口的混杂模式



```bash
ip netns exec ns1 ethtool -S veth-ns1

ip netns exec ns1 ip link show veth-ns1

# 创建网桥
ip link add name br-xxx type bridge
ip link set br0 up

brctl addbr br0
```



```bash
ip link add veth1a type veth peer name veth1b
ip link add veth2a type veth peer name veth2b

ip netns add ns1
ip netns add ns2

ip link set veth1a netns ns1
ip link set veth2a netns ns2

ip addr add 10.244.1.3 dev veth1b
ip addr add 10.244.2.2 dev veth2b
ip netns exec ns1 ip addr add 10.244.1.2/24 dev veth1a
ip netns exec ns2 ip addr add 10.244.2.3/24 dev veth2a

ip link set veth1b up
ip link set veth2b up
ip netns exec ns1 ip link set veth1a up
ip netns exec ns2 ip link set veth2a up

ip route add 10.244.1.2 dev veth1b
ip route add 10.244.2.3 dev veth2b
ip netns exec ns1 route add -net 10.244.2.0/24 dev veth1a
ip netns exec ns2 route add -net 10.244.1.0/24 dev veth2a

echo 1 > /proc/sys/net/ipv4/conf/veth1b/proxy_arp
echo 1 > /proc/sys/net/ipv4/conf/veth2b/proxy_arp

echo 1 > /proc/sys/net/ipv4/ip_forward


```



```bash
ip netns exec ns1 ping 10.244.2.3
ip netns exec ns2 ping 10.244.1.3
ip netns exec ns1 ip a

# 抓包
ip netns exec ns1 tcpdump -vv -i veth1a

tcpdump -vv -i veth1b

ip netns exec ns2 tcpdump -vv -i veth2a

tcpdump -vv -i veth2b


https://juejin.cn/post/7156502840888786952#heading-21
https://mp.weixin.qq.com/s/CBx-t4-00zGsIpHDY7nuFQ
https://www.koenli.com/fcdddb4a.html
http://arthurchiao.art/blog/cilium-life-of-a-packet-pod-to-service-zh/
https://zhuanlan.zhihu.com/p/594586203


ip netns exec ns1 arp -s 10.244.2.3 02:3a:51:b9:a7:f5 -i veth1a
export GO111MODULE=on
go run -exec sudo main.go bpf_bpfel.go -n veth1b -l veth2b -m 16:dc:ca:ed:68:e8 -i 10.244.2.3


ip netns exec ns2 arp -s 10.244.1.3 a6:27:bb:7f:53:0e -i veth2a
export GO111MODULE=on
go run -exec sudo main.go bpf_bpfel.go -n veth2b -l veth1b -m 2e:e5:91:ea:89:37 -i 10.244.1.3
```



```
7: veth1b@if8: 02:3a:51:b9:a7:f5   8: veth1a@if7: 2e:e5:91:ea:89:37
9: veth2b@if10: a6:27:bb:7f:53:0e  10: veth2a@if9: 16:dc:ca:ed:68:e8 


ip netns exec ns1 tcpdump -s0 -e -i veth1a
tcpdump -s0 -e -i veth1b


tcpdump -s0 -e -i veth2b
```



## tcpdump 参数

`-a`：将网络地址和广播地址转变成名字；

`-d`：将匹配信息包的代码以人们能够理解的汇编格式给出；

`-dd`：将匹配信息包的代码以 c 语言程序段的格式给出；

`-ddd`：将匹配信息包的代码以十进制的形式给出；

`-e`：在输出行打印出数据链路层的头部信息，包括源 mac 和目的 mac，以及网络层的协议；

`-f`：将外部的 Internet 地址以数字的形式打印出来；

`-l`：使标准输出变为缓冲行形式；

`-n`：指定将每个监听到数据包中的域名转换成 IP 地址后显示，不把网络地址转换成名字；

`-nn`：指定将每个监听到的数据包中的域名转换成 IP、端口从应用名称转换成端口号后显示

`-t`：在输出的每一行不打印时间戳；

`-v`：输出一个稍微详细的信息，例如在 ip 包中可以包括 ttl 和服务类型的信息；

`-vv`：输出详细的报文信息；

`-c`：在收到指定的包的数目后，tcpdump 就会停止；

`-F`：从指定的文件中读取表达式,忽略其它的表达式；

`-i`：指定监听的网络接口；

`-p`：将网卡设置为非混杂模式，不能与 host 或 broadcast 一起使用

`-r`：从指定的文件中读取包(这些包一般通过-w 选项产生)；

`-w`：直接将包写入文件中，并不分析和打印出来；

`-s`：snaplen snaplen 表示从一个包中截取的字节数。0 表示包不截断，抓完整的数据包。默认的话 tcpdump 只显示部分数据包,默认 68 字节。

`-T`：将监听到的包直接解释为指定的类型的报文，常见的类型有 rpc （远程过程调用）和 snmp（简单网络管理协议；）

`-X`：告诉 tcpdump 命令，需要把协议头和包内容都原原本本的显示出来（tcpdump 会以 16 进制和 ASCII 的形式显示），这在进行协议分析时是绝对的利器。









路由表：https://blog.csdn.net/kikajack/article/details/80457841



Veth-pair: https://www.cnblogs.com/bakari/p/10613710.html

https://www.jianshu.com/p/369e50201bce

https://cloud.tencent.com/developer/article/1755146

https://weiliang-ms.github.io/wl-awesome/2.%E5%AE%B9%E5%99%A8/%E5%AE%B9%E5%99%A8%E5%8E%9F%E7%90%86/ns-net.html#veth-pair%E4%BB%8B%E7%BB%8D





https://www.zhaohuabing.com/post/2020-03-12-linux-network-virtualization/

https://www.modb.pro/db/50733

https://www.zhihu.com/column/c_1579391544305356800

https://blog.csdn.net/PPPPPPPKD/article/details/121487347

https://zhuanlan.zhihu.com/p/558785823

https://github.com/wenbin8/doc/blob/master/%E5%88%86%E5%B8%83%E5%BC%8F/CloudNative/%E7%BD%91%E7%BB%9C%E5%9F%BA%E7%A1%80/02-Linux%20%E7%BD%91%E7%BB%9C%E5%9F%BA%E7%A1%80%EF%BC%88Network%20Namespase%E3%80%81veth%20pair%E3%80%81bridge%E3%80%81Iptables%EF%BC%89.md



https://www.feiyiblog.com/2020/04/05/flannel%E8%B7%A8%E4%B8%BB%E6%9C%BA%E9%80%9A%E4%BF%A1%E9%83%A8%E7%BD%B2/

https://blog.51cto.com/liujingyu/5251994



解决两个命名空间ping不通的问题

https://gobomb.github.io/post/learning-linux-veth-and-bridge/

https://www.cnblogs.com/bakari/p/10613710.html





详细解释CNI：https://domc.me/2021/10/17/cilium_0_to_0_1/

https://arthurchiao.art/blog/understanding-ebpf-datapath-in-cilium-zh

https://www.infoq.cn/article/FH2Fk1t70wksRuLLUtJf?utm_source=related_read_bottom&utm_medium=article

http://blog.nsfocus.net/cilium-network-0713/

https://www.koenli.com/fcdddb4a.html

https://davidlovezoe.club/wordpress/archives/851

https://morningspace.github.io/lab/

https://morningspace.github.io/tech/k8s-net-mimic-docker/

https://morven.life/posts/create-your-own-cni-with-golang/

https://blog.hdls.me/16131164540989.html



https://www.kancloud.cn/pshizhsysu/network/2202538

https://www.cnblogs.com/goldsunshine/p/10740928.html

https://www.zhihu.com/column/c_1283157687438626816

https://imliuda.com/post/1175

https://www.lixueduan.com/posts/kubernetes/02-cluster-network/

https://www.cnblogs.com/lijin543/p/17352634.html

https://juejin.cn/post/6994825163757846565

http://team.jiunile.com/blog/2020/11/k8s-cilium-service.html





https://blog.51cto.com/liujingyu/5270338



```bash
ip link add veth1a type veth peer name veth1b
ip link add veth2a type veth peer name veth2b

ip netns add ns1
ip netns add ns2

ip link set veth1a netns ns1
ip link set veth2a netns ns2

ip netns exec ns1 ip addr add 10.244.1.2/24 dev veth1a
ip netns exec ns2 ip addr add 10.244.2.3/24 dev veth2a

ip link set veth1b up
ip link set veth2b up
ip netns exec ns1 ip link set veth1a up
ip netns exec ns2 ip link set veth2a up

ip link add veth_net type veth peer name veth_host
ip link set veth_net up
ip link set veth_host up
ip addr add 10.244.0.1 dev veth_host

ip netns exec ns1 route add -net 10.244.3.1/24 dev veth1a
ip netns exec ns1 route add default gw 10.244.0.1 dev veth1a
ip netns exec ns2 route add -net 10.244.0.1/24 dev veth2a
ip netns exec ns2 route add default gw 10.244.0.1 dev veth2a

ip netns exec ns1 arp -s 10.244.0.1 -i veth1a
ip netns exec ns2 arp -s 10.244.0.1 -i veth2a
```

