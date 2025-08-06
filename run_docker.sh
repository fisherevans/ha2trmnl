#!/bin/bash
set -euo pipefail

TAG=fisherevans/ha2trmnl
CONFIG_PATH=$(realpath config.yaml)

echo "[run] Running Docker image with config: $CONFIG_PATH"
docker run --rm --init \
  -v "$CONFIG_PATH":/config/config.yaml \
  "$TAG" "$@"