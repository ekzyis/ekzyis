#!/usr/bin/env bash

function sync() {
  go run .
  rsync -avh public/ dev.ekzyis.com:/var/www/dev.ekzyis --delete
}

sync
while inotifywait -r -e modify html/; do
  sync
done
