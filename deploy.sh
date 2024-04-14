#!/usr/bin/env bash

set -e

hugo
rsync -avh public/ ekzy.is:/var/www/ek --delete --dry-run

echo
read -p "Continue deploy? [yn] " yn
echo
[ "$yn" == "y" ] && rsync -avh public/ ekzy.is:/var/www/ek --delete
