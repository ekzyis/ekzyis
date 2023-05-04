#!/usr/bin/env bash

rsync -v public/*.{html,css,ico,png,webmanifest} vps:/var/www/ekzyis
