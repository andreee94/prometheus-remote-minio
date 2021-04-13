#!/bin/bash

if [[ "$EUID" = 0 ]]; then
    echo "OK, already root"
else
    echo "Run again with sudo."
    exit
fi

# source ./env

sudo -E docker-compose build

sudo -E docker-compose push
