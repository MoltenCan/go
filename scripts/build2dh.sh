#!/bin/sh
echo "build 2 docker hub"

[ -z $1 ] && {
    echo "need a target"
    echo " starter |"
    exit 1
}

CMD_DIR="cmd/$1"
[ -e $CMD_DIR ] || {
    echo "$1 does not exist"
    exit 1
}

platforms="linux/arm/v7,linux/arm/v6,linux/amd64,linux/arm64/v8"
ver="0.0.1"

cd $CMD_DIR
docker buildx build \
    --platform "${platforms}" \
    --tag moltencan/${1}:${ver} \
    --tag moltencan/${1}:latest \
    --push \
    .
