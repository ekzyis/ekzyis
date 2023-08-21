Title:        WireGuard Packet Forwarding
Date:         2023-08-21
ReadingTime:  10 minutes
Sats:         0

---

# introduction

In my [previous blog post](/blog/20230809-demystifying-wireguard-and-iptables),
I have shown you how to setup your own VPN with [WireGuard](https://wireguard.com/) and [`iptables`](https://wiki.archlinux.org/title/iptables).
We have established a point-to-point connection between two peers where one peer (10.172.16.1) was reachable from the internet:

![](/blog/img/point-to-point.png)

Today, I will explain how more peers can be added to our VPN.
One peer ("the router") will be configured to forward packets between all other peers ("the end devices").
Therefore, our VPN will become a star network:

![](/blog/img/star_network.png)

You could then [install](https://www.wireguard.com/install/#android-play-store-direct-apk-file) WireGuard on your mobile device
and reach all other machines in your VPN from anywhere with internet connection.

To get a deeper understanding how this forwarding works, we will take a brief look at the network traffic with `tcpdump` after we configured our network.

---

# initial configuration

The end devices will start with a point-to-point connection to the router with 10.172.16.1 as its internal IP address.

The WireGuard and firewall configuration for these could therefore look like this:

**_WireGuard configuration for end device:_**

```
[Interface]
Address = 10.172.16.x/32
PrivateKey = <PRIVATE KEY OF PEER>

[Peer]
AllowedIPs = 10.172.16.1/32
PublicKey = <PUBLIC KEY OF ROUTER>
Endpoint = 139.162.153.133:51913
PersistentKeepalive = 25
```

**_firewall configuration for end device:_**

```
-P INPUT DROP
-P FORWARD DROP
-P OUTPUT DROP
-A INPUT -m state --state ESTABLISHED -j ACCEPT
-A INPUT -i wg0 -j ACCEPT
-A INPUT -s 139.162.153.133 -i enp3s0 -p udp -m udp --sport 51913 -j ACCEPT
-A OUTPUT -d 139.162.153.133/32 -p tcp -m tcp --dport 22 -j ACCEPT
-A OUTPUT -o wg0 -j ACCEPT
-A OUTPUT -d 139.162.153.133/32 -o enp3s0 -p udp -m udp --dport 51913 -j ACCEPT
```

Since the router is connected to multiple peers, it will have multiple <code class="bg-transparent">[Peer]</code> sections in its configuration:

**_WireGuard configuration for router:_**

```
[Interface]
Address = 10.172.16.1/32
PrivateKey = <PRIVATE KEY OF ROUTER>
ListenPort = 51913

[Peer]
AllowedIPs = 10.172.16.2/32
PublicKey = <PUBLIC KEY OF PEER>

[Peer]
AllowedIPs = 10.172.16.4/32
PublicKey = <PUBLIC KEY OF PEER>

...

[Peer]
AllowedIPs = 10.172.16.25/32
PublicKey = <PUBLIC KEY OF PEER>
```

**_firewall configuration for router:_**


```
-P INPUT DROP
-P FORWARD DROP
-P OUTPUT DROP
-A INPUT -p tcp -m tcp --dport 22 -j ACCEPT
-A INPUT -i wg0 -j ACCEPT
-A INPUT -i eth0 -p udp -m udp --dport 51913 -j ACCEPT
-A OUTPUT -m state --state ESTABLISHED -j ACCEPT
-A OUTPUT -o wg0 -j ACCEPT
```

<small>_If these configurations are confusing to you, read my [previous blog post](/blog/20230809-demystifying-wireguard-and-iptables)._</small>

---

# end device configuration

The only change we have to do on the end devices is to route all IP addresses within the VPN to the router peer.
We configure this in the WireGuard configuration file at _/etc/wireguard/wg0.conf_:

```diff
  [Peer]
- AllowedIPs = 10.172.16.1/32
+ AllowedIPs = 10.172.16.0/24
  PublicKey = &ltPUBLIC KEY OF ROUTER&gt
  Endpoint = 139.162.153.133:51913
  PersistentKeepalive = 25
```

To apply these changes, we run this command <span id="ft-0b">[[0]](#ft-0)</span>:

```
$ wg syncconf wg0 <(wg-quick strip wg0)
```

---

# router configuration

## firewall

The router requires no additional WireGuard configuration.

However, the router firewall needs to allow forwarding packets inside the VPN:

```diff
  -P INPUT DROP
  -P FORWARD DROP
  -P OUTPUT DROP
  -A INPUT -p tcp -m tcp --dport 22 -j ACCEPT
  -A INPUT -i wg0 -j ACCEPT
  -A INPUT -i eth0 -p udp -m udp --dport 51913 -j ACCEPT
  -A OUTPUT -m state --state ESTABLISHED -j ACCEPT
  -A OUTPUT -o wg0 -j ACCEPT
+ -A FORWARD -i wg0 -o wg0 -j ACCEPT
```

## kernel parameters

We also need to allow IP forwarding in the kernel. We can configure kernel parameters with `sysctl`.

To see the current value for IP forwarding:

```
$ sysctl net.ipv4.ip_forward
net.ipv4.ip_forward = 0
```

This basically does the same as checking the content of the file _/proc/sys/net/ipv4/ip\_forward_:

```
$ cat /proc/sys/net/ipv4/ip_forward
0
```

To change the setting, we can use `sysctl -w`:

```
$ sysctl -w net.ipv4.ip_forward=0
```

However, this change is not persistent. To make sure the new setting is kept after reboot, we need to modify _/etc/sysctl.conf_.
We can do this by appending <code class="bg-transparent">net.ipv4.ip_forward = 1</code> to the file:

```
$ echo 'net.ipv4.ip_forward = 1' >> /etc/sysctl.conf
```

To reload settings:

```
$ sysctl -p
```

---

# inspecting the network traffic with tcpdump

We now should have a connection between every peer via the router.

To confirm this, we ping one machine (10.172.16.25) from another machine (10.172.16.6):

```
$ ping 10.172.16.25
PING 10.172.16.25 (10.172.16.25) 56(84) bytes of data.
64 bytes from 10.172.16.25: icmp_seq=1 ttl=63 time=138 ms
64 bytes from 10.172.16.25: icmp_seq=2 ttl=63 time=165 ms
64 bytes from 10.172.16.25: icmp_seq=3 ttl=63 time=182 ms
64 bytes from 10.172.16.25: icmp_seq=4 ttl=63 time=206 ms
64 bytes from 10.172.16.25: icmp_seq=5 ttl=63 time=229 ms
```

We will inspect this traffic on 10.172.16.1 (the router) and 10.172.16.6 (the pinging machine) with `tcpdump` now <span id="ft-1b">[[1]](#ft-1)</span>.

On the pinging machine, we will get the following output for the virtual network interface <span id="ft-2b">[[2]](#ft-2)</span>:

```tcpdump
$ tcpdump -tni wg0
tcpdump: verbose output suppressed, use -v[v]... for full protocol decode
listening on wg0, link-type RAW (Raw IP), snapshot length 262144 bytes
IP 10.172.16.6 > 10.172.16.25: ICMP echo request, id 19, seq 1, length 64
IP 10.172.16.25 > 10.172.16.6: ICMP echo reply, id 19, seq 1, length 64
IP 10.172.16.6 > 10.172.16.25: ICMP echo request, id 19, seq 2, length 64
IP 10.172.16.25 > 10.172.16.6: ICMP echo reply, id 19, seq 2, length 64
IP 10.172.16.6 > 10.172.16.25: ICMP echo request, id 19, seq 3, length 64
IP 10.172.16.25 > 10.172.16.6: ICMP echo reply, id 19, seq 3, length 64
IP 10.172.16.6 > 10.172.16.25: ICMP echo request, id 19, seq 4, length 64
IP 10.172.16.25 > 10.172.16.6: ICMP echo reply, id 19, seq 4, length 64
IP 10.172.16.6 > 10.172.16.25: ICMP echo request, id 19, seq 5, length 64
IP 10.172.16.25 > 10.172.16.6: ICMP echo reply, id 19, seq 5, length 64
```

and this for the physical network interface (filtered by UDP packets from/to port 51913):

```tcpdump
$ tcpdump -tni enp3s0 'udp and port 51913'
tcpdump: verbose output suppressed, use -v[v]... for full protocol decode
listening on enp3s0, link-type EN10MB (Ethernet), snapshot length 262144 bytes
IP 192.168.178.146.51941 > 139.144.78.247.51913: UDP, length 128
IP 139.144.78.247.51913 > 192.168.178.146.51941: UDP, length 128
IP 192.168.178.146.51941 > 139.144.78.247.51913: UDP, length 128
IP 139.144.78.247.51913 > 192.168.178.146.51941: UDP, length 128
IP 192.168.178.146.51941 > 139.144.78.247.51913: UDP, length 128
IP 139.144.78.247.51913 > 192.168.178.146.51941: UDP, length 128
IP 192.168.178.146.51941 > 139.144.78.247.51913: UDP, length 128
IP 139.144.78.247.51913 > 192.168.178.146.51941: UDP, length 128
IP 192.168.178.146.51941 > 139.144.78.247.51913: UDP, length 128
IP 139.144.78.247.51913 > 192.168.178.146.51941: UDP, length 128
```

On the router, we get this output for the virtual network interface:

```tcpdump
$ tcpdump -tni wg0
tcpdump: verbose output suppressed, use -v[v]... for full protocol decode
listening on wg0, link-type RAW (Raw IP), snapshot length 262144 bytes
IP 10.172.16.6 > 10.172.16.25: ICMP echo request, id 21, seq 1, length 64
IP 10.172.16.6 > 10.172.16.25: ICMP echo request, id 21, seq 1, length 64
IP 10.172.16.25 > 10.172.16.6: ICMP echo reply, id 21, seq 1, length 64
IP 10.172.16.25 > 10.172.16.6: ICMP echo reply, id 21, seq 1, length 64
IP 10.172.16.6 > 10.172.16.25: ICMP echo request, id 21, seq 2, length 64
IP 10.172.16.6 > 10.172.16.25: ICMP echo request, id 21, seq 2, length 64
IP 10.172.16.25 > 10.172.16.6: ICMP echo reply, id 21, seq 2, length 64
IP 10.172.16.25 > 10.172.16.6: ICMP echo reply, id 21, seq 2, length 64
IP 10.172.16.6 > 10.172.16.25: ICMP echo request, id 21, seq 3, length 64
IP 10.172.16.6 > 10.172.16.25: ICMP echo request, id 21, seq 3, length 64
IP 10.172.16.25 > 10.172.16.6: ICMP echo reply, id 21, seq 3, length 64
IP 10.172.16.25 > 10.172.16.6: ICMP echo reply, id 21, seq 3, length 64
IP 10.172.16.6 > 10.172.16.25: ICMP echo request, id 21, seq 4, length 64
IP 10.172.16.6 > 10.172.16.25: ICMP echo request, id 21, seq 4, length 64
IP 10.172.16.25 > 10.172.16.6: ICMP echo reply, id 21, seq 4, length 64
IP 10.172.16.25 > 10.172.16.6: ICMP echo reply, id 21, seq 4, length 64
IP 10.172.16.6 > 10.172.16.25: ICMP echo request, id 21, seq 5, length 64
IP 10.172.16.6 > 10.172.16.25: ICMP echo request, id 21, seq 5, length 64
IP 10.172.16.25 > 10.172.16.6: ICMP echo reply, id 21, seq 5, length 64
IP 10.172.16.25 > 10.172.16.6: ICMP echo reply, id 21, seq 5, length 64
```

and this for the physical network interface _(public IP addresses of end devices redacted)_:

```tcpdump
$ tcpdump -tni eth0 'udp and port 51913'
tcpdump: verbose output suppressed, use -v[v]... for full protocol decode
listening on eth0, link-type EN10MB (Ethernet), snapshot length 262144 bytes
IP 87.161.X.X.51941 > 139.144.78.247.51913: UDP, length 128
IP 139.144.78.247.51913 > 109.43.X.X.11007: UDP, length 128
IP 109.43.X.X.11007 > 139.144.78.247.51913: UDP, length 128
IP 139.144.78.247.51913 > 87.161.X.X.51941: UDP, length 128
IP 87.161.X.X.51941 > 139.144.78.247.51913: UDP, length 128
IP 139.144.78.247.51913 > 109.43.X.X.11007: UDP, length 128
IP 109.43.X.X.11007 > 139.144.78.247.51913: UDP, length 128
IP 139.144.78.247.51913 > 87.161.X.X.51941: UDP, length 128
IP 87.161.X.X.51941 > 139.144.78.247.51913: UDP, length 128
IP 139.144.78.247.51913 > 109.43.X.X.11007: UDP, length 128
IP 109.43.X.X.11007 > 139.144.78.247.51913: UDP, length 128
IP 139.144.78.247.51913 > 87.161.X.X.51941: UDP, length 128
IP 87.161.X.X.51941 > 139.144.78.247.51913: UDP, length 128
IP 139.144.78.247.51913 > 109.43.X.X.11007: UDP, length 128
IP 109.43.X.X.11007 > 139.144.78.247.51913: UDP, length 128
IP 139.144.78.247.51913 > 87.161.X.X.51941: UDP, length 128
IP 87.161.X.X.51941 > 139.144.78.247.51913: UDP, length 128
IP 139.144.78.247.51913 > 109.43.X.X.11007: UDP, length 128
IP 109.43.X.X.11007 > 139.144.78.247.51913: UDP, length 128
IP 139.144.78.247.51913 > 87.161.X.X.51941: UDP, length 128
```

As you can see, for every ICMP echo request, we get an ICMP echo reply from 10.172.16.25 in the `tcpdump` output of 10.172.16.6.
We also see that we receive a response for every UDP packet sent to port 51913 of 139.144.78.247.

However, in the `tcpdump` output of 10.172.16.1 for <code class="bg-transparent">wg0</code>,
we see every packet twice. Why is that?

The reason for this is that `tcpdump` captures every incoming and outgoing packet <span id="ft-3b">[[3]](#ft-3)</span>.

Since we are forwarding packets on the virtual network interface, the packet for 10.172.16.25 from 10.172.16.6 will be captured on 10.172.16.1
when it enters the interface from 10.172.16.6 and when it leaves to 10.172.16.25.

The same happens when the response from 10.172.16.25 arrives at 10.172.16.1 and is forwarded to 10.172.16.6.

We can see this forwarding happening better in the `tcpdump` output for the physical network interface. There, we see that for every ICMP echo request, we see one UDP
packet from 87.161.X.X to 139.144.78.247 and then another from 139.144.78.247 to 109.43.X.X.

As you also can see, the forwarding does not change internal IP addresses. The packets still contain the same source and destination IP addresses.
This works because we have configured the interfaces of all peers to route all IP addresses towards the router peer.
So when 10.172.16.25 sees the ICMP echo request from 10.172.16.6 via 10.172.16.1, it will just respond to 10.172.16.6 again via 10.172.16.1.

The UDP packet IP addresses do seem changed, though. However, this is not network address translation (NAT) but just packet encapsulation.
As mentioned in my [previous blog post](/blog/20230809-demystifying-wireguard-and-iptables), the physical network interfaces have no notion of internal IP addresses
thus every UDP packet is sent using publicly routable IP addresses over the wire.

This means we are not doing any network address translation here hence the name forwarding.

Another important detail is that we do not need the following FORWARD rules in the router:

```
-A FORWARD -s 10.172.16.0/24 -i wg0 -o eth0 -j ACCEPT
-A FORWARD -d 10.172.16.0/24 -i eth0 -o wg0 -j ACCEPT
```

These rules are only required if we want to forward VPN traffic to <code class="bg-transparent">eth0</code> which faces the public internet.
Therefore, these rules would be part of a configuration to give internet access via 10.172.16.1 to peers inside the VPN.

---

Thanks for reading my second blog post! If you want to read more content like this, please consider subscribing via [RSS](https://dev.ekzyis.com/blog/rss.xml).

In the next blog post, we will use network address translation with the SNAT and DNAT targets in `iptables` for port forwarding.
This will make it possible to expose an internal service (like an HTTP server for example) to the public internet.

---

<small>
  <span id="ft-0">[[0]](#ft-0b) `wg-quick strip wg0` returns the config in the format that `wg` can parse.
  This is necessary because `wg-quick` <q>adds a few extra configuration values to the format understood by `wg` in order to
  configure additional attributes of an interface</q>. ([source](https://man.archlinux.org/man/wg-quick.8.en#CONFIGURATION))
  </span><br />
  <span id="ft-1">[[1]](#ft-1b) Since 10.172.16.25 is a mobile device, I didn't include `tcpdump` output from that device. It also wouldn't include anything interesting.
  </span><br />
  <span id="ft-2">[[2]](#ft-2b) I used `-t` to not include timestamps and `-n` to not resolve IP addresses to hostnames. `-i <interface>` selects the interface we want to tap.
  </span><br />
  <span id="ft-3">[[3]](#ft-3b) To be precise, `tcpdump` captures incoming packets _before_ firewall processing
  while outgoing packets will be captured _after_ firewall processing. This means `tcpdump` will capture incoming packets that will be dropped by the firewall
  whereas outgoing packets that were dropped will not show up.
  See [here](https://wiki.archlinux.org/title/Network_Debugging#Tcpdump) and [here](https://superuser.com/a/925332) for more information.
  </span>
</small>
