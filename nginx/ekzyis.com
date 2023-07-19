server {
    server_name ekzyis.com;
    listen      80;
    listen      [::]:80;

    return 301 https://ekzyis.com$request_uri;
}

server {
    server_name         ekzyis.com;
    listen              443 ssl;
    listen              [::]:443 ssl;

    ssl_certificate     /etc/letsencrypt/live/ekzyis.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/ekzyis.com/privkey.pem;

    root /var/www/ekzyis;
    index index.html;

    error_page 404 /404.html;

    include letsencrypt.conf;
}
