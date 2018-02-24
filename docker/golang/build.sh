#!/bin/bash
set -e
echo "[build.sh:building binary]"

# linux86_64
cd $BUILDPATH && \
   GOOS=linux GOARCH=amd64 go build -o ckp && \
   mv ./ckp /usr/bin/ && \
   rm -rf /tmp/*

chmod u+x /usr/bin/ckp

echo "[build.sh:launching binary]"
ckp