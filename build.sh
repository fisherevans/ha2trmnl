#!/bin/bash
set -euo pipefail

TAG=fisherevans/ha2trmnl:latest

echo "[build] Building Docker image as $TAG"
docker build -t "$TAG" -f docker/Dockerfile .