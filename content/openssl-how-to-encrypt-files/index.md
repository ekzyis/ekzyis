---
title: How to encrypt and decrypt files using openssl and a password
date: 2025-04-07
sn_id: 937345
tags: crypto openssl
publish: false
---

Every time I feel like I should save something encrypted—especially because my [NixOS installation is not using FDE](https://nixos.wiki/wiki/Full_Disk_Encryption), but that's a story for another day—, I realize I have no idea how to do this except that it should be possible with `openssl` which is quite embarrassing as the founder of ~crypto and ~security.

So I am writing this post to never need to look at the [massive manual for `openssl`](https://man.archlinux.org/man/openssl.1ssl.en) just to encrypt a file again.

**Encryption**

Encrypt file with AES-256-CTR, use [PBKDF2](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#pbkdf2) as the password hashing algorithm and use base64 encoding for the output:

```
$ openssl enc -aes-256-ctr -pbkdf2 -a -in <file>
```

_will prompt for password_

**Decryption**

Same as above with but `-d`:

```
$ openssl enc -d -aes-256-ctr -pbkdf2 -a -in <file>
```

---

Btw, I used AES-256-CTR here because I know it _should_ be safe and I don't want to research all the other ciphers I could use. Maybe ChaCha20 would be safer??

Which cipher would you use to store, say, an API key? Would you even use `openssl` for that since it's quite low-level and you can probably easily shoot yourself in the foot with it?
