**💻 popular commands**

```
# enter nix shell
$ nix develop

# serve site at localhost:8080
$ nix run

# website configuration
$ s3cmd ws-create --ws-index=index.html --ws-error=404.html s3://www.ekzy.is

# deploy site
$ ./deploy

# configure cors
$ s3cmd setcors cors.xml s3://www.ekzy.is

# request cert manually and upload via Linode web interface
$ sudo certbot certonly --manual

# stacker.news links for @-mentions
$ sed -E 's/(^|[^a-zA-Z0-9])@([a-zA-Z0-9_-]+)/\1[@\2](https:\/\/stacker.news\/\2)/g'
```

**📑 useful links**

- [Upload your Static Site to Linode Object Storage | Linode Docs](https://www.linode.com/docs/guides/host-static-site-object-storage/#upload-your-static-site-to-linode-object-storage)
