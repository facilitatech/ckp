#!/bin/bash
set -e
echo "[build.sh:building binary]"
cd $BUILDPATH && go build -o /app && rm -rf /tmp/*

cp /app /usr/bin/
chmod u+x /usr/bin/app

echo "[build.sh:launching binary]"
/app
