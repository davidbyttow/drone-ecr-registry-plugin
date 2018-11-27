#!/bin/bash
set -ex

ROOT=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )
cd "${ROOT}"

mkdir -p $1

cd "${ROOT}"

package_root="github.com/davidbyttow/drone-ecr-registry-plugin"

version_ldflags=""
if [[ -n "${2}" ]]; then
  version_ldflags="-X ${package_root}/version.Version=${2}"
fi

if [[ -n "${3}" ]]; then
  version_ldflags="$version_ldflags -X ${package_root}/version.GitCommitSHA=${3}"
fi

GOOS=$TARGET_GOOS GOARCH=$TARGET_GOARCH CGO_ENABLED=0 \
       	go build -installsuffix cgo -a -ldflags "-s ${version_ldflags}" \
       	-o $1/drone-ecr-registry-plugin \
        ./cmd
