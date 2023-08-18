#!/usr/bin/env bash

function sync() {
  ENV=development make run $@
  date +%s.%N > public/hot-reload
  rsync -avh public/ dev.ekzyis.com:/var/www/dev.ekzyis --delete
}

function cleanup() {
    rm -f public/hot-reload
}
trap cleanup EXIT

sync -B
while inotifywait -r -e modify html/ blog/ *.go; do
  sync
done
