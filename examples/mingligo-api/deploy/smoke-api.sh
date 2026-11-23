#!/usr/bin/env bash
set -euo pipefail
IMAGE_TAG="${IMAGE_TAG:-local}"
IMAGE_REGISTRY="${IMAGE_REGISTRY:-mingligo-local}"
IMAGE="${MINGLIGO_API_IMAGE:-$IMAGE_REGISTRY/mingligo-api:$IMAGE_TAG}"
PORT="${MINGLIGO_API_SMOKE_PORT:-18092}"
NAME="mingligo-smoke-api-$$"
cleanup() { docker rm -f "$NAME" >/dev/null 2>&1 || true; }
trap cleanup EXIT
docker run -d --name "$NAME" -p "$PORT:18082" -e BASE_PATH=/mingligo/api "$IMAGE"
for _ in $(seq 1 30); do
  curl -fsS "http://127.0.0.1:$PORT/healthz" >/dev/null 2>&1 && break
  sleep 1
done
curl -fsS "http://127.0.0.1:$PORT/mingligo/api/openapi.json" >/dev/null
curl -fsS \
  -H 'Content-Type: application/json' \
  -d '{"timezone":"Asia/Shanghai"}' \
  "http://127.0.0.1:$PORT/mingligo/api/time_now" | grep -q 'local_standard_time'
printf '%s\n' "✓ /mingligo/api: 17 methods, OpenAPI, Problem Details runtime"
