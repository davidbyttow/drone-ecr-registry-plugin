#!/bin/bash

# This script is used for compilation of a specific variant.
# Specify GOOS as $1, GOARCH as $2
# Binaries are placed into ./bin/$GOOS-$GOARCH/docker-ecr-registry-plugin

ROOT=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )
cd "${ROOT}"

. ./scripts/shared_env

export TARGET_GOOS="$1"
export TARGET_GOARCH="$2"
VERSION="$3"
GITCOMMIT_SHA="$4"

./scripts/build_binary.sh "./bin/${TARGET_GOOS}-${TARGET_GOARCH}" $VERSION $GITCOMMIT_SHA

echo "Built drone-ecr-registry-plugin for ${TARGET_GOOS}-${TARGET_GOARCH}-${VERSION}"
