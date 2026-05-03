---
title: OpenSSL Cheatsheet
date: 2025-04-07
frontpage: true
---

## Encryption & Decryption

<sub>[`openssl-enc`](https://man.archlinux.org/man/openssl-enc.1ssl.en)</sub>

**Encrypt file with AES-256-CTR, base64 encoding, PBKDF2:**

```
$ openssl enc -aes-256-ctr -pbkdf2 -a -in <file>
```

**Decrypt file encrypted with options above:**

```
$ openssl enc -d -aes-256-ctr -pbkdf2 -a -in <file>
```

## SSL

<sub>[`openssl-x509`](https://man.archlinux.org/man/openssl-x509.1ssl.en)</sub>

**Print certificate of website in text form:**

```
$ echo "QUIT" | openssl s_client -connect www.ekzy.is:443 2>/dev/null | openssl x509 -text -noout
```

**Common certificate printing options**

Use these instead of `-text`:

* `-enddate`: certificate dates
* `-ext=subjectAltName`: print Subject Alternative Name

## CSPRNG

<sub>[`openssl-rand`](https://man.archlinux.org/man/openssl-rand.1ssl.en)</sub>

**32 random bytes in hex or base64:**

```
$ openssl rand -hex 32
$ openssl rand -base64 32
```
