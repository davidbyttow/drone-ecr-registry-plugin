#!/bin/bash
set -ex

# Specify PATH to Dockerfile as $1, image tag as $2
ROOT=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )

docker build -f $1 -t $2 ${ROOT}
