#!/bin/bash
set -euo pipefail

REPO="fisherevans/ha2trmnl"

# Determine tag
if [[ $# -gt 0 ]]; then
  NEW_TAG="$1"
else
  LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
  IFS='.' read -r MAJOR MINOR PATCH <<<"${LAST_TAG#v}"
  NEW_TAG="v${MAJOR}.${MINOR}.$((PATCH + 1))"
fi

echo "[publish] Using tag: ${NEW_TAG}"

# Create buildx builder if needed
if ! docker buildx inspect multiarch-builder &>/dev/null; then
  docker buildx create --name multiarch-builder --use
fi

# Build and push multi-arch image
echo "[publish] Building and pushing image for linux/amd64 and linux/arm64"
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --tag "${REPO}:${NEW_TAG}" \
  --tag "${REPO}:latest" \
  --file docker/Dockerfile \
  --push \
  .

# Tag the Git repo
echo "[publish] Creating Git tag ${NEW_TAG}"
git tag "${NEW_TAG}"
git push origin "${NEW_TAG}"

echo "[publish] Done."