#!/bin/bash -e
j2 nginx.conf.j2 > /etc/nginx/nginx.conf
j2 nginx-site.conf.j2 > /etc/nginx/conf.d/default.conf
exec "$@"
