#!/bin/sh

set -e

echo "\033[44m\033[97mTesting\033[0m Debian"
docker run --rm -it "$(docker build -q --no-cache -f debian/Dockerfile ../../)"

echo "\033[44m\033[97mTesting\033[0m Ubuntu"
docker run --rm -it "$(docker build -q --no-cache -f ubuntu/Dockerfile ../../)"

echo "\033[44m\033[97mTesting\033[0m Red Hat"
docker run --rm -it "$(docker build -q --no-cache -f redhat/Dockerfile ../../)"

echo "\033[42m\033[97mSuccess\033[0m"
