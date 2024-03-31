#!/usr/bin/env bash

set -e

ENV=production make -B
rsync -avh .well-known public/ ekzyis.com:/var/www/ekzyis --delete --dry-run

echo
read -p "Continue deploy? [yn] " yn
echo
[ "$yn" == "y" ] && rsync -avh .well-known public/ ekzyis.com:/var/www/ekzyis --delete
