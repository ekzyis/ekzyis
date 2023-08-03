#!/usr/bin/env bash

function sync() {
  go run . --dev
  date +%s.%N > public/hot-reload
  rsync -avh public/ dev.ekzyis.com:/var/www/dev.ekzyis --delete
}

function cleanup() {
    rm -f public/hot-reload
}
trap cleanup EXIT

sync
while inotifywait -r -e modify html/; do
  sync
done
