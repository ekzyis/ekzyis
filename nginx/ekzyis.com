proxy_cache_path /var/run/nginx-cache/jscache levels=1:2 keys_zone=jscache:100m inactive=30d  use_temp_path=off max_size=100m;

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
    try_files $uri.html $uri $uri/ =404;

    error_page 404 /404.html;

    resolver 9.9.9.9;
    set $plausible_script_url https://plausible.io/js/script.js;
    set $plausible_event_url https://plausible.io/api/event;

    location = /js/script.js {
      proxy_pass $plausible_script_url;
      proxy_set_header Host plausible.io;

      # Tiny, negligible performance improvement. Very optional.
      proxy_buffering on;

      # Cache the script for 6 hours, as long as plausible.io returns a valid response
      proxy_cache jscache;
      proxy_cache_valid 200 6h;
      proxy_cache_use_stale updating error timeout invalid_header http_500;

      # Optional. Adds a header to tell if you got a cache hit or miss
      add_header X-Cache $upstream_cache_status;
    }

    location = /api/event {
      proxy_pass $plausible_event_url;
      proxy_set_header Host plausible.io;
      proxy_buffering on;
      proxy_http_version 1.1;

      proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
      proxy_set_header X-Forwarded-Host  $host;
    }

    include letsencrypt.conf;
}
