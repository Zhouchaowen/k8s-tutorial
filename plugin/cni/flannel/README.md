## 创建虚拟网络

### veth-pair

```bash
# 查看网络命名空间列表
ip netns list

# 添加网络命名空间
ip netns a ns0
ip netns a ns1

# 创建一对 veth-pair veth0 veth1
ip l a veth0 type veth peer name veth1

# 将veth0和veth1分别加入两个命名空间
ip l s veth0 netns ns0
ip l s veth1 netns ns1

# 将veth0和veth1分别配上ip
ip netns exec ns0 ip a a 10.1.1.2/24 dev veth0
ip netns exec ns0 ip l s veth0 up

ip netns exec ns1 ip a a 10.1.1.3/24 dev veth1
ip netns exec ns1 ip l s veth1 up

# 从veth0 ping veth1
ip netns exec ns1 ping 10.1.1.3

# 删除网络命名空间
ip netns del ns0
ip netns del ns1
```



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