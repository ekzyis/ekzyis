Title:        WireGuard Port Forwarding
Date:         2023-09-25
ReadingTime:  10 minutes
Sats:         1803
Comments:     https://stacker.news/items/265524

---

# introduction

Today, we will implement port forwarding with [`iptables`](https://wiki.archlinux.org/title/iptables). We will do this to expose a service to the public internet which is running inside a VPN.

The exposed service will be a basic HTTP server running on a host in a home network. This means that the host itself already has a connection to the public internet (as most hosts in a home network do) and is not only connected to all other hosts inside the VPN via a VPN router.

However, since the HTTP server is running inside the VPN, we will need to use port forwarding inside our VPN instead of on the router with which we access the public internet. Port forwarding on that home router won't work since it can't access our VPN service (even though the host it's running on is inside the home network!).

Essentially, we will update the network topology from the [previous blog post](/blog/20230821-wireguard-packet-forwarding) with a connection to the public internet:

![](/blog/img/star_network_with_internet_access.png)

As always, after reading my blog post, you will not only be able to setup your own HTTP server inside a VPN; but will actually understand what happens at the network level. With this knowledge, you will be able to adapt and find the best solution for your specific needs.

Therefore, we start with a primer about NAT which stands for _Network Address Translation_.

---

# NAT primer

Network Address Translation was developed to mitigate the issue of not having enough IPv4 addresses for every device which wants to connect to the internet. IPv4 addresses are only 32 bits long which means there are only 2<sup>32</sup> = 4,294,967,296 addresses available. This may sound like a lot but when considering that more and more devices are connected to the internet (which isn't going to stop anytime soon) and that the world population is around 8 billion people, there simply aren't enough IPv4 addresses to go around. Hence a solution needed to be found.

With NAT, one can hide multiple devices behind a router and thus behind a single IPv4 address. The hosts connected to the internet via this router are assigned private IPv4 addresses and thus form a private network. Private IPv4 addresses are addresses which are not routable on the public internet. Routers on the public internet would simply not forward packets containing these IP addresses. The following IPv4 address ranges were reserved for such private use by [IANA](https://www.iana.org/) in [RFC1918](https://datatracker.ietf.org/doc/html/rfc1918#section-3):

- 10.0.0.8/8
- 172.16.0.0/12
- 192.168.0.0/16

This mitigates the problem since these private IP addresses only have to be unique per private network and not globally. That's also the reason why these addresses are not routable. How should a router know how to route these packets? There isn't a unique device with this IP address!

To enable internet access for hosts with only a private IPv4 address, every packet with a destination not inside the private network gets forwarded to the _NAT gateway_. Usually, this is the router. The router then replaces the private source IP address with its own public IP address <span id="ft-0b">[[0]](#ft-0)</span> and the port with another random port:

![](/blog/img/nat_send.png)

The NAT gateway then stores this replacement in a NAT table:

![](/blog/img/nat_table.png)

For arriving packets, this table is consulted to reverse the translation:

![](/blog/img/nat_recv.png)

The NAT IP address and NAT port are required to reverse this process.
Without them, we would not be able to distinguish multiple connections from the private network to the same destination IP address.

This method to allow multiple hosts inside a private network access to the internet is called _Source NAT_ (SNAT) since we change the source IP address when initiating the connection.

If we want to allow hosts from the public internet to access hosts in a private network, we use _Destination NAT_ (DNAT) since we change the destination IP address when the connection is initiated.

The reversal in both methods is automatically handled using the NAT table.

---

# Port forwarding

We will now apply this knowledge to expose an HTTP server running inside a VPN. Our setup will work like this:

![](/blog/img/nat_http.png)

We will have to use DNAT _and_ SNAT since we need to change the destination IP address to 10.172.16.2 _and_ the source IP address to 10.172.16.1 since we can only route internal IP addresses (10.172.16.0/24) over the virtual network interface to the HTTP server.

## initial configuration

We will start with the following configuration of the VPN router which will also act as a NAT gateway:

_/etc/wireguard/wg0.conf @ 10.172.16.1:_
```
[Interface]
Address = 10.172.16.1/32
PrivateKey = r3M073+s3cR37+fouaQZbP5QqfgwypHjKGBNmztxNEc=
ListenPort = 51913

[Peer]
AllowedIPs = 10.172.16.2/32
PublicKey = /wH4OzafBUJVvRGzK8itUweV/GpwoUzn7OS99lr7gHI=
```

_firewall configuration @ 10.172.16.1:_
```
-P INPUT DROP
-P FORWARD DROP
-P OUTPUT DROP
-A INPUT -p tcp -m tcp --dport 22 -j ACCEPT
-A INPUT -i wg0 -j ACCEPT
-A INPUT -i eth0 -p udp -m udp --dport 51913 -j ACCEPT
-A FORWARD -i wg0 -o wg0 -j ACCEPT
-A OUTPUT -m state --state ESTABLISHED -j ACCEPT
-A OUTPUT -o wg0 -j ACCEPT
```

IP forwarding is also enabled in the kernel:

```
$ sysctl net.ipv4.ip_forward
net.ipv4.ip_forward = 1
```

The host on which we will run the HTTP server is configured like this:

_/etc/wireguard/wg0.conf @ 10.172.16.2:_
```
[Interface]
Address = 10.172.16.2/32
PrivateKey = l0c4l+s3cR37+RDr+dJdgX/ACeRQLANiduQRJK9O23A=

[Peer]
AllowedIPs = 10.172.16.0/24
PublicKey = GL33DRrI8/2yAT6+r5mTtBLd7CoErAAsio3yNqQ3K1M=
Endpoint = 93.184.216.34:51913
PersistentKeepalive = 25
```

_firewall configuration @ 10.172.16.2:_

```
-P INPUT DROP
-P FORWARD DROP
-P OUTPUT DROP
-A INPUT -m state --state ESTABLISHED -j ACCEPT
-A INPUT -i wg0 -j ACCEPT
-A INPUT -s 93.184.216.34 -i enp3s0 -p udp -m udp --sport 51913 -j ACCEPT
-A OUTPUT -d 93.184.216.34/32 -p tcp -m tcp --dport 22 -j ACCEPT
-A OUTPUT -o wg0 -j ACCEPT
-A OUTPUT -d 93.184.216.34/32 -o enp3s0 -p udp -m udp --dport 51913 -j ACCEPT
```

## HTTP server on 10.172.16.2 <> 10.172.16.1

With the existing configuration, we can run an HTTP server on 10.172.16.2 and access it from 10.172.16.1.

To keep it simple, we will use the built-in HTTP server from Python:

_10.172.16.2:_
```
$ python -m http.server -b 10.172.16.2 8000
```

As mentioned, we can already access the HTTP server from 10.172.16.1 with the existing configuration:

_10.172.16.1:_
```
$ curl -I 10.172.16.2:8000
HTTP/1.0 200 OK
Server: SimpleHTTP/0.6 Python/3.11.5
Date: Sun, 24 Sep 2023 22:28:14 GMT
Content-type: text/html; charset=utf-8
Content-Length: 187
```

## HTTP server on 10.172.16.2 <> 93.184.216.34

Our goal is to access the HTTP server using the public IP address of the VPN server.
With the following command, we try exactly this. We try to access the HTTP server running on the same host but inside the VPN over the public internet:

_10.172.16.2:_
```
$ curl -v 93.184.216.34:8000
* processing: 93.184.216.34:8000
*   Trying 93.184.216.34:8000...
```

However, currently, we get no response since the firewall of 93.184.216.34 drops the incoming packets.

We can also see this in the captured network traffic since there are no further packets:

_10.172.16.1:_
```
$ tcpdump -i any -n '(host 54.147.66.132 or net 10.172.16.0/24) and not port 22'
tcpdump: data link type LINUX_SLL2
tcpdump: verbose output suppressed, use -v[v]... for full protocol decode
listening on any, link-type LINUX_SLL2 (Linux cooked v2), snapshot length 262144 bytes
22:52:25.007994 eth0  In  IP 54.147.66.132.47746 > 93.184.216.34.8000: Flags [S], seq 224277161, win 64240, options [mss 1452,sackOK,TS val 3994473945 ecr 0,nop,wscale 7], length 0
```

_(54.147.66.132 is the public IP address used by 10.172.16.2)_

To open port 8000 and use DNAT, we update the firewall configuration like this:

```diff
  -P INPUT DROP
  -P FORWARD DROP
  -P OUTPUT DROP
  -A INPUT -p tcp -m tcp --dport 22 -j ACCEPT
  -A INPUT -i wg0 -j ACCEPT
  -A INPUT -i eth0 -p udp -m udp --dport 51913 -j ACCEPT
  -A FORWARD -i wg0 -o wg0 -j ACCEPT
+ -A FORWARD -i eth0 -o wg0 -j ACCEPT
  -A OUTPUT -m state --state ESTABLISHED -j ACCEPT
  -A OUTPUT -o wg0 -j ACCEPT
+ -t nat -A PREROUTING -p tcp -m tcp --dport 8000 -j DNAT --to-destination 10.172.16.2:8000
```

We will now see that the destination IP address of the packet is replaced <span id="ft-1b">[[1]](#ft-1)</span>:

```diff
  23:00:35.336176 eth0  In  IP 54.147.66.132.38308 > 93.184.216.34.8000: Flags [S], seq 511352029, win 64240, options [mss 1452,sackOK,TS val 3994964273 ecr 0,nop,wscale 7], length 0
+ 23:00:35.336226 wg0   Out IP 54.147.66.132.38308 > 10.172.16.2.8000: Flags [S], seq 511352029, win 64240, options [mss 1452,sackOK,TS val 3994964273 ecr 0,nop,wscale 7], length 0
```

As mentioned, we can't route public IP addresses inside our VPN. Therefore, we also need to use SNAT to translate the source IP address to the internal IP address:

```diff
  -t nat -A PREROUTING -p tcp -m tcp --dport 8000 -j DNAT --to-destination 10.172.16.2:8000
+ -t nat -A POSTROUTING -d 10.172.16.2/32 -o wg0 -p tcp -m tcp --dport 8000 -j SNAT --to-source 10.172.16.1
```

We now see that the source IP address is also translated and further packets are captured:

```diff
  23:14:28.136307 eth0  In  IP 54.147.66.132.54866 > 93.184.216.34.8000: Flags [S], seq 1430183439, win 64240, options [mss 1452,sackOK,TS val 3995797073 ecr 0,nop,wscale 7], length 0
- 23:14:28.136339 wg0   Out IP 54.147.66.132.54866 > 10.172.16.2.8000: Flags [S], seq 1430183439, win 64240, options [mss 1452,sackOK,TS val 3995797073 ecr 0,nop,wscale 7], length 0
+ 23:14:28.136339 wg0   Out IP 10.172.16.1.54866 > 10.172.16.2.8000: Flags [S], seq 1430183439, win 64240, options [mss 1452,sackOK,TS val 3995797073 ecr 0,nop,wscale 7], length 0
+ 23:14:28.136375 eth0  Out IP 93.184.216.34.51913 > 54.147.66.132.38785: UDP, length 96
+ 23:14:28.146021 eth0  In  IP 54.147.66.132.38785 > 93.184.216.34.51913: UDP, length 96
+ 23:14:28.146053 wg0   In  IP 10.172.16.2.8000 > 10.172.16.1.54866: Flags [S.], seq 3292104528, ack 1430183440, win 64296, options [mss 1380,sackOK,TS val 445665275 ecr 3995797073,nop,wscale 7], length 0
```

The last two packets are the response from the HTTP server. However, the response does not get forwarded to the <code class="bg-transparent">eth0</code> interface. That's because we only added a rule to allow forwarding from <code class="bg-transparent">eth0</code> to <code class="bg-transparent">wg0</code> but not vice versa. After changing this:

```diff
  -A FORWARD -i eth0 -o wg0 -j ACCEPT
+ -A FORWARD -i wg0 -o eth0 -j ACCEPT
```

we receive the response now:

```diff
  23:31:36.029921 eth0  In  IP 54.147.66.132.42710 > 93.184.216.34.8000: Flags [S], seq 2290103144, win 64240, options [mss 1452,sackOK,TS val 3996824967 ecr 0,nop,wscale 7], length 0
  23:31:36.030009 wg0   Out IP 10.172.16.1.54866 > 10.172.16.2.8000: Flags [S], seq 2290103144, win 64240, options [mss 1452,sackOK,TS val 3996824967 ecr 0,nop,wscale 7], length 0
  23:31:36.030073 eth0  Out IP 93.184.216.34.59194 > 54.147.66.132.38785: UDP, length 96
  23:31:36.039876 eth0  In  IP 54.147.66.132.38785 > 93.184.216.34.59194: UDP, length 96
  23:31:36.039924 wg0   In  IP 10.172.16.2.8000 > 10.172.16.1.42710: Flags [S.], seq 4227529946, ack 2290103145, win 64296, options [mss 1380,sackOK,TS val 446693169 ecr 3996824967,nop,wscale 7], length 0
+ 23:31:36.039947 eth0  Out IP 93.184.216.34.8000 > 54.147.66.132.42710: Flags [S.], seq 4227529946, ack 2290103145, win 64296, options [mss 1380,sackOK,TS val 446693169 ecr 3996824967,nop,wscale 7], length 0
```

The output of <code class="bg-transparent">curl</code> further confirms this:

_10.172.16.2_:
```
$ curl -I 93.184.216.34:8000
HTTP/1.0 200 OK
Server: SimpleHTTP/0.6 Python/3.11.5
Date: Sun, 24 Sep 2023 23:36:26 GMT
Content-type: text/html; charset=utf-8
Content-Length: 187
```

## final configuration

The final firewall configuration of 10.172.16.1 is this:

```diff
  -P INPUT DROP
  -P FORWARD DROP
  -P OUTPUT DROP
  -A INPUT -p tcp -m tcp --dport 22 -j ACCEPT
  -A INPUT -i wg0 -j ACCEPT
  -A INPUT -i eth0 -p udp -m udp --dport 51913 -j ACCEPT
+ -A FORWARD -i eth0 -o wg0 -j ACCEPT
+ -A FORWARD -i wg0 -o eth0 -j ACCEPT
  -A OUTPUT -m state --state ESTABLISHED -j ACCEPT
  -A OUTPUT -o wg0 -j ACCEPT
+ -t nat -A PREROUTING -p tcp -m tcp --dport 8000 -j DNAT --to-destination 10.172.16.2:8000
+ -t nat -A POSTROUTING -d 10.172.16.2/32 -o wg0 -p tcp -m tcp --dport 8000 -j SNAT --to-source 10.172.16.1
```

---

Congratulations! You now know how network address translation works and how port forwarding is implemented using DNAT and SNAT.

---

<small>
  <span id="ft-0">[[0]](#ft-0b) A router in a private network has two IP addresses: a private address which is commonly the first IP address in the used ranged (so 192.168.0.1 if the range 192.168.0.0/16 is used) and a public one assigned by an internet service provider.
  </span><br />
  <span id="ft-1">[[1]](#ft-1b) We can see that it's the same packet since everything is the same (especially the sequence number and TS val) except the destination IP address.
</small>