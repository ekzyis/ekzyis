---
title: Post-Quantum TLS vs. MTU?
date: 2025-07-13
sn_id: 1037999
frontpage: true
---

**Update Aug 07, 2025:**

_I figured out that my MTU was too high. The Arch Wiki mentions the problem [here](https://wiki.archlinux.org/title/WireGuard#Adjusting_the_MTU_value):_

> In certain cases larger MTU values can lead to unstable or intermittent connection because of unreliable Path MTU discovery (PMTU) along the route.

_After I set the MTU of the virtual network interface to 1380, the problems went away._

---

For a long time, I couldn't figure out why I couldn't access my self-hosted [vaultwarden](https://github.com/dani-garcia/vaultwarden) instance in a VPN via the browser (Brave) sometimes even though I could always access it fine via cURL or the Bitwarden CLI. When it didn't work, the site would just load until the connection times out.

**But today, I figured it out a little bit more: Post-Quantum Cryptography in the TLS v1.3 handshake made the packets so big that my network interface must have choked—or something like that.**

I arrived at this intermediate conclusion by comparing the browser's TLS v1.3 handshake with the one from cURL. I noticed the browser's Client Hello is a lot bigger (1866 vs 517) and has a lot of TCP retransmissions:

![](./tcp_retransmission.webp)
Why is the Client Hello so big? Is it related to the TCP transmissions?

I also noticed that if I forced Firefox—I couldn't figure this out with Brave—to use TLS v1.2 by setting `security.tls.version.max` to 3 in the advanced config (that you can visit if you type about:config into the address bar), the site loaded immediately. So it was definitely related to TLS v1.3, but specifically the implementation in Brave and Firefox, since I could use TLS v1.3 fine with cURL.

[^1]: I couldn't figure out how to do this in Brave via brave://flags/ so I tested this with Firefox.

I then looked further into why it was so big and noticed the unknown key share 4588. Thanks to this [blog post](https://www.netmeister.org/blog/tls-hybrid-kex.html), I learned that this is a post-quantum cryptography thing.

![](./tcp_wow.webp)
Unknown key share 4588 takes up most of the Client Hello TCP packet

Fortunately, Firefox also had a setting to disable this via `security.tls.enable_kyber`.

When I did this, boom, multiple hours of debugging came to a conclusion that I was happy with, at least for now. When I searched for "PQC vs MTU", I found [this blog post](https://dadrian.io/blog/posts/pqc-signatures-2024/). Apparently, PQC has the issue that it's a lot bigger over the wire to what we're used to:

> In more concrete terms, for the server-sent messages, [Cloudflare found](https://blog.cloudflare.com/pq-2024/) that every 1K of additional data added to the server response caused median HTTPS handshake latency increase by around 1.5%. For the ClientHello, Chrome saw a 4% increase in TLS handshake latency when they deployed ML-KEM, which takes up approximate 1K of additional space in the ClientHello. This pushed the size of the ClientHello greater than the standard maximum transmission unit (MTU) of packets on the Internet, ~1400 bytes, causing the ClientHello to be fragmented over two underlying transport layer (TCP or UDP) packets4.

In some way, debugging continues though: I don't understand why the Client Hello wasn't simply fragmented. Afaict, it wasn't which would explain the TCP retransmissions? The MTU of my physical network interface is set to 1500 and my virtual network interface is set to 1380.

However, after I restarted my virtual network interface, it works now in Brave and in Firefox with Kyber enabled and I think it's still not fragmented—or maybe Wireshark doesn't show me that or I don't know what to look for ¯\\\_(ツ)_/¯

Anyway, at least I now know how to reliably access my password manager via the browser, lol.