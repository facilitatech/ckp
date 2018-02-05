#!/bin/bash
set -e
echo "[build.sh:building binary]"

cd $BUILDPATH && GOOS=linux GOARCH=amd64 go build -o ckp_linux86_64 && rm -rf /tmp/*
cd $BUILDPATH && GOOS=darwin GOARCH=amd64 go build -o ckp_darwinx86_64 && rm -rf /tmp/*

cp ckp_linux86_64 /usr/bin/
cp ckp_darwinx86_64 /usr/bin/

chmod u+x /usr/bin/ckp_linux86_64

echo "[build.sh:launching binary]"
ckp_linux86_64
