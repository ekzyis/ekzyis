---
title: The Curious Case of Digital Signatures
date: 2024-02-24T23:53:31.373Z
sn_id: 436752
tags: crypto
---

# Good questions demand good answers

[@Natalia](https://stacker.news/Natalia) asked [some good questions](https://stacker.news/items/436029) in the [@saloon](https://stacker.news/saloon) today:

> Playing around PGP last night ( yes, late into the game - it takes time to sort things out one by one. ) so here is my understanding of the logic of verifying software:
>
> 1. you import the specific APP's public key from the terminal
> 2. download both the software and the signature file (asc)
> 3. copy and paste a line of code into the terminal to verify
>
> **But how does it work in the background - is it verifying the signature file matches with the public key? and how things be if someone changed the app, the signature then unmatched with the public key?** 🤔

I like such questions since they motivate me to dig deeper. I actually didn't really know anymore how it works in detail so I wasn't able to answer these questions with confidence but I would have loved to. I mean, I am the founder of [~crypto](https://stacker.news/~crypto)! What kind of LARPer would I be if I can't even explain how `gpg --verify` works?

So I decided to be less of a LARPer and do some research. In this post, I present you the result of my research and my attempt to explain this in an easy to understand way.[^1]

I'll use [Sparrow Wallet](https://sparrowwallet.com/) and a person that is learning from first principles how signatures are verified with and by `gpg --verify` (basically me after getting [nerd sniped](https://xkcd.com/356/) by [@Natalia](https://stacker.news/Natalia)'s questions) as an example.

# Purpose of Digital Signatures

![913.gif](https://m.stacker.news/17250)

Digital signatures are commonly used to ensure the integrity and authenticity of software. When you verify a digital signature, you make sure that the software was created by the person you trust and think it was created by (authenticity) and that it was not modified (integrity).

You usually download the signature from the same location as the software. For example, this location is https://sparrowwallet.com/download for Sparrow Wallet:

![2024-02-23-195709_1920x1080_scrot.png](https://m.stacker.news/17249)

The signature is the file `sparrow-1.8.2-manifest.txt.asc`. It contains the following content:

```
-----BEGIN PGP SIGNATURE-----
Comment: GPGTools - http://gpgtools.org

iQIzBAABCgAdFiEE1NDTIC/AaEmiV7ON6UYYM0xnS0AFAmWo/vYACgkQ6UYYM0xn
S0B40RAAteFi0MBNKY80txOOFmcChAAgcLSnaH3Hd8AaKjijFdvYGvUAebP1EMBj
58K2p14aHq6rdEn2it8aQV4pO9sk9l9q6Fp0YObkJoWkcLL4r1fy4RVxgYFTo7D4
UkpigS4ulxrot5ChSzVxD5HTOsCFk1xrrWr77vWE2PiZ4w+L8d1+6rRdzPp6Yftl
7rF3Hb5x4aiZVKBqoSYbmBiPybLnCdXP0gxSYd4JZeyQ44XSOirNEO07w+tdwYPP
1fR09nBdcnF/BD1SPjCKL3P0YqWE04vJ+7XesdesMReAcxAntjlfRmGAlbbX0v60
MzNoM6vN/H9rfY6WXCJ5JJpdufulm11GbYAbC8rbSBNTXwxVY6czhMmG2cuQ7Euy
AJ+1sPzh5U+hWAELOxI80atkQ8KSmw0kBnEP2H4bNMB4C48gU3tX2rEmvIZTHTN7
QLUr08y5FscWftlwJ1uqkGVr32JKZNnoLgL0VVlwlPquHNeX5OjNn6nWc3w2AQ7q
j8tDZlzLnvFhvvRL4l1EO65PEg0zoL0LJVRSuoEs85o6ZKSJ0rUH9+jO1+4kUnzx
ypR4KM40swCrP4DpuNLuc4p4Os0XFDdRKxcx2q3PH06+N1soAh5NJiXZRsOlEl9q
NitWcDrzqKdZfYTevuTzbDMk/EikdRzfIe4KgcsutEiaKk7opQc=
=m0CZ
-----END PGP SIGNATURE-----
```

Let's assume you don't trust this location for some reason. You might be a new user and thus don't know if this is the official, legit website to download Sparrow. Or you simply don't want to rely solely on [SSL/TLS certificates](https://www.cloudflare.com/learning/ssl/what-is-an-ssl-certificate/).

However, you want to learn how to verify digital signatures as mentioned. You want to be a smart cookie!

# The Insanity of [Web of Trust](https://en.wikipedia.org/wiki/Web_of_trust)

![download (2).jpeg](https://m.stacker.news/17274)

You start with importing the public key as explained in the download page:

```
$ curl https://keybase.io/craigraw/pgp_keys.asc | gpg --import
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  5588  100  5588    0     0  14944      0 --:--:-- --:--:-- --:--:-- 15021
gpg: key E94618334C674B40: public key "Craig Raw <craig@sparrowwallet.com>" imported
gpg: Total number processed: 1
gpg:               imported: 1
```

**`One important detail is here that the public key is not provided via the same location. If the location is compromised, it could simply replace the public key with a public key controlled by an adversary which in turn compromises the whole scheme.`**

However, the instructions to download the public key are provided at the same location. So couldn't they also have been replaced?

Exactly right! However, we can go to https://keybase.io/craigraw and check if there is some useful information that makes the public key more trustworthy[^2]:

![2024-02-23-233532_1920x1080_scrot.png](https://m.stacker.news/17267)

We find some linked profiles and followers as sources of trust. We find proofs next to the linked profiles.

[Here](https://gist.github.com/craigraw/20311f22d738fd9f8a1ee0007f633f5a) is the proof that the Github account [craigraw](https://github.com/craigraw) controls this account on Keybase:

![2024-02-23-233921_1920x1080_scrot.png](https://m.stacker.news/17268)

and [here](https://twitter.com/craigraw/status/1179755136387883008) is a tweet proving the same for [@craigraw](https://twitter.com/craigraw):

![2024-02-23-233905_600x318_scrot.png](https://m.stacker.news/17269)

Finally, we check that the key fingerprint matches the fingerprint shown in the proofs:

![2024-02-23-235858_1920x1080_scrot.png](https://m.stacker.news/17272)

Checking the fingerprint of the key we received via `gpg --list-keys --keyid-format long`:

```
$ gpg --list-keys --keyid-format long craig@sparrowwallet.com
pub   rsa4096/E94618334C674B40 2019-10-03 [SC] [expires: 2027-09-18]
      D4D0D3202FC06849A257B38DE94618334C674B40
uid                 [ unknown] Craig Raw <craig@sparrowwallet.com>
sub   rsa4096/3FB2D3544ECAF6DA 2019-10-03 [E] [expires: 2027-09-18]
```

The fingerprint `D4D0D3202FC06849A257B38DE94618334C674B40` matches. Great!

After you did your due diligence like a real believer in Web of Trust, you can use `gpg --edit-key craig@sparrowwallet.com` to specify how much you trust this key now:

```
$ gpg --edit-key craig@sparrowwallet.com
gpg (GnuPG) 2.4.3; Copyright (C) 2023 g10 Code GmbH
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.


pub  rsa4096/E94618334C674B40
     created: 2019-10-03  expires: 2027-09-18  usage: SC
     trust: unknown       validity: unknown
sub  rsa4096/3FB2D3544ECAF6DA
     created: 2019-10-03  expires: 2027-09-18  usage: E
[ unknown] (1). Craig Raw <craig@sparrowwallet.com>
[ revoked] (2)  Craig Raw <craigraw@gmail.com>

gpg>
```

Enter the command `trust` to edit the key trust:

```
gpg> trust
pub  rsa4096/E94618334C674B40
     created: 2019-10-03  expires: 2027-09-18  usage: SC
     trust: unknown       validity: unknown
sub  rsa4096/3FB2D3544ECAF6DA
     created: 2019-10-03  expires: 2027-09-18  usage: E
[ unknown] (1). Craig Raw <craig@sparrowwallet.com>
[ revoked] (2)  Craig Raw <craigraw@gmail.com>

Please decide how far you trust this user to correctly verify other users' keys
(by looking at passports, checking fingerprints from different sources, etc.)

  1 = I don't know or won't say
  2 = I do NOT trust
  3 = I trust marginally
  4 = I trust fully
  5 = I trust ultimately
  m = back to the main menu

Your decision?
```

_The option 5 "I trust ultimately" should only be used with your own keys._

# Sanity checks

![download.jpeg](https://m.stacker.news/17253)

After you have ~~decided that you don't care about Web of Trust~~ verified the received public key, you download the signature now. But you don't download the software just yet. You want to do a [sanity check](https://en.wikipedia.org/wiki/Sanity_check):

> _If there is no software to verify, running `gpg --verify` shouldn't make sense, right? But what would it do?_

So like the curious person that you are and that assumes sanity is not a given but must be earned, you run `gpg --verify <signature>` in a folder which contains only the signature (`sparrow-1.8.2-manifest.txt.asc`) and expect the command to fail in some way. If you are on the right track, the error message should also not surprise you:

```
$ ls
total 4.0K
-rw-r--r-- 1 ekzyis ekzyis 873 Feb 23 15:42 sparrow-1.8.2-manifest.txt.asc
$ gpg --verify sparrow-1.8.2-manifest.txt.asc
gpg: no signed data
gpg: can't hash datafile: No data
```

Looks like `gpg` indeed threw an error! Congratulations, you are not insane (yet) and maybe even on a good path to absolute wisdom (or at least being a smart cookie). Additionally, it seems like it even threw two errors: "no signed data" and "can't hash datafile: No data".

You are not sure what they mean though. So you start reading [the documentation](https://www.gnupg.org/gph/en/manual/r697.html) of `gpg --verify` since [RTFM](https://archive.ph/KbrZ6) is always a good start to understand things better:

> `verify <signature> <document>`
>
> This command verifies a document against a signature to ensure that the document has not been altered since the signature was created. If `<signature>` is omitted, gpg will look in `<document>` for a clearsign signature.

Since we only gave a single argument to `gpg --verify`, `gpg` tried to find a _clearsign signature_ in the provided file. But what is a clearsign signature? You start reading [more documentation](https://www.gnupg.org/gph/en/manual/x135.html), this time for `gpg --sign`:

> **Clearsigned documents**
>
> A common use of digital signatures is to sign usenet postings or email messages. In such situations it is undesirable to compress the document while signing it. The option --clearsign causes the document to be wrapped in an ASCII-armored signature but otherwise does not modify the document.
>
> ```
> alice% gpg --clearsign doc
>
> You need a passphrase to unlock the secret key for
> user: "Alice (Judge) <alice@cyb.org>"
> 1024-bit DSA key, ID BB7576AC, created 1999-06-04
>
> -----BEGIN PGP SIGNED MESSAGE-----
> Hash: SHA1
>
> [...]
> -----BEGIN PGP SIGNATURE-----
> Version: GnuPG v0.9.7 (GNU/Linux)
> Comment: For info see http://www.gnupg.org
>
> iEYEARECAAYFAjdYCQoACgkQJ9S6ULt1dqz6IwCfQ7wP6i/i8HhbcOSKF4ELyQB1
> oCoAoOuqpRqEzr4kOkQqHRLE/b8/Rw2k
> =y6kj
> -----END PGP SIGNATURE-----
> ```

Aha! You look at the signature that you have and see that it only seems to include a signature but not a document or message. There is no `-----BEGIN PGP SIGNED MESSAGE-----`:

```
-----BEGIN PGP SIGNATURE-----
Comment: GPGTools - http://gpgtools.org

iQIzBAABCgAdFiEE1NDTIC/AaEmiV7ON6UYYM0xnS0AFAmWo/vYACgkQ6UYYM0xn
S0B40RAAteFi0MBNKY80txOOFmcChAAgcLSnaH3Hd8AaKjijFdvYGvUAebP1EMBj
58K2p14aHq6rdEn2it8aQV4pO9sk9l9q6Fp0YObkJoWkcLL4r1fy4RVxgYFTo7D4
UkpigS4ulxrot5ChSzVxD5HTOsCFk1xrrWr77vWE2PiZ4w+L8d1+6rRdzPp6Yftl
7rF3Hb5x4aiZVKBqoSYbmBiPybLnCdXP0gxSYd4JZeyQ44XSOirNEO07w+tdwYPP
1fR09nBdcnF/BD1SPjCKL3P0YqWE04vJ+7XesdesMReAcxAntjlfRmGAlbbX0v60
MzNoM6vN/H9rfY6WXCJ5JJpdufulm11GbYAbC8rbSBNTXwxVY6czhMmG2cuQ7Euy
AJ+1sPzh5U+hWAELOxI80atkQ8KSmw0kBnEP2H4bNMB4C48gU3tX2rEmvIZTHTN7
QLUr08y5FscWftlwJ1uqkGVr32JKZNnoLgL0VVlwlPquHNeX5OjNn6nWc3w2AQ7q
j8tDZlzLnvFhvvRL4l1EO65PEg0zoL0LJVRSuoEs85o6ZKSJ0rUH9+jO1+4kUnzx
ypR4KM40swCrP4DpuNLuc4p4Os0XFDdRKxcx2q3PH06+N1soAh5NJiXZRsOlEl9q
NitWcDrzqKdZfYTevuTzbDMk/EikdRzfIe4KgcsutEiaKk7opQc=
=m0CZ
-----END PGP SIGNATURE-----
```

You scroll a bit down in the documentation and see that it also mentions _detached signatures_:

> **Detached signatures**
>
> A signed document has limited usefulness. Other users must recover the original document from the signed version, and even with clearsigned documents, the signed document must be edited to recover the original. Therefore, there is a third method for signing a document that creates a detached signature. A detached signature is created using the --detach-sig option.

Well, that sounds interesting. You are now happy about your reading comprehension skills and are pretty sure that what you have is a detached signature. This explains both errors from before. There was no data included and `gpg` then tried to find the document that was signed in some way but couldn't find it. Poor `gpg`.

The online documentation however does not seem to explain where `gpg` tries to find the document. However, there is also the command `man` which stands for the M in RTFM. So you give `man gpg` a shot and search for `--verify`. You find this note among many other notes:

> **--verify**:
>
>              [...]
>
>              Note:  If the option --batch is not used, gpg may assume that a single argument is a file with a detached signature, and it will try to find a matching data file by stripping certain suffixes.
>              Using this historical feature to verify a detached signature is strongly discouraged; you should always specify the data file explicitly.

You ignore the "strongly discouraged" and want to do another sanity check:

> _If I understand correctly what I just read, `gpg` tried to remove "certain suffixes" to find the data file. I wonder what will happen if I download the software and put it in the same folder ..._

So you do exactly that and run `gpg --verify` again:

```
$ ls
total 106M
-rw-r--r-- 1 ekzyis ekzyis  873 Feb 23 15:42 sparrow-1.8.2-manifest.txt.asc
-rw-r--r-- 1 ekzyis ekzyis 106M Feb 23 20:17 sparrow-1.8.2-x86_64.tar.gz
$ gpg --verify sparrow-1.8.2-manifest.txt.asc
gpg: no signed data
gpg: can't hash datafile: No data
```

Mhh, same error. So something is wrong. You double-check the file names and realize that when you strip `.txt.asc`, you get `sparrow-1.8.2-manifest` and not `sparrow-1.8.2-x86_64.tar.gz`! So `gpg` could still not find what you seem to call "software", the documentation seems to call "document" and the program itself seems to call "datafile" (consistency must be a virtue) simply by stripping suffixes. Everything makes enough sense again—for now.

You check the download page again and see that there is a file named `sparrow-1.8.2-manifest.txt`.

> _That looks very promising. It's the same file name as the signature just without `.asc` at the end. Now it should work._

You download it and look inside:

```
$ cat sparrow-1.8.2-manifest.txt
541b82e83ec928e7b54490211fa5e9d3b315394d895aebeab9a9696ce24c378d *Sparrow-1.8.2-aarch64.dmg
508a5ba212c2393a9d4cd122c2a8aa8b16e593c640bc3e4ac111ee574b2fd82d *Sparrow-1.8.2-x86_64.dmg
1035fc5663e53ce3a7d87a63a1dadc33bb180ceb3154549b251407b09db202b7 *Sparrow-1.8.2.exe
8ff8e70886af91c8758509097a1047a05341b29a2407934c49b6885ff71c60c2 *Sparrow-1.8.2.zip
afee92405d3013d24add5e3099cea86f5201829b8f5f485f26b753af109619b9 *sparrow-1.8.2-1.aarch64.rpm
22b72b62865d8e9d0dcf71dc0907a3945f7a5c859257db4f05ee5b815c4bc56d *sparrow-1.8.2-1.x86_64.rpm
97d43340c6938dbd513c62e3528ae61f57cefe308cd950b0d0cf3f7a3a36993f *sparrow-1.8.2-aarch64.tar.gz
ffb7f86e978ab312dcc169d7c70b33048256dcf3280aecb1cb9392124eac31b9 *sparrow-1.8.2-x86_64.tar.gz
e69873062837f90745b384d53afeb590cb0058061da3148cbad35604b3e4b1bd *sparrow-server-1.8.2-1.aarch64.rpm
3fb4900e078427c230dd1372e1b56abdaa71653c5247d6ca4605e9fdda90bffe *sparrow-server-1.8.2-1.x86_64.rpm
6278fe3b8261576bb9e597ac0e1a5305781917a6716dfc8c5ebb8b244c1a8805 *sparrow-server-1.8.2-aarch64.tar.gz
0f2025a3fee1057fef763ff1ead471b4468766c0db99e07d8cd1de26e359852d *sparrow-server-1.8.2-x86_64.tar.gz
c538a93d22698a46218859fbf936731484683e188da33a2e35dac053388a5e14 *sparrow-server_1.8.2-1_amd64.deb
817469b84741832d4b58b0f6305139cf740ba6b91ca9ff5e65dac75f1fa0221e *sparrow-server_1.8.2-1_arm64.deb
ef461e5b75681e39fde235ee75d0bb23d4cfbeb30c312383f17596115d0a9188 *sparrow_1.8.2-1_amd64.deb
f81b259799d3dbf7ac75555b59ac719256fcfacb716c653ca445f21f327add5c *sparrow_1.8.2-1_arm64.deb
```

You see a lot of stuff. You don't care about it for now since you want to keep your sanity. You simply run `gpg --verify` again and ...

```
$ ls
total 106M
-rw-r--r-- 1 ekzyis ekzyis 1.5K Feb 23 16:12 sparrow-1.8.2-manifest.txt
-rw-r--r-- 1 ekzyis ekzyis  873 Feb 23 15:42 sparrow-1.8.2-manifest.txt.asc
-rw-r--r-- 1 ekzyis ekzyis 106M Feb 23 20:17 sparrow-1.8.2-x86_64.tar.gz
$ gpg --verify sparrow-1.8.2-manifest.txt.asc
gpg: assuming signed data in 'sparrow-1.8.2-manifest.txt'
gpg: Signature made Thu 18 Jan 2024 11:35:34 AM CET
gpg:                using RSA key D4D0D3202FC06849A257B38DE94618334C674B40
gpg: Good signature from "Craig Raw <craig@sparrowwallet.com>" [unknown]
gpg: WARNING: This key is not certified with a trusted signature!
gpg:          There is no indication that the signature belongs to the owner.
Primary key fingerprint: D4D0 D320 2FC0 6849 A257  B38D E946 1833 4C67 4B40
```

... it worked! You get "Good signature". That sounds good. But you wonder about the warning:

```
gpg: WARNING: This key is not certified with a trusted signature!
gpg:          There is no indication that the signature belongs to the owner.
```

Haven't you done your due diligence in the previous chapter? Well, I have and I still get this error for some reason, lol. Anyway ... you see that the Sparrow download page mentions exactly this case:

> Note that you may get a message similar to the following:
>
> ```
> gpg: WARNING: This key is not certified with a trusted signature!
> gpg:          There is no indication that the signature belongs to the owner.
> ```
>
> This simply means that you have not explicitly marked the public key as trusted in your own instance of GPG. In this case it is good practice to check the key against other sources, for example https://keybase.io/craigraw (click on the link next to the key icon to see the full public key). You can read more about validating keys in the [GnuPG Privacy Handbook](https://www.gnupg.org/gph/en/manual/x334.html).

But we want to focus on verification for now. Never lose your sanity from multiple directions at once.

But aren't we already done? We got a good signature so we verified the software right?

You remember that the detached signature didn't look for the downloaded software since it didn't match the file without the `.asc` suffix. So something is off. You run another sanity check by deleting the software (or moving it somewhere else) and run `gpg --verify` for the fourth time:

```
$ ls
total 8.0K
-rw-r--r-- 1 ekzyis ekzyis 1.5K Feb 23 16:12 sparrow-1.8.2-manifest.txt
-rw-r--r-- 1 ekzyis ekzyis  873 Feb 23 15:42 sparrow-1.8.2-manifest.txt.asc
$ gpg --verify sparrow-1.8.2-manifest.txt.asc
gpg: assuming signed data in 'sparrow-1.8.2-manifest.txt'
gpg: Signature made Thu 18 Jan 2024 11:35:34 AM CET
gpg:                using RSA key D4D0D3202FC06849A257B38DE94618334C674B40
gpg: Good signature from "Craig Raw <craig@sparrowwallet.com>" [unknown]
gpg: WARNING: This key is not certified with a trusted signature!
gpg:          There is no indication that the signature belongs to the owner.
Primary key fingerprint: D4D0 D320 2FC0 6849 A257  B38D E946 1833 4C67 4B40
```

It still works. That was expected but it's weird.

> _We didn't even need the software for `gpg --verify`. How does that make sense? I thought we want to verify the authenticity and integrity of it?_

You feel betrayed but continue reading the download page and find that it mentions this:

> You have now verified the signature of the manifest file, which ensures integrity and authenticity of the manifest file - not the binaries! Next, depending on your operating system, you must re-compute the sha256 hash of the archive with `shasum -a 256 <filename>`. First, download the installation for your operating system (if you haven’t done so already). Then follow the steps below to compare it with the corresponding one in the manifest file, and ensure they match exactly.

> _Aha, so I verified the manifest file._

Looking at the manifest file again, you see that it contains a line for every binary that one can download. The binary that you downloaded is therefore also included (else your feelings of betrayal would have intensified):

```
ffb7f86e978ab312dcc169d7c70b33048256dcf3280aecb1cb9392124eac31b9 *sparrow-1.8.2-x86_64.tar.gz
```

You have some experience with hashes and can recognize hashes like `ffb7f86e978ab312dcc169d7c70b33048256dcf3280aecb1cb9392124eac31b9` when you see them. But you didn't need to since the download page already mentioned that the file contains sha256 hashes.

It also mentioned that we should run `shasum -a 256 <filename>`. So we do exactly that and expect that it returns the same hash:

```
$ shasum -a 256 sparrow-1.8.2-x86_64.tar.gz
ffb7f86e978ab312dcc169d7c70b33048256dcf3280aecb1cb9392124eac31b9  sparrow-1.8.2-x86_64.tar.gz
```

That looks about right! But it's pretty annoying to compare hashes by hand. The download page also has a solution for that. It mentions the following command for Linux users:

```
sha256sum --check sparrow-1.8.2-manifest.txt --ignore-missing
```

So we run this and expect more intimacy via hand-holding:

```
$ sha256sum --check sparrow-1.8.2-manifest.txt --ignore-missing
sparrow-1.8.2-x86_64.tar.gz: OK
```

It says "OK". Well, that's not very intimate like what a lover probably would say but still sounds like a friend. Everything seems to be okay!

# Conclusion

So what we just did was to basically verify the authenticity and integrity of the file that contained the hashes for all binaries with `gpg --verify`. When the hashes could be trusted, we could use them to make sure that the software was not tampered with. But why not simply provide a digital signature for the binary itself?

I actually don't know. But my educated guess is that it's related to convenience. Instead of providing a signature for every binary, the hashes are signed. Using `sha256sum --check` with `--ignore-missing` then simply ignores all files that don't exist. So I am basically guessing that there is no way to do something similar with digital signatures. Maybe someone knows more?

# Wait ... but why (does it work)?

![17250.gif](https://m.stacker.news/17277)

Essentially, if you believe in cryptographic hashes (like during bitcoin mining), you believe in digital signatures.

_If you don't know what a hash is (and especially not cryptographic ones), read [this](https://blog.cloudflare.com/how-developers-got-password-security-so-wrong). If you still don't know, you can ask questions in the comments. I am sure someone will be willing to answer them._

A digital signature is usually nothing else than a cryptographic hash that was "encrypted" with a secret key:

> The private key can be used to create a digital signature for any piece of data using a digital signature algorithm. This typically involves taking a cryptographic hash of the data and operating on it mathematically using the private key. Anyone with the public key can check that this signature was created using the private key and the appropriate signature validation algorithm. A digital signature is a powerful tool because it allows you to publicly vouch for any message.

— blog.cloudflare.com, [_ECDSA: The digital signature algorithm of a better internet_](https://blog.cloudflare.com/ecdsa-the-digital-signature-algorithm-of-a-better-internet)

This means that during verification with the public key, we reverse the encryption with the private key to reveal the hash. Then, we hash the signed data ourselves and compare the hashes. If the hashes match, the integrity was verified and since public key cryptography was used, we also verified the authenticity of the hashes.

---

Makes sense? Still reading?

If so, you potentially even learned a lot about how _digital certificates_ work on accident. They are basically digital signatures that verify the identity of entities via trusted third parties ([Certificate Authorities](https://www.digicert.com/blog/what-is-a-certificate-authority)) and with more metadata associated. At least that's my simplistic explanation. But how the [Public Key Infrastructure](https://www.okta.com/identity-101/public-key-infrastructure/) works is a topic for another day.



[^1]: I'll assume some basic Linux skills since those are more often than not very handy. You might even acquire the "required" Linux skills by accident while reading this just by the context the commands are used in—just like human languages.

[^2]: This is a step that is often overlooked or not even mentioned. That the public key is trustworthy is the foundation for this whole scheme. You SHOULD NOT simply trust keys. You should do some basic due diligence. Ask friends if they have the same public key. If you're fucked, at least they are fucked, too. Try to find key fingerprints and compare them using `gpg --list-keys --keyid-format long`. But surely someone else would have noticed if the key is compromised, right? Right? But I get it. It's just too convenient to simply trust keys in most cases.
