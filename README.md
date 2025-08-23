![](./content/about/psychedelic-digital-sky.jpg)

<div align="center">

🌐 [ekzy.is](https://ekzy.is) | 💬 [nostr](https://njump.me/npub16x07c4qz05yhqe2gy2q2u9ax359d2lc0tsh6wn3y70dmk8nv2j2s96s89d) [signal](https://signal.me/#eu/QQJWrLHuZ-qRrNxo8x1CygWeU9ITJkrCkHg7Sm0vx4WfxB9y5PJM-aPINkauSLxb) | 🔑 [EBAF 75DA 7279 CB48](https://keybase.io/ekzyis)

</div>

```
$ s3cmd ws-create --ws-index=index.html --ws-error=404.html s3://www.ekzy.is
$ s3cmd sync --guess-mime-type --no-mime-magic --acl-public --delete-removed --delete-after public/ s3://www.ekzy.is
$ sudo certbot certonly --manual
```

https://www.linode.com/docs/guides/host-static-site-object-storage/#upload-your-static-site-to-linode-object-storage